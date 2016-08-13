package main

import (
	"bufio"
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
	p := filepath.Join(gitRepo, path, "REVIEWER")
	// fmt.Println(p)

	out := strings.Join(rs, "\n")
	// fmt.Println(out)

	if err := ioutil.WriteFile(p, []byte(out), 0644); err != nil {
		panic(err)
	}
}

// cat ../_output/owner.txt | ./genpatch --gitrepo="$GOPATH/src/k8s.io/kubernetes"

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
