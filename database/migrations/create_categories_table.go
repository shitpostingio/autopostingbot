package database

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
)

func run() {
	// Create table for model `Category`
	// will append "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci" to the SQL statement when creating table `categories`
	db
	.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	.CreateTable(&Category{})

	// Add unique index for name
	db
	.Model(&Post{})
	.AddUniqueIndex("categories_name_unique", "name")
}

func drop() {
	// Drop Category table
	db.DropTableIfExists(&Category{})
}