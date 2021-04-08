package push

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type FakeRepo struct {
}

func (f *FakeRepo) CommitObject(h plumbing.Hash) (*object.Commit, error) {
	return &object.Commit{
		Message: "message",
	}, nil
}

func (f *FakeRepo) Worktree() (*git.Worktree, error) {
	return nil, nil
}

func TestPrint(t *testing.T) {
	dt := time.Now()
	push := Push{
		repo: &FakeRepo{},
		t:    dt,
		hash: "hash",
	}
	msg, err := push.Print()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := fmt.Sprintf("%v: hash message\n", dt)
	if msg != expected {
		t.Errorf("%s != %s", msg, expected)
	}
}
