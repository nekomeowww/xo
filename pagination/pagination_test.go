package pagination

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	assert := assert.New(t)

	pagination := New(1, 20, 101)
	assert.EqualValues(0, pagination.Offset())

	pagination = New(6, 20, 101)
	assert.EqualValues(6, pagination.MaxPage)
	assert.EqualValues(100, pagination.Offset())
	assert.True(pagination.Valid())

	pagination = New(7, 20, 101)
	assert.EqualValues(6, pagination.MaxPage)
	assert.False(pagination.Valid())

	pagination = New(1, 20, 0)
	assert.False(pagination.Valid())
}

func TestPaginationLimit(t *testing.T) {
	assert := assert.New(t)

	pa := New(1, 20, 10)
	assert.True(pa.Valid() && pa.Limit() == 10)

	pa = New(4, 30, 100)
	assert.Equal(int64(10), pa.Limit())
	assert.True(pa.Valid() && pa.Limit() == 10)
}

func TestLogically(t *testing.T) {
	t.Run("DefaultInt64", func(t *testing.T) {
		assert := assert.New(t)

		list := make([]int64, 100)
		for i := range list {
			list[i] = int64(i)
		}

		slicedList := Logically(list, 1, 10)
		assert.Len(slicedList, 10)
		assert.ElementsMatch(list[:10], slicedList)
	})

	t.Run("DefaultString", func(t *testing.T) {
		assert := assert.New(t)

		list := make([]string, 100)
		for i := range list {
			list[i] = fmt.Sprintf("%d", i)
		}

		slicedList := Logically(list, 1, 10)
		assert.Len(slicedList, 10)
		assert.ElementsMatch(list[:10], slicedList)
	})

	t.Run("DefaultStruct", func(t *testing.T) {
		assert := assert.New(t)

		type T struct {
			ID string
		}

		list := make([]*T, 100)

		for i := range list {
			list[i] = &T{
				ID: fmt.Sprintf("%d", i),
			}
		}

		slicedList := Logically(list, 1, 10)
		assert.Len(slicedList, 10)
		assert.ElementsMatch(list[:10], slicedList)
	})

	t.Run("InvalidPagination", func(t *testing.T) {
		assert := assert.New(t)

		slicedList := Logically([]int64{}, 0, 10)
		assert.Len(slicedList, 0)

		slicedList = Logically([]int64{}, 1, 10)
		assert.Len(slicedList, 0)
	})
}
