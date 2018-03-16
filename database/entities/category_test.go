package database

import (
	"reflect"
	"testing"

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
