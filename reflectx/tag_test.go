package reflectx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/sincospro/x/reflectx"
)

type Foo struct {
	Field  int `   name:"tagName,default='0'"    json:"tagName,omitempty"   :  `
	Field1 int `   `      // empty tag
	Field2 int `name:"\\` // not meet end tag
	Field3 int `name:abc`
}

func TestTagValueAndFlags(t *testing.T) {
	ft, _ := reflect.ValueOf(Foo{}).Type().FieldByName("Field")

	nameTag, _ := ft.Tag.Lookup("name")
	t.Log("name tag", nameTag)

	jsonTag, _ := ft.Tag.Lookup("json")
	t.Log("json tag", jsonTag)

	key, flags := ParseTagKeyAndFlags(nameTag)
	NewWithT(t).Expect(key).To(Equal("tagName"))
	NewWithT(t).Expect(flags).To(Equal(map[string]struct{}{"default='0'": {}}))

	key, flags = ParseTagKeyAndFlags(jsonTag)
	NewWithT(t).Expect(key).To(Equal("tagName"))
	NewWithT(t).Expect(flags).To(Equal(map[string]struct{}{"omitempty": {}}))
}

func TestParseStructTag(t *testing.T) {
	v := reflect.ValueOf(Foo{})

	f1 := v.Type().Field(0)
	f2 := v.Type().Field(1)
	f3 := v.Type().Field(2)
	f4 := v.Type().Field(3)

	flags := ParseStructTag(string(f1.Tag))
	NewWithT(t).Expect(flags["name"]).To(Equal("tagName,default='0'"))
	NewWithT(t).Expect(flags["json"]).To(Equal("tagName,omitempty"))

	flags = ParseStructTag(string(f2.Tag))
	NewWithT(t).Expect(len(flags)).To(Equal(0))

	flags = ParseStructTag(string(f3.Tag))
	NewWithT(t).Expect(len(flags)).To(Equal(0))

	flags = ParseStructTag(string(f4.Tag))
	NewWithT(t).Expect(len(flags)).To(Equal(0))
}
