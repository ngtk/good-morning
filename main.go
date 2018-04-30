package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"os"
	"time"
)

func main() {
	orgName := os.Getenv("GOOD_MORNING_GITHUB_ORG_NAME")
	teamSlug := os.Getenv("GOOD_MORNING_GITHUB_TEAM_SLUG")
	accessToken := os.Getenv("GOOD_MORNING_GITHUB_ACCESS_TOKEN")
	client := createGithubClient(accessToken)
	team := findTeam(client, teamSlug, orgName)
	createDiscusstionIfPossible(client, *team.ID)
}

func createDiscusstionIfPossible(client *github.Client, teamID int64) {
	ctx := context.Background()
	discussions, _, err := client.Teams.ListDiscussions(ctx, teamID, nil)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	newDiscussion := newTeamDiscussionForMorningMeeting()

	for _, discussion := range discussions {
		if *discussion.Title == *newDiscussion.Title {
			fmt.Printf("Already created. You can see the discussion at:\n%v", *discussion.HTMLURL)
			os.Exit(0)
		}
	}

	createdDiscussion, _, err := client.Teams.CreateDiscussion(ctx, teamID, newDiscussion)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Created! You can see the discussion at:\n%v", *createdDiscussion.HTMLURL)
}

func newTeamDiscussionForMorningMeeting() github.TeamDiscussion {
	t := time.Now()
	weekday := t.Weekday()
	var dayrange string
	switch weekday {
	case 0, 6: // Weekend
		fmt.Print("How impatient you are! You should do it at the beginning of the next week.")
		os.Exit(1)
	default: // Weekday
		startDay := t
		endDay := t.Add(time.Duration((5-weekday)*24) * time.Hour) // Next occuring friday
		dayrange = fmt.Sprintf("[%v-%v]", startDay.Format("2006.01.02"), endDay.Format("2006.01.02"))
	}

	title := fmt.Sprintf("%v 朝会", dayrange)
	body := `
## Why

- 朝会の代替
  - 各メンバーのタスク、コンディションをお互いに把握する

## What

各自やることを書いていく
  `
	discussion := github.TeamDiscussion{Title: &title, Body: &body}
	return discussion
}

func findTeam(client *github.Client, teamSlug string, orgName string) *github.Team {
	ctx := context.Background()
	teams, _, err := client.Organizations.ListTeams(ctx, orgName, nil)

	if err != nil {
		message := err.(*github.ErrorResponse).Response.Status
		fmt.Print(message)
		os.Exit(1)
	}

	var foundTeam *github.Team
	for _, team := range teams {
		if *team.Slug == teamSlug {
			foundTeam = team
			break
		}
	}
	return foundTeam
}

func createGithubClient(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
