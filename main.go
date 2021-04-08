package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/mindriot101/git-pushlogs/errors"
	"github.com/mindriot101/git-pushlogs/pushlogs"
)

func main() {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		log.Fatalf("error opening repo: %v", err)
	}
	p, err := pushlogs.New(repo)
	if err != nil {
		log.Fatalf("error creating pushlogs: %v", err)
	}
	pushes, err := p.Pushes()
	if err != nil {
		log.Fatalf("error fetching pushes: %v", err)
	}
	for _, push := range pushes {
		msg, err := push.Print()
		if err != nil {
			if aerr, ok := err.(*errors.Error); ok {
				if strings.Contains(aerr.Msg, "error fetching message for commit") {
					continue
				}
			}
			log.Printf("error generating log message: %v", err)
			continue
		}
		fmt.Print(msg)
	}
}
