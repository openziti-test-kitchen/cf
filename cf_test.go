package cf

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	basic := &struct {
		StringValue string
	}{}

	var data = map[string]interface{}{
		"string_value": "oh, wow!",
	}

	err := Bind(basic, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "oh, wow!", basic.StringValue)
}

func TestRenaming(t *testing.T) {
	renamed := &struct {
		SomeInt int `cf:"some_int_,+required"`
	}{}

	var data = map[string]interface{}{
		"some_int_": 46,
	}

	err := Bind(renamed, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, 46, renamed.SomeInt)
}

func TestStringArray(t *testing.T) {
	withArray := &struct {
		StringArray []string
	}{}

	var data = map[string]interface{}{
		"string_array": []string{"one", "two", "three"},
	}

	err := Bind(withArray, data, DefaultOptions())
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"one", "two", "three"}, withArray.StringArray)
}

func TestIntArray(t *testing.T) {
	withArray := &struct {
		IntArray []int
	}{}

	var data = map[string]interface{}{
		"int_array": []int{1, 2, 3, 4, 5, 6},
	}

	err := Bind(withArray, data, DefaultOptions())
	assert.Nil(t, err)
}

func TestRequired(t *testing.T) {
	required := &struct {
		Required int `cf:"+required"`
	}{}

	data := make(map[string]interface{})

	err := Bind(required, data, DefaultOptions())
	assert.NotNil(t, err)
}

type nestedType struct {
	Name  string
	Count int
}

func newNestedType() *nestedType {
	return &nestedType{Name: "oh, wow!", Count: 33} // defaults
}

func TestNestedPtr(t *testing.T) {
	root := &struct {
		Id     string
		Nested *nestedType
	}{}

	var data = map[string]interface{}{
		"id": "TestNested",
		"nested": map[string]interface{}{
			"name": "Different",
		},
	}

	opt := DefaultOptions().AddInstantiator(reflect.TypeOf(nestedType{}), func() interface{} { return newNestedType() })

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "TestNested", root.Id)
	assert.NotNil(t, root.Nested)
	assert.Equal(t, "Different", root.Nested.Name)
	assert.Equal(t, 33, root.Nested.Count)
}

func TestNestedValue(t *testing.T) {
	root := &struct {
		Id     string
		Nested nestedType
	}{}

	var data = map[string]interface{}{
		"id": "TestNested",
		"nested": map[string]interface{}{
			"name": "Different",
		},
	}

	opt := DefaultOptions().AddInstantiator(reflect.TypeOf(nestedType{}), func() interface{} { return newNestedType() })

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "TestNested", root.Id)
	assert.NotNil(t, root.Nested)
	assert.Equal(t, "Different", root.Nested.Name)
	assert.Equal(t, 33, root.Nested.Count)
}

func TestNestedWithTypeWiring(t *testing.T) {
	root := &struct {
		Id     string
		Nested *nestedType
	}{}

	var data = map[string]interface{}{
		"id": "TestNested",
		"nested": map[string]interface{}{
			"name": "Different",
		},
	}

	opt := DefaultOptions()
	opt.AddInstantiator(reflect.TypeOf(nestedType{}), func() interface{} { return newNestedType() })
	opt.AddWiring(reflect.TypeOf(nestedType{}), func(cf interface{}) error {
		if v, ok := cf.(*nestedType); ok {
			v.Count = v.Count * 2
		}
		return nil
	})

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "TestNested", root.Id)
	assert.NotNil(t, root.Nested)
	assert.Equal(t, "Different", root.Nested.Name)
	assert.Equal(t, 66, root.Nested.Count) // type wiring
}

func TestNestedWithDefaultInstantiator(t *testing.T) {
	root := &struct {
		Id     string
		Nested *nestedType
	}{}

	var data = map[string]interface{}{
		"id": "TestNested",
		"nested": map[string]interface{}{
			"name": "Different",
		},
	}

	err := Bind(root, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "TestNested", root.Id)
	assert.NotNil(t, root.Nested)
	assert.Equal(t, "Different", root.Nested.Name)
	assert.Equal(t, 0, root.Nested.Count) // type wiring
}

func TestStructTypeArray(t *testing.T) {
	root := &struct {
		Id      string
		Nesteds []*nestedType
	}{}

	var data = map[string]interface{}{
		"id": "StructTypeArray",
		"nesteds": []map[string]interface{}{
			{"name": "a"},
			{"name": "b"},
		},
	}

	err := Bind(root, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "StructTypeArray", root.Id)
	assert.Equal(t, 2, len(root.Nesteds))
	assert.Equal(t, "a", root.Nesteds[0].Name)
	assert.Equal(t, "b", root.Nesteds[1].Name)
}

func TestAnonymousStruct(t *testing.T) {
	root := &struct {
		Id     string
		Nested struct {
			Name string
		}
	}{}

	var data = map[string]interface{}{
		"id": "AnonymousStruct",
		"nested": map[string]interface{}{
			"name": "oh, wow!",
		},
	}

	err := Bind(root, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "AnonymousStruct", root.Id)
	assert.Equal(t, "oh, wow!", root.Nested.Name)
}
