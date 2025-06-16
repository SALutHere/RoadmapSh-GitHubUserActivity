package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

// Returns a username from specified argument
func GetUsername() string {
	if len(os.Args) != 2 {
		log.Fatal("Wrong count of arguments. Please, use: github-activity <username>")
	}
	username := os.Args[1]
	return username
}

// Returns a list of public GitHub-user's events by its username
func GetSortedEvents(username string) []Event {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting response: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading of response's body: %v", err)
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		log.Fatalf("Error unmarshaling response's body: %v", err)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})

	return events
}

// Returns grouped consecutive events of the same type. Grouping: by login and by event type
func GroupEvents(events []Event) [][]Event {
	groups := [][]Event{}
	currentGroup := []Event{}
	for _, event := range events {
		if len(currentGroup) == 0 {
			currentGroup = append(currentGroup, event)
			continue
		}
		if event.Type != currentGroup[len(currentGroup)-1].Type ||
			event.Actor.Login != currentGroup[len(currentGroup)-1].Actor.Login {
			groups = append(groups, currentGroup)
			currentGroup = []Event{}
			currentGroup = append(currentGroup, event)
			continue
		}
		currentGroup = append(currentGroup, event)
	}
	groups = append(groups, currentGroup)

	return groups
}

// Returns user-friendly output from groups of events
func RepresentOutput(groups [][]Event) string {
	if len(groups[0]) == 0 {
		return "Specified user has no public activity"
	}

	var output string

	for _, group := range groups {
		switch group[0].Type {
		case "CommitCommentEvent":
			count := len(group)
			login := group[0].Actor.Login
			repo := group[0].Repo.Name
			record := fmt.Sprintf("- %s commented commit %d times (%s)\n", login, count, repo)
			output += record
		case "CreateEvent":
			for _, event := range group {
				login := event.Actor.Login
				refType := event.Payload.RefType
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s created a %s (%s)\n", login, refType, repo)
				output += record
			}
		case "DeleteEvent":
			for _, event := range group {
				login := event.Actor.Login
				refType := event.Payload.RefType
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s deleted a %s (%s)\n", login, refType, repo)
				output += record
			}
		case "ForkEvent":
			for _, event := range group {
				login := event.Actor.Login
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s forked %s\n", login, repo)
				output += record
			}
		case "GollumEvent":
			count := len(group)
			login := group[0].Actor.Login
			repo := group[0].Repo.Name
			record := fmt.Sprintf("- %s updated wiki %d times (%s)\n", login, count, repo)
			output += record
		case "IssueCommentEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s comment on some issue(%s)\n", login, action, repo)
				output += record
			}
		case "IssueEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s an issue (%s)\n", login, action, repo)
				output += record
			}
		case "MemberEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s collaboration (%s)\n", login, action, repo)
				output += record
			}
		case "PublicEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s %s\n", login, action, repo)
				output += record
			}
		case "PullRequestEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s pull request (%s)\n", login, action, repo)
				output += record
			}
		case "PullRequestReviewEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s pull request review (%s)\n", login, action, repo)
				output += record
			}
		case "PullRequestReviewCommentEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s pull request review comment (%s)\n", login, action, repo)
				output += record
			}
		case "PullRequestReviewThreadEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s comment thread on pull request (%s)\n", login, action, repo)
				output += record
			}
		case "PushEvent":
			login := group[0].Actor.Login
			commitsCount := len(group[0].Payload.Commits)
			repo := group[0].Repo.Name
			record := fmt.Sprintf("- %s pushed %d commits (%s)\n", login, commitsCount, repo)
			output += record
		case "ReleaseEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s %s\n", login, action, repo)
				output += record
			}
		case "SponsorshipEvent":
			for _, event := range group {
				login := event.Actor.Login
				action := event.Payload.Action
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s %s sponsorship listing on %s\n", login, action, repo)
				output += record
			}
		case "WatchEvent":
			for _, event := range group {
				login := event.Actor.Login
				repo := event.Repo.Name
				record := fmt.Sprintf("- %s starred %s\n", login, repo)
				output += record
			}
		}
	}

	return output
}

func main() {
	username := GetUsername()
	events := GetSortedEvents(username)
	eventGroups := GroupEvents(events)
	output := RepresentOutput(eventGroups)
	fmt.Println(output)
}
