package reflectx_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/xoctopus/x/reflectx"
)

func TestTagValueAndFlags(t *testing.T) {
	tag := reflect.StructTag(`   name:"tagName,default='0'"    json:"tagName,omitempty"   :  `)
	nameTag, _ := tag.Lookup("name")
	t.Log("name tag", nameTag)

	jsonTag, _ := tag.Lookup("json")
	t.Log("json tag", jsonTag)

	key, flags := ParseTagKeyAndFlags(nameTag)
	NewWithT(t).Expect(key).To(Equal("tagName"))
	NewWithT(t).Expect(flags).To(Equal(map[string]struct{}{"default='0'": {}}))

	key, flags = ParseTagKeyAndFlags(jsonTag)
	NewWithT(t).Expect(key).To(Equal("tagName"))
	NewWithT(t).Expect(flags).To(Equal(map[string]struct{}{"omitempty": {}}))

	tag = `name:",default"`
	key, flags = ParseTagKeyAndFlags(tag.Get("name"))
	NewWithT(t).Expect(key).To(Equal(""))
	NewWithT(t).Expect(flags).To(Equal(map[string]struct{}{"default": {}}))
}

func TestParseStructTag(t *testing.T) {
	flags := ParseStructTag(`   name:"tagName,default='0'"    json:"tagName,omitempty"   :  `)
	NewWithT(t).Expect(flags["name"]).To(Equal("tagName,default='0'"))
	NewWithT(t).Expect(flags["json"]).To(Equal("tagName,omitempty"))

	flags = ParseStructTag(`  `)
	NewWithT(t).Expect(len(flags)).To(Equal(0))

	flags = ParseStructTag(`name:"\\`)
	NewWithT(t).Expect(len(flags)).To(Equal(0))

	flags = ParseStructTag(`name:abc`)
	NewWithT(t).Expect(len(flags)).To(Equal(0))
}
