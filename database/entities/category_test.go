package database

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // We need SQLite to perform migrations
)

func TestCategoryMayBeRelatedToPosts(t *testing.T) {
	t.Log("Check if Category has a collection of Post in Posts")
	category := new(Category)
	kind := reflect.TypeOf(category.Posts).Kind()

	if kind != reflect.Slice {
		t.Errorf("Expected slice, found %s instead.", kind)
	}
}

func TestCategoryHasUniqueName(t *testing.T) {
	t.Log("A Category should not be created if another one has the same name")

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic("failed to connect database")
	}

	db.CreateTable(&Category{})
	db.Create(&Category{Name: "foo"})

	if err := db.Create(&Category{Name: "foo"}).Error; err == nil {
		t.Errorf("Expected false, got true when creating another category with the same name")
	}

	defer db.Close()
}
