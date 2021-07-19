package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/google/go-github/v37/github"
)

// Code Repo struct
type Package struct {
	Name             string `json:"name"`
	WatcherCount     int    `json:"watchers"`
	StarsCount       int    `json:"starts"`
	ForksCount       int    `json:"forks"`
	IssuesCount      int    `json:"issues"`
	PullRequestCount int    `json:"prs"`
	Fork             bool   `json:"fork"`
	Language         string `json:"language"`
	LastUpdatedBy    string `json:"last_update_time"`
	RepoUrl          string `json:"url"`
	Description      string `json:"description"`
}

func fetchRepoDetails(owner string) ([]*Package, error) {
	client := github.NewClient(nil)
	var githubRepos []*github.Repository
	var repositories []*Package

	listOps := &github.RepositoryListOptions{
		Sort:        "updated",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// TODO: Fetch all the repositories under the owner. handle pagination
	//nextPageToken := ""
	for {
		repos, resp, err := client.Repositories.List(context.Background(), owner, listOps)
		if err != nil {
			return repositories, errors.New(fmt.Sprintf("Problem in getting repository information %v\n", err))
		}

		githubRepos = append(githubRepos, repos...)
		if resp.NextPageToken == "" {
			break
		}
	}

	fmt.Println(fmt.Sprintf("Fetch %d repositories under organization %v: ", len(githubRepos), owner))
	for _, repo := range githubRepos {
		p := &Package{
			Name:          *repo.Name,
			RepoUrl:       *repo.HTMLURL,
			WatcherCount:  *repo.WatchersCount,
			StarsCount:    *repo.StargazersCount,
			ForksCount:    *repo.ForksCount,
			IssuesCount:   *repo.OpenIssuesCount,
			Fork:          *repo.Fork,
			LastUpdatedBy: repo.GetUpdatedAt().Format("02-01-2006"),
		}

		if repo.Description != nil {
			p.Description = *repo.Description
		}

		if repo.Language != nil {
			p.Language = *repo.Language
		}

		repositories = append(repositories, p)
	}

	return repositories, nil
}

func main() {
	var username string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	repositories, err := fetchRepoDetails(username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = writeToCsv(repositories)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

func writeToCsv(repositories []*Package) error {
	repositoriesFile, err := os.OpenFile("devstats.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer repositoriesFile.Close()

	// Go to the start of the file
	if _, err := repositoriesFile.Seek(0, 0); err != nil {
		panic(err)
	}

	// Get all structs as CSV string
	//csvContent, err := gocsv.MarshalString(&repositories)
	err = gocsv.MarshalFile(&repositories, repositoriesFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}

	return nil
}
