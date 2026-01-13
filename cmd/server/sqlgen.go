package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite 驱动
	"github.com/rei0721/rei0721/internal/models"
	"github.com/rei0721/rei0721/pkg/cli"
	"github.com/rei0721/rei0721/pkg/sqlgen"
)

type SqlGenCommand struct{}

func (c *SqlGenCommand) Name() string {
	return "sqlgen"
}

func (c *SqlGenCommand) Description() string {
	return "SQL example generator"
}

func (c *SqlGenCommand) Usage() string {
	return "sqlgen [type=<sqlite/mysql/postgres>] [mode=<ddl/model/script>]"
}

func (c *SqlGenCommand) Flags() []cli.Flag {
	return []cli.Flag{
		{
			Name:        "type",
			ShortName:   "t",
			Type:        cli.FlagTypeString,
			Required:    false,
			Default:     "mysql",
			Description: "DB type (sqlite/mysql/postgres)",
		},
		{
			Name:        "mode",
			ShortName:   "m",
			Type:        cli.FlagTypeString,
			Required:    false,
			Default:     "script",
			Description: "Mode: ddl/model/script (script=Model→SQL script)",
		},
	}
}

func (c *SqlGenCommand) Execute(ctx *cli.Context) error {
	runSqlGen(ctx)
	return nil
}

func runSqlGen(ctx *cli.Context) {
	dbType := ctx.GetString("type")

	switch dbType {
	case "sqlite":
		// TODO
		// 初始化生成器
		gen := sqlgen.New(&sqlgen.Config{
			Dialect: sqlgen.MySQL,
			Pretty:  false,
		})

		sql, _ := gen.Table(&models.User{})

		fmt.Println(sql, " --- hhhh")
	case "mysql":
		// TODO
		// 初始化生成器
		gen := sqlgen.New(&sqlgen.Config{
			Dialect: sqlgen.MySQL,
			Pretty:  false,
		})

		sql, _ := gen.Table(&models.User{})

		fmt.Println(sql, " --- hhhh")
	case "postgres":
		// TODO
	default:
		println("❌ 不支持的数据库类型:", dbType)
		println("支持的类型: sqlite, mysql, postgres")
		println()
		println("可用模式:")
		println("  --mode=ddl   从数据库生成 DDL 脚本")
		println("  --mode=model 从数据库生成 Go Model")
		println("  --mode=sql   从 Model 生成 SQL 脚本")
	}
}
