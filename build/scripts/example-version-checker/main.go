package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var exlcudedPatterns = [...]string{"*.md", "*.yaml", "OWNERS", ".gitignore"}

func main() {
	byteContents, err := os.ReadFile("/Users/edwardmoulsdale/Projects/agones/examples/xonotic/Makefile")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	contents := string(byteContents)
	version, err := getVersionFromMakefile(contents)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println(version)
}

func getVersionFromMakefile(contents string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(contents))
	for scanner.Scan() {
		line := scanner.Text()
		if lineContainsVersion(line) {
			return getVersionFromLine(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Could not get version: %v\n", err)
	}

	return "", errors.New("No version string was found.")
}

func lineContainsVersion(line string) bool {
	return strings.HasPrefix(line, "version")
}

func getVersionFromLine(line string) (string, error) {
	split := strings.SplitN(line, ":=", 2)
	if len(split) != 2 {
		return "", fmt.Errorf("Bad version line: %s\n", line)
	}
	trimmed := strings.TrimSpace(split[1])
	if trimmed == "" {
		return "", errors.New("Version can not be empty.")
	}
	return trimmed, nil
}
