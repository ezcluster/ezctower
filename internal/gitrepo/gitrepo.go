package gitrepo

import (
	"errors"
	"ezcluster/tower/internal/config"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-logr/logr"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type GitRepo struct {
	log       logr.Logger
	branch    string
	progress  sideband.Progress
	auth      *http.BasicAuth
	repoPath  string
	repo      *git.Repository
	committer *config.Committer
}

func New(progress sideband.Progress) (*GitRepo, error) {
	gitRepo := &GitRepo{
		log:       config.Log.WithName("gitRepo"),
		branch:    config.Conf.Branch,
		progress:  progress,
		committer: &config.Conf.Committer,
	}
	if config.Conf.Auth.Username != "" && config.Conf.Auth.Token != "" {
		gitRepo.auth = &http.BasicAuth{
			Username: config.Conf.Auth.Username,
			Password: config.Conf.Auth.Token,
		}
	}
	repoName, err := extractRepoName(config.Conf.Repo)
	if err != nil {
		return nil, err
	}
	gitRepo.repoPath = filepath.Join(config.Conf.Workdir, repoName)
	info, err := os.Stat(gitRepo.repoPath)
	if err != nil {
		gitRepo.log.V(1).Info("Clone repo", "url", config.Conf.Repo, "location", gitRepo.repoPath, "branch", config.Conf.Branch)
		// repo does not exists. Must clone
		cloneOptions := &git.CloneOptions{
			URL:           config.Conf.Repo,
			SingleBranch:  true,
			ReferenceName: plumbing.NewBranchReferenceName(config.Conf.Branch),
			Progress:      gitRepo.progress,
			Auth:          gitRepo.auth,
		}
		gitRepo.repo, err = git.PlainClone(gitRepo.repoPath, false, cloneOptions)
		if err != nil {
			return nil, err
		}
	} else {
		if !info.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", gitRepo.repoPath)
		}
		gitRepo.log.V(1).Info("Open repo", "location", gitRepo.repoPath)
		gitRepo.repo, err = git.PlainOpen(gitRepo.repoPath)
		if err != nil {
			return nil, err
		}
	}
	return gitRepo, nil
}

func extractRepoName(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	base := path.Base(u.Path)
	return base[0 : len(base)-4], nil
}

// Pull return a flag true if the pull was effective. false if the local repo was already up to date
func (r *GitRepo) Pull() (bool, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("unable to retrieve Worktree(): %w", err)
	}
	r.log.V(1).Info("pull repo")
	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		Progress:      r.progress,
		Auth:          r.auth,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(config.Conf.Branch),
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return false, err
	}
	return err == nil, nil
}

// AbsPath return the full path from a repo local path
func (r *GitRepo) AbsPath(path string) string {
	return filepath.Join(r.repoPath, path)
}

// Add add the file to the index, to be included in the next commit
func (r *GitRepo) Add(path string) error {
	w, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("unable to retrieve Worktree(): %w", err)
	}
	r.log.V(1).Info("Add file", "file", path)
	_, err = w.Add(path)
	return err
}

func (r *GitRepo) Commit(message string) error {
	w, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("unable to retrieve Worktree(): %w", err)
	}
	r.log.V(1).Info("Commit", "message", message)
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  r.committer.Name,
			Email: r.committer.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *GitRepo) Push() error {
	r.log.V(1).Info("Push")
	return r.repo.Push(&git.PushOptions{
		Auth: r.auth,
	})
}

func (r *GitRepo) GetLastHashLog(path string) (string, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return "", err
	}
	r.log.V(1).Info("log()", "path", path)

	iter, err := r.repo.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
		PathFilter: func(p string) bool {
			return strings.HasPrefix(p, path)
		},
		Since: &time.Time{},
		All:   true,
	})
	if err != nil {
		return "", err
	}
	commit, err := iter.Next()
	if err != nil {
		return "", err
	}
	return commit.Hash.String(), nil
}
