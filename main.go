package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/secoura/eventreplay/processor/file"
	"github.com/secoura/eventreplay/processor/float"
	"github.com/secoura/eventreplay/processor/hex"
	"github.com/secoura/eventreplay/processor/integer"
	"github.com/secoura/eventreplay/processor/ipv4"
	"github.com/secoura/eventreplay/processor/ipv6"
	"github.com/secoura/eventreplay/processor/list"
	"github.com/secoura/eventreplay/processor/mac"
	processorString "github.com/secoura/eventreplay/processor/string"
	"github.com/secoura/eventreplay/processor/timestamp"
	"github.com/secoura/eventreplay/processor/uuid"

	"github.com/secoura/eventreplay/output/console"

	"github.com/secoura/eventreplay/config"
)

func printHelp() {
	args := []string{
		"[-c <config file>]",
	}
	fmt.Printf("Usage: " + os.Args[0] + "\n\t" +
		strings.Join(args, "\n\t") +
		"\n")
	os.Exit(1)
}

func printError(err error) {
	log.Fatalf("[ERROR] %v", err)
}

var cfg config.Config
var debug *bool
var regexCache = make(map[string]*regexp.Regexp)

func main() {
	var err error

	configFile := flag.String("c", "", "config file")
	debug = flag.Bool("debug", false, "enable debugging mode")
	flag.Parse()

	if os.Getenv("PLUGIN_CONFIG") != "" {
		if err = tryParseFromEnvVar(); err != nil {
			printError(err)
		}
	}
	if *configFile == "" {
		printHelp()
	} else {
		if err = tryParseFromFile(*configFile); err != nil {
			printError(err)
		}
	}

	// Validate samples were provided
	if len(cfg.Samples) == 0 {
		fmt.Printf("Config file does not have any 'samples' defined.")
	}

	// Wait for close signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT)

	generateSamples()
	<-sigCh

	if *debug {
		fmt.Printf("Closing ...\n")
	}
}

func tryParseFromFile(cfgFile string) error {
	f, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &cfg)
}

func tryParseFromEnvVar() error {
	pluginConf := os.Getenv("PLUGIN_CONFIG")
	if err := json.Unmarshal([]byte(pluginConf), &cfg); err != nil {
		return errors.Errorf("failed to parse PLUGIN_CONFIG: %v", err)
	}
	return nil
}

func generateSamples() {
	for _, s := range cfg.Samples {
		go func() {
			if err := runSampleLoop(s); err != nil {
				log.Printf("ERROR: %v\n", err)
			}
		}()
	}
}

func runSampleLoop(sampleCfg config.SampleConfig) error {
	if sampleCfg.InputFile == "" {
		return errors.Errorf("input file not defined for sample")
	}

	// Open input file
	sampleFile, err := os.Open(sampleCfg.InputFile)
	if err != nil {
		return err
	}
	sampleData, err := ioutil.ReadAll(sampleFile)
	if err != nil {
		return err
	}
	sampleFile.Close()

	// Delimit the sample data
	delimiterRegex := regexp.MustCompile(sampleCfg.Delimiter)
	events := delimiterRegex.Split(string(sampleData), -1)

	if *debug {
		log.Printf("got %d events\n", len(events))
	}

	// Check limit
	limit := len(events)
	if sampleCfg.Count > 0 {
		limit = sampleCfg.Count
	}

	// Create timer to loop
	timer := time.NewTimer(sampleCfg.Interval)
	for {
		select {
		case <-timer.C:
			for i := 0; i < limit; i++ {
				ev := events[i]
				for _, repl := range sampleCfg.Replacements {
					ev = processEvent(ev, sampleCfg, repl)
				}
				outputEvent(ev, sampleCfg)
				timer.Reset(sampleCfg.Interval)
			}
		}
	}
}

func processEvent(ev string, cfg config.SampleConfig, repl config.Replacement) string {
	var tokRegex *regexp.Regexp
	var ok bool

	if tokRegex, ok = regexCache[repl.Token]; !ok {
		tokRegex = regexp.MustCompile(repl.Token)
		regexCache[repl.Token] = tokRegex
	}

	var replacementVal string
	switch repl.Type {
	case "guid":
		replacementVal = uuid.ProcessEvent()
	case "ipv4":
		replacementVal = ipv4.ProcessEvent()
	case "ipv6":
		replacementVal = ipv6.ProcessEvent()
	case "mac":
		replacementVal = mac.ProcessEvent()
	case "integer":
		replacementVal = integer.ProcessEvent(repl.Replacement)
	case "float":
		replacementVal = float.ProcessEvent(repl.Replacement)
	case "string":
		replacementVal = processorString.ProcessEvent(repl.Replacement)
	case "hex":
		replacementVal = hex.ProcessEvent(repl.Replacement)
	case "list":
		replacementVal = list.ProcessEvent(repl.Replacement)
	case "timestamp":
		replacementVal = timestamp.ProcessEvent(repl.Replacement, cfg.EarliestTime, cfg.LatestTime)
	case "file":
		replacementVal = file.ProcessEvent(repl.Replacement)
	case "static":
		replacementVal = repl.Replacement.(string)
	default:
		log.Fatalf("unknown processor '%s'", repl.Type)
	}

	matches := tokRegex.FindStringSubmatch(ev)
	if len(matches) > 0 {
		return strings.Replace(ev, matches[len(matches)-1], replacementVal, 1)
	}
	return ev
}

func outputEvent(ev string, cfg config.SampleConfig) {
	switch cfg.OutputMode {
	case "console":
		console.ProcessEvent(ev, cfg.OutputCodec, cfg.Identifier)
	default:
		log.Fatalf("unknown output '%s'", cfg.OutputMode)
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
