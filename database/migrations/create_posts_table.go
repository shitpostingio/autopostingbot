package database

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
)

func run() {
	// Create table for model `Post`
	// will append "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci" to the SQL statement when creating table `posts`
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").CreateTable(&Post{})

	// Add Foreign key to reference the id on users table with cascade onupdate
	db.Model(&Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "CASCADE")
}

func drop() {
	// Drop Post table
	db.DropTableIfExists(&Post{})
}
