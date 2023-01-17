package gitrepo

import (
	"errors"
	"ezcluster/tower/internal/config"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-logr/logr"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type GitRepo struct {
	log      logr.Logger
	branch   string
	progress sideband.Progress
	auth     *http.BasicAuth
	repoPath string
	repo     *git.Repository
}

func New(progress sideband.Progress) (*GitRepo, error) {
	gitRepo := &GitRepo{
		log:      config.Log.WithName("gitRepo"),
		branch:   config.Conf.Branch,
		progress: progress,
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
		gitRepo.log.Info("Clone repo", "url", config.Conf.Repo, "location", gitRepo.repoPath, "branch", config.Conf.Branch)
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
		gitRepo.log.Info("Open repo", "location", gitRepo.repoPath)
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
		return false, err
	}
	r.log.Info("pull repo")
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
