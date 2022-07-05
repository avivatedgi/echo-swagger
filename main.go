package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/avivatedgi/echo-swagger/echo_swagger"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	log.SetLevel(log.DebugLevel)

	var infoFile *os.File = nil

	output := os.Stdout

	directory := flag.String("dir", "", "Directory to scan for request handlers")
	pattern := flag.String("pattern", "./...", "Package pattern to scan for request handlers")

	flag.Func("out", "Path to file output to write in the geerated OpenAPI specifications", func(s string) error {
		if s == "-" {
			return nil
		}

		f, err := os.OpenFile(s, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}

		output = f
		return nil
	})

	flag.Func("info", "`path` to info file", func(s string) error {
		f, err := os.Open(s)
		if err != nil {
			return err
		}

		infoFile = f
		return nil
	})

	// Parse the actual arguments from the command line
	flag.Parse()

	// Validate the arguments from the command line
	if infoFile == nil {
		log.Fatal("info file is required!")
	} else if directory == nil || *directory == "" {
		log.Fatal("directory is required!")
	} else if pattern == nil || *pattern == "" {
		log.Fatal("pattern is required!")
	}

	// Read the info file
	infoData, err := ioutil.ReadAll(infoFile)
	if err != nil {
		log.Fatal("Failed to read info file `", infoFile.Name(), "`, error = ", err)
	}

	// Unmarshal the info file
	info := echo_swagger.Info{}
	if err := yaml.Unmarshal(infoData, &info); err != nil {
		log.Fatal("Failed to unmarshal info file, error = ", err)
	}

	if output != os.Stdout {
		defer output.Close()
	}

	parser := echo_swagger.NewContext()

	// Parse the directory
	err = parser.ParseDirectory(*directory, *pattern)
	if err != nil {
		log.Fatal("Failed to parse directory ", directory, ", error = ", err)
	}

	parser.OpenAPI.Info = info

	// Marshal the generated OpenAPI specifications
	data, err := yaml.Marshal(parser.OpenAPI)
	if err != nil {
		log.Fatal("Failed to marshal generated OpenAPI specifications: ", err)
	}

	// Write the OpenAPI specifications to the output file
	if _, err := output.Write(data); err != nil {
		log.Fatal("Failed to write generated OpenAPI specifications to file ", output.Name(), ", error = ", err)
	}
}
