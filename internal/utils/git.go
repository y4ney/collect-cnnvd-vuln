package utils

import (
	"golang.org/x/xerrors"
	"gopkg.in/src-d/go-git.v4"
)

type Git struct {
	URL        string
	Dir        string
	RemoteName string
}

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
