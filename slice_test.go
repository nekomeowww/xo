package xo

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToMap(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	type foo struct {
		ID   int
		Name string
	}

	inputs := make([]*foo, 2)
	for i := 0; i < 2; i++ {
		inputs[i] = &foo{
			ID:   i,
			Name: "A",
		}
	}

	result := ToMap(inputs, func(item *foo) int { return item.ID })
	for _, elem := range inputs {
		v, ok := result[elem.ID]
		assert.True(ok)
		assert.Equal(elem, v)
	}
}

func TestSliceSlices(t *testing.T) {
	t.Parallel()

	t.Run("Int64", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)
		require := require.New(t)

		var s1 []int64
		length := 2
		result := SliceSlices(s1, length)
		require.Len(result, 1)
		assert.Empty(result[0])

		s1 = make([]int64, 0)
		result = SliceSlices(s1, length)
		require.Len(result, 1)
		assert.Empty(result[0])

		s1 = []int64{1, 2, 3}
		result = SliceSlices(s1, length)
		require.Len(result, 2)
		assert.Equal(result, [][]int64{{1, 2}, {3}})

		s1 = []int64{1, 2, 3, 4}
		result = SliceSlices(s1, length)
		require.Len(result, 2)
		assert.Equal(result, [][]int64{{1, 2}, {3, 4}})

		s1 = make([]int64, 0)
		for i := 0; i < 10000; i++ {
			s1 = append(s1, int64(i+1))
		}

		length = 10
		result = SliceSlices(s1, length)
		require.Len(result, 10000/length)

		for _, v := range result {
			require.Len(v, length)
		}

		s1 = make([]int64, 0)
		for i := 0; i < 10000; i++ {
			s1 = append(s1, int64(i+1))
		}

		length = 100000
		result = SliceSlices(s1, length)

		actualLength := 10000 / length
		if actualLength <= 0 {
			actualLength = 1
		}

		require.Len(result, actualLength)

		for _, v := range result {
			require.Len(v, len(s1))
		}
	})

	t.Run("Struct", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)
		require := require.New(t)

		type testType struct {
			id int64
		}

		var s1 []testType
		length := 2
		result := SliceSlices(s1, length)
		require.Len(result, 1)
		assert.Empty(result[0])

		s1 = make([]testType, 0)
		result = SliceSlices(s1, length)
		require.Len(result, 1)
		assert.Empty(result[0])

		s1 = []testType{
			{id: 1},
			{id: 2},
			{id: 3},
		}
		result = SliceSlices(s1, length)
		require.Len(result, 2)
		assert.Equal(result, [][]testType{
			{{id: 1}, {id: 2}},
			{{id: 3}},
		})

		s1 = []testType{
			{id: 1},
			{id: 2},
			{id: 3},
			{id: 4},
		}
		result = SliceSlices(s1, length)
		require.Len(result, 2)
		assert.Equal(result, [][]testType{
			{{id: 1}, {id: 2}},
			{{id: 3}, {id: 4}},
		})

		s1 = make([]testType, 0)
		for i := 0; i < 10000; i++ {
			s1 = append(s1, testType{id: int64(i + 1)})
		}

		length = 10
		result = SliceSlices(s1, length)
		require.Len(result, 10000/length)

		for _, v := range result {
			require.Len(v, length)
		}
	})
}

func TestJoin(t *testing.T) {
	t.Parallel()

	t.Run("Int64", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		a := []int64{1, 2, 3}
		assert.Equal("1,2,3", Join(a, ","))
	})

	t.Run("Float64", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		a := []float64{1.2, 1.0, 1}
		assert.Equal("1.2,1,1", Join(a, ","))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		a := []error{errors.New("1"), errors.New("2"), errors.New("3")}
		assert.Equal("1,2,3", Join(a, ","))
	})
}

