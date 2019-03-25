package file

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

var openedFiles = make(map[string][]string)

func ProcessEvent(ev string, tokenRegex *regexp.Regexp, replacement string) string {
	var replacementSamples []string
	var ok bool

	// Check if we've opened the file already
	if replacementSamples, ok = openedFiles[replacement]; !ok {
		if err := openFile(replacement); err != nil {
			log.Fatalf("Could not open replacement sample file '%s': %v", replacement, err)
		}
	}

	// Grab random
	replacementSamples = openedFiles[replacement]
	n := rand.Int() % len(replacementSamples)

	matches := tokenRegex.FindStringSubmatch(ev)
	if len(matches) > 0 {
		return strings.Replace(ev, matches[len(matches) - 1], replacementSamples[n], 1)
	}
	return ev
}

func openFile(sampleFile string) error {
	f, err := os.Open(sampleFile)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	samples := strings.Split(string(data), "\n")
	openedFiles[sampleFile] = samples
	return nil
}
