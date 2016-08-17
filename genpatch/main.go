package main

import (
	"bufio"
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var gitRepo string

func init() {
	flag.StringVar(&gitRepo, "gitrepo", "", "")
	flag.Parse()
}

func patch(path string, rs []string) {
	p := filepath.Join(gitRepo, path, "OWNERS")
	// fmt.Println(p)

	b := &bytes.Buffer{}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		b.WriteString("assignees:\n")
		for _, r := range rs {
			b.WriteString("  - ")
			b.WriteString(r)
			b.WriteString("\n")
		}

		err := ioutil.WriteFile(p, b.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
		return
	}

	f, err := os.OpenFile(p, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	names := map[string]struct{}{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) != 2 || fields[0] != "-" {
			continue
		}
		name := fields[1]
		names[name] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for _, r := range rs {
		if _, ok := names[r]; ok {
			continue
		}
		b.WriteString("  - ")
		b.WriteString(r)
		b.WriteString("\n")
	}
	_, err = f.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

// How to run it:
//   cat $OWNER_FILE | ./genpatch --gitrepo="$GOPATH/src/k8s.io/kubernetes"
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		segs := strings.Split(line, ",")
		if len(segs) != 2 {
			panic("unexpected")
		}
		pathSegs := strings.Split(segs[0], ":")
		if len(pathSegs) != 2 {
			panic("unexpected")
		}
		path := strings.TrimSpace(pathSegs[1])

		reviwerSegs := strings.Split(segs[1], ":")
		if len(reviwerSegs) != 2 {
			panic("unexpected")
		}
		reviewerListStr := strings.TrimSpace(reviwerSegs[1])
		reviewers := strings.Split(reviewerListStr[1:len(reviewerListStr)-1], " ")

		patch(path, reviewers)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
