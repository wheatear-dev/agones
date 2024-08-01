package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const cloneDir = "clone"
const examplesDir = "../../../examples" // sort this out later
const gitUrl = "https://github.com/googleforgames/agones.git"
const targetBranch = "refs/heads/main"

var exlcudedPatterns = [...]string{"*.md", "*.yaml", "OWNERS", ".gitignore"}

func main() {
	names, err := getExampleDirNames(examplesDir)
	if err != nil {
		fmt.Print(err)
	} else {
		for _, name := range names {
			fmt.Println(name)
		}
	}
}

func cloneRepo(gitDir string) (*git.Repository, error) {
	cloneOptions := &git.CloneOptions{
		URL: gitUrl,
	}
	return git.PlainClone(gitDir, false, cloneOptions)
}

func fetchTargetCommit(repo *git.Repository) *object.Commit {
	targetRef, err := repo.Reference(plumbing.ReferenceName(targetBranch), true)
	if err != nil {
		log.Fatalf("Could refernce to main: %v\n", err)
	}

	targetCommit, err := repo.CommitObject(targetRef.Hash())
	if err != nil {
		log.Fatalf("Failed to get HEAD on main: %v\n", err)
	}

	return targetCommit
}

func dirIsExample(dirName string) bool {
	makefileName := fmt.Sprintf("%s/Makefile", dirName)
	if _, err := os.Stat(makefileName); err == nil {
		return true
	} else {
		return false
	}
}

func getExampleDirNames(baseDir string) ([]string, error) {
	dirNames := make([]string, 0)

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return dirNames, fmt.Errorf("Count not open directory: %v\n", err)
	}

	for _, entry := range entries {
		name := fmt.Sprintf("%s/%s", baseDir, entry.Name())
		if dirIsExample(name) {
			dirNames = append(dirNames, name)
		}
	}
	return dirNames, nil
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
