package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}
	switch e := event.(type) {
	case *github.PingEvent:
		w.Write([]byte("PONG!"))
	case *github.PushEvent:
		// this is a commit push, do something with it
		fmt.Printf("Head: %v\n", e.GetHead())
		fmt.Printf("Ref: %v\n", e.GetRef())
		fmt.Printf("Sender: %v\n", e.GetSender())
		fmt.Printf("Puhser: %v\n", e.GetPusher())
	case *github.PullRequestEvent:
		// this is a pull request, do something with it
		fmt.Printf("Action: %v\n", *e.Action)
		fmt.Printf("Number: %v\n", *e.Number)
		fmt.Printf("PR ID: %v\n", *e.PullRequest.ID)
		fmt.Printf("PR No: %v\n", *e.PullRequest.Number)
	case *github.WatchEvent:
		// https://developer.github.com/v3/activity/events/types/#watchevent
		// someone starred our repository
		fmt.Printf("Event: %v\n", e)
		if e.Action != nil && *e.Action == "starred" {
			fmt.Printf("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}
	default:
		log.Printf("skipping event type %s\n", github.WebHookType(r))
		return
	}
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/webhook", ServeHTTP)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
