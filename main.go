package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/avivatedgi/echo-swagger/echo_swagger"
	"gopkg.in/yaml.v3"
)

type arrayFlag []string

func (i *arrayFlag) String() string {
	return ""
}

func (i *arrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var infoFile *os.File = nil

	output := os.Stdout
	directories := arrayFlag{}

	flag.Var(&directories, "dir", "Directories to scan for swagger routes, make sure to pass declaration files before the routes")
	flag.Func("out", "`path` to file output (default STDOUT)", func(s string) error {
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
	flag.Parse()

	if infoFile == nil {
		fmt.Fprintf(os.Stderr, "Info file is required\n")
		os.Exit(1)
	}

	infoData, err := ioutil.ReadAll(infoFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read info file `%s`, error = `%s`\n", infoFile.Name(), err)
		os.Exit(1)
	}

	info := echo_swagger.Info{}
	if err := yaml.Unmarshal(infoData, &info); err != nil {
		panic(err)
	}

	if output != os.Stdout {
		defer output.Close()
	}

	parser := echo_swagger.New()

	openapi, err := parser.ParseDirectories([]string(directories))
	if err != nil {
		panic(err)
	}

	openapi.Info = info

	data, err := yaml.Marshal(openapi)
	if err != nil {
		panic(err)
	}

	output.Write(data)
}