func TestJoinWithConverter(t *testing.T) {
	t.Parallel()

	t.Run("Int64", func(t *testing.T) {
		t.Parallel()

		t.Run("fmt", func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			a := []int64{1, 2, 3}
			assert.Equal("1,2,3", JoinWithConverter(a, ",", func(v int64) string {
				return fmt.Sprintf("%d", v)
			}))
		})

		t.Run("strconv", func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			a := []int64{1, 2, 3}
			assert.Equal("1,2,3", JoinWithConverter(a, ",", func(v int64) string {
				return strconv.FormatInt(v, 10)
			}))
		})
	})

	t.Run("Float64", func(t *testing.T) {
		t.Parallel()

		t.Run("fmt", func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			a := []float64{1.2, 1.0, 1}
			assert.Equal("1.20,1.00,1.00", JoinWithConverter(a, ",", func(v float64) string {
				return fmt.Sprintf("%.2f", v)
			}))
		})

		t.Run("strconv", func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			a := []float64{1.2, 1.0, 1}
			assert.Equal("1.2,1,1", JoinWithConverter(a, ",", func(v float64) string {
				return strconv.FormatFloat(v, 'f', -1, 64)
			}))
		})
	})
}

func TestIntersection(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	t.Run("Number", func(t *testing.T) {
		t.Parallel()

		a1 := []int64{1, 1}
		a2 := []int64{2, 2}
		intersection := Intersection(a1, a2)
		require.Empty(intersection)

		a1 = []int64{1, 2}
		a2 = []int64{2, 3}
		intersection = Intersection(a1, a2)
		require.NotEmpty(intersection)
		assert.ElementsMatch([]int64{2}, Intersection(a1, a2))

		a1 = []int64{1, 2}
		a2 = []int64{3, 4}
		intersection = Intersection(a1, a2)
		require.Empty(intersection)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		a1 := []string{"1", "1"}
		a2 := []string{"2", "2"}
		intersection := Intersection(a1, a2)
		require.Empty(intersection)

		a1 = []string{"1", "2"}
		a2 = []string{"2", "3"}
		intersection = Intersection(a1, a2)
		require.NotEmpty(intersection)
		assert.ElementsMatch([]string{"2"}, Intersection(a1, a2))

		a1 = []string{"1", "2"}
		a2 = []string{"3", "4"}
		intersection = Intersection(a1, a2)
		require.Empty(intersection)
	})

	t.Run("Struct", func(t *testing.T) {
		t.Parallel()

		type testStruct struct {
			id int64
		}

		a1 := []testStruct{{id: 1}, {id: 1}}
		a2 := []testStruct{{id: 2}, {id: 2}}
		intersection := Intersection(a1, a2)
		require.Empty(intersection)

		a1 = []testStruct{{id: 1}, {id: 2}}
		a2 = []testStruct{{id: 2}, {id: 3}}
		intersection = Intersection(a1, a2)
		require.NotEmpty(intersection)
		assert.ElementsMatch([]testStruct{{id: 2}}, Intersection(a1, a2))

		a1 = []testStruct{{id: 1}, {id: 2}}
		a2 = []testStruct{{id: 3}, {id: 4}}
		intersection = Intersection(a1, a2)
		require.Empty(intersection)
	})
}

func TestFindDuplicates(t *testing.T) {
	t.Parallel()

	t.Run("Number", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		list := []int64{1, 2, 3, 4, 5}
		repeatList := FindDuplicates(list)
		assert.Len(repeatList, 0)

		list = []int64{1, 1, 2}
		repeatList = FindDuplicates(list)
		assert.ElementsMatch([]int64{1}, repeatList)

		list = []int64{1, 1, 2, 2}
		repeatList = FindDuplicates(list)
		assert.ElementsMatch([]int64{1, 2}, repeatList)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		list := []string{"1", "2", "3", "4", "5"}
		repeatList := FindDuplicates(list)
		assert.Len(repeatList, 0)

		list = []string{"1", "1", "2"}
		repeatList = FindDuplicates(list)
		assert.ElementsMatch([]string{"1"}, repeatList)

		list = []string{"1", "1", "2", "2"}
		repeatList = FindDuplicates(list)
		assert.ElementsMatch([]string{"1", "2"}, repeatList)
	})
}
