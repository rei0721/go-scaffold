package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rei0721/go-scaffold/pkg/rbac/models"
	"github.com/rei0721/go-scaffold/pkg/rbac/repository"
	"github.com/rei0721/go-scaffold/pkg/rbac/service"
	"github.com/rei0721/go-scaffold/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 这个示例展示如何使用 pkg/rbac 包设置完整的 RBAC 系统
func main() {
	// 1. 初始化数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 2. 自动迁移表结构
	err = db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// 3. 创建 Repository
	repo := repository.NewGormRBACRepository(db)

	// 4. 创建 Service
	rbacService := service.NewRBACService(repo)

	ctx := context.Background()

	// 5. 创建角色
	fmt.Println("=== 创建角色 ===")
	adminRole, err := rbacService.CreateRole(ctx, &types.CreateRoleRequest{
		Name:        "admin",
		Description: "系统管理员",
		Status:      1,
	})
	if err != nil {
		log.Fatal("failed to create admin role:", err)
	}
	fmt.Printf("创建角色: %s (ID: %d)\n", adminRole.Name, adminRole.ID)

	editorRole, err := rbacService.CreateRole(ctx, &types.CreateRoleRequest{
		Name:        "editor",
		Description: "内容编辑者",
		Status:      1,
	})
	if err != nil {
		log.Fatal("failed to create editor role:", err)
	}
	fmt.Printf("创建角色: %s (ID: %d)\n", editorRole.Name, editorRole.ID)

	// 6. 创建权限
	fmt.Println("\n=== 创建权限 ===")
	permissions := []struct {
		name        string
		resource    string
		action      string
		description string
	}{
		{"users:read", "users", "read", "读取用户"},
		{"users:write", "users", "write", "创建/更新用户"},
		{"posts:read", "posts", "read", "读取文章"},
		{"posts:write", "posts", "write", "创建/更新文章"},
		{"posts:delete", "posts", "delete", "删除文章"},
	}

	createdPerms := make([]*models.Permission, 0, len(permissions))
	for _, p := range permissions {
		perm, err := rbacService.CreatePermission(ctx, &types.CreatePermissionRequest{
			Name:        p.name,
			Resource:    p.resource,
			Action:      p.action,
			Description: p.description,
			Status:      1,
		})
		if err != nil {
			log.Fatal("failed to create permission:", err)
		}
		createdPerms = append(createdPerms, perm)
		fmt.Printf("创建权限: %s (%s)\n", perm.Name, perm.Description)
	}

	// 7. 为角色分配权限
	fmt.Println("\n=== 分配权限给角色 ===")

	// admin 拥有所有权限
	fmt.Println("为 admin 角色分配所有权限")
	for _, perm := range createdPerms {
		err = rbacService.AssignPermission(ctx, adminRole.ID, perm.ID)
		if err != nil {
			log.Fatal("failed to assign permission to admin:", err)
		}
	}

	// editor 只能读写文章
	fmt.Println("为 editor 角色分配文章读写权限")
	for _, perm := range createdPerms {
		if perm.Resource == "posts" && (perm.Action == "read" || perm.Action == "write") {
			err = rbacService.AssignPermission(ctx, editorRole.ID, perm.ID)
			if err != nil {
				log.Fatal("failed to assign permission to editor:", err)
			}
		}
	}

	// 8. 为用户分配角色
	fmt.Println("\n=== 分配角色给用户 ===")
	userID1 := int64(1) // 假设的用户 ID
	userID2 := int64(2)

	err = rbacService.AssignRole(ctx, userID1, adminRole.ID)
	if err != nil {
		log.Fatal("failed to assign admin role to user1:", err)
	}
	fmt.Printf("用户 %d 被分配为 admin\n", userID1)

	err = rbacService.AssignRole(ctx, userID2, editorRole.ID)
	if err != nil {
		log.Fatal("failed to assign editor role to user2:", err)
	}
	fmt.Printf("用户 %d 被分配为 editor\n", userID2)

	// 9. 检查权限
	fmt.Println("\n=== 权限检查 ===")

	// 检查 user1 (admin) 的权限
	testPermissionCheck(rbacService, ctx, userID1, "admin", []struct {
		resource string
		action   string
	}{
		{"users", "write"},
		{"posts", "delete"},
	})

	// 检查 user2 (editor) 的权限
	testPermissionCheck(rbacService, ctx, userID2, "editor", []struct {
		resource string
		action   string
	}{
		{"posts", "write"},  // 应该有权限
		{"posts", "delete"}, // 应该没有权限
		{"users", "write"},  // 应该没有权限
	})

	fmt.Println("\n✅ 示例完成！")
}

func testPermissionCheck(svc service.RBACService, ctx context.Context, userID int64, roleName string, checks []struct {
	resource string
	action   string
}) {
	for _, check := range checks {
		hasPermission, err := svc.CheckPermission(ctx, userID, check.resource, check.action)
		if err != nil {
			log.Fatal("failed to check permission:", err)
		}

		status := "❌ 无权限"
		if hasPermission {
			status = "✅ 有权限"
		}
		fmt.Printf("%s - %s:%s %s\n", roleName, check.resource, check.action, status)
	}
}
