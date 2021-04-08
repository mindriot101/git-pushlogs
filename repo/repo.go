package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repo interface {
	CommitObject(plumbing.Hash) (*object.Commit, error)
	Worktree() (*git.Worktree, error)
}
