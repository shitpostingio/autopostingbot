package entities

import (
	"reflect"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite" // We need SQLite to perform migrations
)

func TestUserMayBeRelatedToPosts(t *testing.T) {
	t.Log("Check if user can access to every posts he has made")
	user := new(User)
	kind := reflect.TypeOf(user.Posts).Kind()

	if kind != reflect.Slice {
		t.Errorf("Expected slice, found %s instead.", kind)
	}
}
