package push

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mindriot101/git-pushlogs/errors"
	"github.com/mindriot101/git-pushlogs/repo"
)

type Push struct {
	repo repo.Repo
	t    time.Time
	hash string
}

func (p *Push) Print() (string, error) {
	msg, err := p.commitMessage(p.hash)
	if err != nil {
		return "", errors.Errorf(fmt.Sprintf("error fetching message for commit %s", p.hash), err)
	}
	h := commitHeader(msg)
	return fmt.Sprintf("%v: %s %s\n", p.t, p.shortHash(), h), nil
}

func (p *Push) shortHash() string {
	if len(p.hash) >= 16 {
		return p.hash[:16]
	} else {
		return p.hash
	}
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

func New(line string, repo repo.Repo) (Push, error) {
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

func timeFromUnix(s string) (*time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, errors.Errorf(fmt.Sprintf("error parsing %s to an integer", s), err)
	}
	tm := time.Unix(i, 0)
	return &tm, nil
}
