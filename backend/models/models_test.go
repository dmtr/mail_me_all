package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserListSort(t *testing.T) {
	list := UserList{
		TwitterUserSearchResult{TwitterID: "789"}, TwitterUserSearchResult{TwitterID: "123"}, TwitterUserSearchResult{TwitterID: "456"}}
	sort.Sort(list)
	assert.Equal(t, "123", list[0].TwitterID)
	assert.Equal(t, "456", list[1].TwitterID)
	assert.Equal(t, "789", list[2].TwitterID)
}

func TestUserListDiff(t *testing.T) {
	list1 := UserList{
		TwitterUserSearchResult{TwitterID: "789"}, TwitterUserSearchResult{TwitterID: "123"}, TwitterUserSearchResult{TwitterID: "456"}}

	list2 := UserList{
		TwitterUserSearchResult{TwitterID: "789"}, TwitterUserSearchResult{TwitterID: "123"}}

	list3 := UserList{
		TwitterUserSearchResult{TwitterID: "78991"}, TwitterUserSearchResult{TwitterID: "1234"}}

	list4 := UserList{}

	list := list1.Diff(list2)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "456", list[0].TwitterID)

	list = list2.Diff(list1)
	assert.Equal(t, 0, len(list))

	list = list2.Diff(list3)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "789", list[0].TwitterID)
	assert.Equal(t, "123", list[1].TwitterID)

	list = list1.Diff(list4)
	assert.Equal(t, 3, len(list))
	assert.Equal(t, "789", list[0].TwitterID)
	assert.Equal(t, "123", list[1].TwitterID)
	assert.Equal(t, "456", list[2].TwitterID)
}
