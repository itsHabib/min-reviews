package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/itsHabib/min-reviews/internal/github"
	"github.com/itsHabib/min-reviews/internal/setcover"
)

var (
	repo     string
	prNumber int
	exclude  string
)

func main() {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Fatal("GH_TOKEN env var is required")
	}

	flag.Parse()
	if repo == "" {
		log.Fatal("repo flag is required")
	}
	if prNumber == -1 {
		log.Fatal("pr flag is required")
	}

	exclusionMap := make(map[string]struct{})
	if exclude != "" {
		excludeArr := strings.Split(exclude, ",")
		for i := range excludeArr {
			if excludeArr[i] == "" {
				continue
			}
			exclusionMap[excludeArr[i]] = struct{}{}
		}
	}

	c, err := github.NewClient(token)
	if err != nil {
		log.Fatal(err)
	}

	pr, err := c.GetPR(repo, prNumber)
	if err != nil {
		log.Fatalf("unable to get pr #%d: %v", prNumber, err)
	}

	userTeamMap := make(GithubUniverse)
	var users []string
	for i := range pr.RequestedTeams {
		team := pr.RequestedTeams[i]
		members, err := c.GetTeamMembers(team.Url)
		if err != nil {
			log.Fatalf("unable to get members: %s, %v", team.Name, err)
		}
		var userMembers []string
		for j := range members {
			if _, ok := exclusionMap[members[j].Login]; ok {
				continue
			}

			if !contains(users, members[j].Login) {
				userMembers = append(userMembers, members[j].Login)
			}
		}

		userTeamMap[team.Name] = userMembers
		users = append(users, userMembers...)
	}

	solver, err := setcover.NewSolver(users, userTeamMap)
	if err != nil {
		log.Fatal(err)
	}

	solutions := solver.MinCover()
	log.Println("num solutions: ", len(solutions))
	log.Println("solutions: ", solutions)
}

func init() {
	flag.StringVar(&repo, "repo", "", "github repo to get pr from")
	flag.StringVar(&exclude, "exclude", "", "comma separated list of user names to exclude")
	flag.IntVar(&prNumber, "pr", -1, "pr number to get")
}

type GithubUniverse map[string][]string

func (u GithubUniverse) Covers(users []string) bool {
	if u == nil {
		return false
	}

	teams := make(map[string]struct{})
	for k := range u {
		teams[k] = struct{}{}
	}

	for i := range users {
		for k := range teams {
			if contains(u[k], users[i]) {
				delete(teams, k)
				break
			}
		}
	}

	return len(teams) == 0
}

func contains[T comparable](slice []T, item T) bool {
	for i := range slice {
		if slice[i] == item {
			return true
		}
	}

	return false
}
