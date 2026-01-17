package main

import (
	"fmt"
	"log"

	"github.com/rei0721/go-scaffold/pkg/rbac"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. 创建临时SQLite数据库（用于演示）
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 2. 创建RBAC实例
	r, err := rbac.New(&rbac.Config{
		DB:          db,
		EnableCache: true,
	})
	if err != nil {
		log.Fatal("failed to create rbac:", err)
	}
	defer r.Close()

	fmt.Println("=== RBAC基本示例 ===\n")

	// 3. 定义策略（角色权限）
	fmt.Println("1. 添加策略：")
	r.AddPolicy("admin", "users", "write")
	r.AddPolicy("admin", "posts", "write")
	r.AddPolicy("editor", "posts", "write")
	r.AddPolicy("viewer", "posts", "read")
	fmt.Println("  - admin 可以 write users")
	fmt.Println("  - admin 可以 write posts")
	fmt.Println("  - editor 可以 write posts")
	fmt.Println("  - viewer 可以 read posts")

	// 4. 分配角色
	fmt.Println("\n2. 分配角色：")
	r.AddRoleForUser("alice", "admin")
	r.AddRoleForUser("bob", "editor")
	r.AddRoleForUser("charlie", "viewer")
	fmt.Println("  - alice -> admin")
	fmt.Println("  - bob -> editor")
	fmt.Println("  - charlie -> viewer")

	// 5. 权限检查
	fmt.Println("\n3. 权限检查：")
	testEnforce(r, "alice", "users", "write")   // true
	testEnforce(r, "alice", "posts", "write")   // true
	testEnforce(r, "bob", "users", "write")     // false
	testEnforce(r, "bob", "posts", "write")     // true
	testEnforce(r, "charlie", "posts", "read")  // true
	testEnforce(r, "charlie", "posts", "write") // false

	// 6. 角色继承示例
	fmt.Println("\n=== 角色继承示例 ===\n")
	r.AddRoleForUser("super_admin", "admin") // super_admin继承admin
	r.AddRoleForUser("alice", "super_admin")
	fmt.Println("alice被分配为super_admin（继承admin）")
	testEnforce(r, "alice", "users", "write") // true（继承自admin）

	// 7. 多租户（域）示例
	fmt.Println("\n=== 多租户（域）示例 ===\n")
	r.AddRoleForUserInDomain("alice", "admin", "tenant1")
	r.AddRoleForUserInDomain("bob", "admin", "tenant2")
	r.AddPolicyWithDomain("admin", "tenant1", "data", "read")
	r.AddPolicyWithDomain("admin", "tenant2", "data", "read")

	fmt.Println("alice在tenant1中是admin")
	fmt.Println("bob在tenant2中是admin")

	ok, _ := r.EnforceWithDomain("alice", "tenant1", "data", "read")
	fmt.Printf("alice能访问tenant1的data吗？%v\n", ok) // true

	ok, _ = r.EnforceWithDomain("alice", "tenant2", "data", "read")
	fmt.Printf("alice能访问tenant2的data吗？%v\n", ok) // false

	// 8. 批量操作示例
	fmt.Println("\n=== 批量操作示例 ===\n")
	rules := [][]string{
		{"manager", "reports", "read"},
		{"manager", "reports", "write"},
		{"manager", "reports", "delete"},
	}
	r.AddPolicies(rules)
	fmt.Println("批量添加了3条策略给manager角色")

	// 9. 查询示例
	fmt.Println("\n=== 查询示例 ===\n")
	roles, _ := r.GetRolesForUser("alice")
	fmt.Printf("alice的角色：%v\n", roles)

	users, _ := r.GetUsersForRole("admin")
	fmt.Printf("拥有admin角色的用户：%v\n", users)

	policies := r.GetPolicy()
	fmt.Printf("所有策略数量：%d\n", len(policies))
}

func testEnforce(r rbac.RBAC, sub, obj, act string) {
	ok, err := r.Enforce(sub, obj, act)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	result := "❌ 拒绝"
	if ok {
		result = "✅ 允许"
	}
	fmt.Printf("  %s %s %s: %s\n", sub, act, obj, result)
}
