package database

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite" // We need SQLite to perform migrations
)

func run() {
	// Create table for model `User`
	// will append "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci" to the SQL statement when creating table `users`
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").CreateTable(&User{})
}

func drop() {
	// Drop Users table
	db.DropTableIfExists(&User{})
}
