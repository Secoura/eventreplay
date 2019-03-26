package file

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
)

var openedFiles = make(map[string][]string)

func ProcessEvent(replacement interface{}) string {
	var repString = replacement.(string)
	var replacementSamples []string
	var ok bool

	// Check if we've opened the file already
	if replacementSamples, ok = openedFiles[repString]; !ok {
		if err := openFile(repString); err != nil {
			log.Fatalf("Could not open replacement sample file '%s': %v", replacement, err)
		}
	}

	// Grab random
	replacementSamples = openedFiles[repString]
	n := rand.Int() % len(replacementSamples)

	return replacementSamples[n]
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
