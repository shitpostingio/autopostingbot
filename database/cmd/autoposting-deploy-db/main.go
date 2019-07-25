package main

import (
	"flag"
	"log"

	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"

	configuration "gitlab.com/shitposting/autoposting-bot/config"
	"gitlab.com/shitposting/autoposting-bot/database/migrations"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	cleanDeploy    bool
	configFilePath string
)

func main() {

	/* PARSE CLI ARGUMENTS */
	flag.StringVar(&configFilePath, "config", "./config.toml", "autoposting-bot configuration file")
	flag.BoolVar(&cleanDeploy, "clean", false, "Drop all tables and create them from scratch")

	flag.Parse()

	/* LOAD CONFIG */
	cfg, err := configuration.Load(configFilePath, false)
	fatalIfErr(err)

	/* CONNECT TO DB */
	db, err := gorm.Open("mysql", cfg.DB.DatabaseConnectionString())
	fatalIfErr(err)

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(db.Close())
		}
	}()

	if !cleanDeploy {
		migrateTables(db)
		return
	}

	dropTables(db)
	createTables(db)
	deployTypes(db)
}

func dropTables(db *gorm.DB) {
	fatalIfErr(migrations.DropFingerprints(db))
	fatalIfErr(migrations.DropPosts(db))
	fatalIfErr(migrations.DropTypes(db))
	fatalIfErr(migrations.DropUsers(db))
}

func createTables(db *gorm.DB) {
	fatalIfErr(migrations.CreateUsers(db))
	fatalIfErr(migrations.CreateTypes(db))
	fatalIfErr(migrations.CreatePosts(db))
	fatalIfErr(migrations.CreateFingerprints(db))
}

func migrateTables(db *gorm.DB) {
	fatalIfErr(migrations.MigrateUsers(db))
	fatalIfErr(migrations.MigrateTypes(db))
	fatalIfErr(migrations.MigratePosts(db))
	fatalIfErr(migrations.MigrateFingerprints(db))
}

func deployTypes(db *gorm.DB) {
	fatalIfErr(db.Create(&entities.Type{Name: "image"}).Error)
	fatalIfErr(db.Create(&entities.Type{Name: "video"}).Error)
	fatalIfErr(db.Create(&entities.Type{Name: "animation"}).Error)
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
