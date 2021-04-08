package pushlogs

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/mindriot101/git-pushlogs/errors"
	"github.com/mindriot101/git-pushlogs/push"
)

type pushlogs struct {
	repo *git.Repository
}

func (p *pushlogs) Pushes() ([]push.Push, error) {
	fn, err := p.logFilename()
	if err != nil {
		return nil, errors.Errorf("error fetching log filename", err)
	}
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.Errorf("cannot find push log file", err)
	}
	defer f.Close()
	pushes := []push.Push{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		push, err := push.NewPush(line, p.repo)
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
