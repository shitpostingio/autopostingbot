package entities

import (
	"reflect"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite" // We need SQLite to perform migrations
)

func TestPostMayBeRelatedToCategories(t *testing.T) {
	t.Log("Check if a single Post has a collection called Categories with Post")
	post := new(Post)
	kind := reflect.TypeOf(post.Categories).Kind()

	if kind != reflect.Slice {
		t.Errorf("Expected slice, found %s instead.", kind)
	}
}
