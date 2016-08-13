package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	gitRepo string
	topDir  string
)

func init() {
	flag.StringVar(&gitRepo, "gitrepo", "", "")
	flag.StringVar(&topDir, "top-dir", "", "")
	flag.Parse()
}

func ExitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
func main() {
	if gitRepo == "" {
		ExitError(errors.New("Need to set git repo"))
	}
	basePath := path.Join(gitRepo, topDir)
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		p, err := filepath.Rel(gitRepo, path)
		if err != nil {
			return err
		}
		for _, ignoredPrefix := range []string{"_output/", "vendor/"} {
			if strings.HasPrefix(p, ignoredPrefix) {
				return nil
			}
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() != "OWNERS" {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		names := []string{}
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "- ") {
				continue
			}
			var name string
			fmt.Sscanf(line, "- %s", &name)
			names = append(names, name)
		}
		if len(names) > 0 {
			fmt.Printf("path: %s, owners: %v\n", filepath.Dir(p), names)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading file (%s): %v\n", path, err)
		}
		return nil
	})
	if err != nil {
		ExitError(err)
	}
}
