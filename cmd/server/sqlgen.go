package main

import (
	_ "github.com/mattn/go-sqlite3" // SQLite 驱动
	"github.com/rei0721/rei0721/pkg/cli"
)

type SqlGenCommand struct{}

func (c *SqlGenCommand) Name() string {
	return "sqlgen"
}

func (c *SqlGenCommand) Description() string {
	return "SQL 代码生成器演示"
}

func (c *SqlGenCommand) Usage() string {
	return "sqlgen"
}

func (c *SqlGenCommand) Flags() []cli.Flag {
	return []cli.Flag{
		{
			Name:        "type",
			ShortName:   "t",
			Type:        cli.FlagTypeString,
			Required:    false,
			Default:     "sqlite",
			Description: "DB type",
		},
		{
			Name:        "mode",
			ShortName:   "m",
			Type:        cli.FlagTypeString,
			Required:    false,
			Default:     "ddl",
			Description: "DB mode ddl / model",
		},
	}
}

func (c *SqlGenCommand) Execute(ctx *cli.Context) error {
	runSqlGen(ctx)
	return nil
}

func runSqlGen(ctx *cli.Context) {
	dbType := ctx.GetString("type")
	dbMode := ctx.GetString("mode")

	switch dbType {
	case "sqlite":
		if dbMode == "ddl" {
			sqlDDLGenSqlite()
			return
		}
		sqlModelGenSqlite()
	case "mysql":
		if dbMode == "ddl" {
			sqlDDLGenMySQL()
			return
		}
		sqlModelGenMySQL()
	case "postgres":
		// TODO: postgres gen
		println("PostgreSQL support coming soon...")
	default:
		println("❌ 不支持的数据库类型:", dbType)
		println("支持的类型: sqlite, mysql, postgres")
	}
}
