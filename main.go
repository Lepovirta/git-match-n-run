package main

import (
	"bytes"
	"flag"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
)

var (
	configFilename string
	fromRef        string
	toRef          string
	shouldRun      bool
)

type matchRunConfig []matchRunEntry

type matchRunEntry struct {
	Pattern string   `yaml:"pattern"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
}

func (m *matchRunEntry) match(filenames []string) (bool, error) {
	r, err := regexp.Compile(m.Pattern)
	if err != nil {
		return false, err
	}
	for _, filename := range filenames {
		if r.MatchString(filename) {
			return true, nil
		}
	}
	return false, nil
}

func (m *matchRunEntry) run() error {
	cmd := exec.Command(m.Command, m.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (m *matchRunEntry) commandString() string {
	var bytes bytes.Buffer
	bytes.WriteString(m.Command)
	for _, arg := range m.Args {
		bytes.WriteString(" ")
		bytes.WriteString(arg)
	}
	return bytes.String()
}

func main() {
	flag.StringVar(&configFilename, "config", "gitmatchrun.yaml", "Config file for git-match-run")
	flag.StringVar(&fromRef, "from", "", "Git ref to start finding changes from")
	flag.StringVar(&toRef, "to", "HEAD", "Git ref to end finding changes")
	flag.BoolVar(&shouldRun, "run", false, "Whether or not to run the matching commands")
	flag.Parse()

	var config matchRunConfig
	if err := readConfig(&config); err != nil {
		log.Panic(err)
	}

	changedFiles, err := getChangedFiles()
	if err != nil {
		log.Panic(err)
	}

	if err := runEntries(config, changedFiles); err != nil {
		log.Panic(err)
	}
}

func readConfig(config *matchRunConfig) error {
	bs, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bs, config)
}

func getChangedFiles() ([]string, error) {
	changedFiles := make([]string, 0, 100)
	cmd := exec.Command("git", "diff", "-z", "--name-only", fromRef, toRef)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	outBuf := bytes.NewBuffer(out)
	for line := ""; err == nil; line, err = outBuf.ReadString(0) {
		if line != "" {
			changedFiles = append(changedFiles, line[:len(line)-1])
		}
	}
	if err != io.EOF {
		return nil, err
	}
	return changedFiles, nil
}

func runEntries(config matchRunConfig, changedFiles []string) error {
	for _, entry := range config {
		matchFound, err := entry.match(changedFiles)
		if err != nil {
			return err
		}

		if matchFound {
			log.Printf("Running %s", entry.commandString())
			if shouldRun {
				if err := entry.run(); err != nil {
					return err
				}
			}
		} else {
			log.Printf("No changes in files matching /%s/", entry.Pattern)
		}
	}
	return nil
}
