package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	githubToken string
	githubOrg   string
	githubRepo  string
	topLevelDir string
)

var now time.Time

var excludedDirList []string

const levelLimit = 3

func ExitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func init() {
	flag.StringVar(&githubToken, "github-token", "", "")
	flag.StringVar(&githubOrg, "github-org", "kubernetes", "")
	flag.StringVar(&githubRepo, "github-repo", "kubernetes", "")
	flag.StringVar(&topLevelDir, "top-dir", "", "")
	flag.Parse()

	now = time.Now()

	// TODO: This is hardcoded currently. parse from a file
	excludedDirList = []string{
		"vendor",
		"contrib/mesos/",
		// exclude generated code: `find . | grep "generated"` + some guessing
		"staging",
		"cmd/libs/go2idl/client-gen",
		"federation/client/clientset_generated",
		"pkg/client/clientset_generated",
	}
}

func main() {
	if githubToken == "" {
		ExitError(errors.New("Github token isn't provided"))
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	fetchOwners(client, topLevelDir, 0)
}

func fetchOwners(client *github.Client, dir string, level int) {
	_, directoryContent, _, err := client.Repositories.GetContents(githubOrg, githubRepo, dir, &github.RepositoryContentGetOptions{})
	if err != nil {
		ExitError(err)
	}
	fetchTopCommitters(client, dir, 3)

	if level >= levelLimit {
		return
	}

	for _, c := range directoryContent {
		if c.Type != nil && *c.Type == "dir" {
			fetchOwners(client, *c.Path, level+1)
		}
	}
}

func fetchTopCommitters(client *github.Client, dir string, limit int) {
	for _, prefix := range excludedDirList {
		if strings.HasPrefix(dir, prefix) {
			return
		}
	}

	opt := &github.CommitsListOptions{
		Path:  dir,
		Since: now.AddDate(0, -12, 0),
		ListOptions: github.ListOptions{
			PerPage: 500,
		},
	}
	rank := map[string]int{}
	for {
		commits, resp, err := client.Repositories.ListCommits(githubOrg, githubRepo, opt)
		if err != nil {
			ExitError(err)
		}
		for _, c := range commits {
			if c.Commit.Message != nil && strings.HasPrefix(*c.Commit.Message, "Merge pull request") {
				continue
			}
			if c.Author == nil || c.Author.Login == nil {
				log.Printf("Author or Author.Login is nil, unexpected commit: %v\n", c.Commit.String())
				continue
			}
			id := *c.Author.Login
			rank[id] = rank[id] + 1
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	cr := committerRank{}
	for id, c := range rank {
		cr = append(cr, &committer{ID: id, CommitCount: c})
	}
	sort.Sort(cr)

	res := []string{}
	for i := 0; i < len(cr); i++ {
		res = append(res, cr[i].ID)
		if len(res) >= limit {
			if i+1 < len(cr) && cr[i+1].CommitCount == cr[i].CommitCount {
				continue
			}
			break
		}
	}
	if len(res) > 0 {
		fmt.Printf("path: %s, reviewers: %v\n", dir, res)
	}
}

type committer struct {
	ID          string
	CommitCount int
}

type committerRank []*committer

func (s committerRank) Len() int { return len(s) }

func (s committerRank) Less(i, j int) bool {
	return s[i].CommitCount > s[j].CommitCount
}

func (s committerRank) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
