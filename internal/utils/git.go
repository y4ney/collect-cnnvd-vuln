package utils

import (
	"fmt"
	"golang.org/x/xerrors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"time"
)

type Git struct {
	URL        string
	Dir        string
	RemoteName string
	Email      string
	Name       string
	Token      string
}

// TODO git clean
func (g *Git) Clone() error {
	// 创建目录
	if err := Mkdir(g.Dir); err != nil {
		return xerrors.Errorf("failed to mkdir %s:%w", g.Dir, err)
	}

	// 克隆项目
	_, err := git.PlainClone(g.Dir, false, &git.CloneOptions{URL: g.URL})
	if err == nil || err == git.ErrRepositoryAlreadyExists {
		return nil
	}

	return xerrors.Errorf("failed to clone %s:%w", g.URL, err)
}

func (g *Git) Pull() error {
	// 打开Git仓库
	repo, err := git.PlainOpen(g.Dir)
	if err != nil {
		return xerrors.Errorf("failed to open %s:%w", g.Dir, err)
	}

	// 获取仓库的工作树
	wt, err := repo.Worktree()
	if err != nil {
		return xerrors.Errorf("failed to get work tree:%w", err)
	}

	// 执行Git pull操作
	err = wt.Pull(&git.PullOptions{RemoteName: g.RemoteName})
	if err == nil || err == git.NoErrAlreadyUpToDate {
		return nil
	}

	return xerrors.Errorf("failed to pull %s:%w", g.URL, err)
}

func (g *Git) Push() error {
	// 打开仓库
	repo, err := git.PlainOpen(g.Dir)
	if err != nil {
		return xerrors.Errorf("failed to open repo path:%w", err)
	}

	// 获取默认远程仓库
	remote, err := repo.Remote(g.RemoteName)
	if err != nil {
		return xerrors.Errorf("failed to get remote repo:%w", err)
	}

	// 验证推送
	err = remote.Push(&git.PushOptions{RemoteName: g.RemoteName, Auth: &http.BasicAuth{Username: g.Token}})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return xerrors.Errorf("failed to push repo:%w", err)
	}

	return nil
}

func (g *Git) Add() error {
	repo, err := git.PlainOpen(g.Dir)
	if err != nil {
		return xerrors.Errorf("failed to open repo path:%w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return xerrors.Errorf("failed to get worktree :%w", err)
	}

	_, err = worktree.Add(".")
	if err != nil {
		return xerrors.Errorf("failed to add :%w", err)
	}
	return nil
}

func (g *Git) Commit() error {
	repo, err := git.PlainOpen(g.Dir)
	if err != nil {
		return xerrors.Errorf("failed to open repo path:%w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return xerrors.Errorf("failed to get worktree :%w", err)
	}

	_, err = worktree.Commit(fmt.Sprintf("update at %v", time.Now().Local()), &git.CommitOptions{
		Author: &object.Signature{Name: g.Name, Email: g.Email, When: time.Now().Local()},
	})
	if err != nil {
		return xerrors.Errorf("failed to commit :%w", err)
	}

	return nil
}
