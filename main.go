package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mindriot101/git-pushlogs/errors"
)

func timeFromUnix(s string) (*time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, errors.Errorf(fmt.Sprintf("error parsing %s to an integer", s), err)
	}
	tm := time.Unix(i, 0)
	return &tm, nil
}

type Push struct {
	repo *git.Repository
	t    time.Time
	hash string
}

func (p *Push) Print() (string, error) {
	msg, err := p.commitMessage(p.hash)
	if err != nil {
		return "", errors.Errorf(fmt.Sprintf("error fetching message for commit %s", p.hash), err)
	}
	h := commitHeader(msg)
	return fmt.Sprintf("%v: %s %s\n", p.t, p.hash[:16], h), nil
}

func commitHeader(msg string) string {
	lines := strings.Split(msg, "\n")
	return lines[0]
}

func (p *Push) commitMessage(hash string) (string, error) {
	h := plumbing.NewHash(hash)
	c, err := p.repo.CommitObject(h)
	if err != nil {
		return "", errors.Errorf(fmt.Sprintf("invalid hash for project %s", hash), err)
	}
	return c.Message, nil
}

func NewPush(line string, repo *git.Repository) (Push, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return Push{}, errors.Errorf(fmt.Sprintf("error parsing line %s, incorrect number of parts", line), nil)
	}
	dt, err := timeFromUnix(parts[0])
	if err != nil {
		return Push{}, errors.Errorf(fmt.Sprintf("cannot understand time %s", parts[0]), err)
	}
	return Push{
		repo: repo,
		t:    *dt,
		hash: parts[1],
	}, nil
}

type pushlogs struct {
	repo *git.Repository
}

func (p *pushlogs) Pushes() ([]Push, error) {
	fn, err := p.logFilename()
	if err != nil {
		return nil, errors.Errorf("error fetching log filename", err)
	}
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.Errorf("cannot find push log file", err)
	}
	defer f.Close()
	pushes := []Push{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		push, err := NewPush(line, p.repo)
		if err != nil {
			return nil, errors.Errorf(fmt.Sprintf("error reading line %s", line), err)
		}
		pushes = append(pushes, push)
	}
	if err := s.Err(); err != nil {
		return nil, errors.Errorf("error reading log file", err)
	}

	return pushes, nil
}

func (p *pushlogs) logFilename() (string, error) {
	wt, err := p.repo.Worktree()
	if err != nil {
		return "", errors.Errorf("error fetching work tree", err)
	}
	root := wt.Filesystem.Root()
	logFilename := path.Join(root, ".git", "push-log")
	return logFilename, nil
}

func New(repo *git.Repository) (*pushlogs, error) {
	return &pushlogs{
		repo: repo,
	}, nil
}

func main() {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		log.Fatalf("error opening repo: %v", err)
	}
	p, err := New(repo)
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
