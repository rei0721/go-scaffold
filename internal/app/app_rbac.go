package app

import (
	"fmt"

	"github.com/rei0721/go-scaffold/pkg/rbac"
)

func (a *App) initRBAC() error {
	// 初始化 RBAC
	var err error
	rbacCfg := &rbac.Config{
		DB:          a.DB.DB(),
		ModelPath:   a.Config.RBAC.ModelPath,
		EnableCache: a.Config.RBAC.EnableCache,
		CacheTTL:    a.Config.RBAC.CacheTTL,
		AutoSave:    a.Config.RBAC.AutoSave,
		TablePrefix: a.Config.RBAC.TablePrefix,
	}
	a.RBAC, err = rbac.New(rbacCfg)
	if err != nil {
		return fmt.Errorf("failed to init rbac: %w", err)
	}

	return nil
}
