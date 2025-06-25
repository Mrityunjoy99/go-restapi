package migration

import (
	"fmt"

	"github.com/Mrityunjoy99/sample-go/cmd/cmdopts"
	"github.com/Mrityunjoy99/sample-go/src/common/config"
	"github.com/Mrityunjoy99/sample-go/src/infrastructure/database"
	"github.com/go-gormigrate/gormigrate/v2"
)

func GetMigrations() []*gormigrate.Migration {
	migrations := []*gormigrate.Migration{
		V20240229232650,
		V20250529214848,
	}

	return migrations
}

func Migrate(a *cmdopts.Arguments) {
	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.Connect(c)
	if err != nil {
		panic(err)
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, GetMigrations())

	switch {
	case cmdopts.IsFlagPassed(cmdopts.MigrateAllOpts.ToString()):
		err = m.Migrate()
	case cmdopts.IsFlagPassed(cmdopts.MigrateToOpts.ToString()):
		version := a.MigrateTo
		err = m.MigrateTo(version)
	case cmdopts.IsFlagPassed(cmdopts.MigrateLastOpts.ToString()):
		err = m.RollbackLast()
	case cmdopts.IsFlagPassed(cmdopts.RollbackToOpts.ToString()):
		version := a.RollbackTo
		err = m.RollbackTo(version)
	default:
		panic("Migration argument is a mandatory field.")
	}

	if err != nil {
		fmt.Printf("\nFailed to migrate: %v", err)
	}

	fmt.Println("Migation successful")
}
