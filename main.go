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

	"github.com/secoura/eventreplay/output/console"
	"github.com/secoura/eventreplay/processor/file"
	"github.com/secoura/eventreplay/processor/timestamp"

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

	if *configFile == "" {
		if err = tryParseFromEnvVar(); err != nil {
			printHelp()
		}
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
	// Try load config from PLUGIN_CONFIG environment variable
	pluginConf := os.Getenv("PLUGIN_CONFIG")
	if pluginConf == "" {
		return errors.New("PLUGIN_CONFIG not set")
	}
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

	switch repl.Type {
	case "timestamp":
		return timestamp.ProcessEvent(ev, tokRegex, repl.Replacement, cfg.EarliestTime, cfg.LatestTime)
	case "file":
		return file.ProcessEvent(ev, tokRegex, repl.Replacement)
	}

	log.Fatalf("unknown processor '%s'", repl.Type)
	return ""
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