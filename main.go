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

	// TODO: ignore list, e.g. "*generated*", "*proto*".
	fetchOwners(client, topLevelDir)
}

func fetchOwners(client *github.Client, dir string) {
	//log.Println("Collecting owners for path:", dir)
	_, directoryContent, _, err := client.Repositories.GetContents(githubOrg, githubRepo, dir, &github.RepositoryContentGetOptions{})
	if err != nil {
		ExitError(err)
	}
	fetchTopCommitters(client, dir, 3)
	for _, c := range directoryContent {
		if c.Type != nil && *c.Type == "dir" {
			fetchOwners(client, *c.Path)
		}
	}
}

func fetchTopCommitters(client *github.Client, dir string, limit int) {
	opt := &github.CommitsListOptions{
		Path:  dir,
		Since: now.AddDate(0, -6, 0),
		ListOptions: github.ListOptions{
			PerPage: 200,
		},
	}
	rank := map[string]int{}
	for {
		commits, resp, err := client.Repositories.ListCommits(githubOrg, githubRepo, opt)
		if err != nil {
			ExitError(err)
		}
		for _, c := range commits {
			if c.Commit.Message == nil {
				log.Printf("Commit.Message is nil, unexpected commit: %v\n", c.Commit.String())
				continue
			}
			if strings.HasPrefix(*c.Commit.Message, "Merge pull request") {
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
	for i := 0; i < limit && i < len(cr); i++ {
		res = append(res, cr[i].ID)
	}
	fmt.Printf("path: %s, owners: %v\n", dir, res)
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
