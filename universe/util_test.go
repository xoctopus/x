package universe_test

import (
	"github.com/xoctopus/x/universe"
	"go/ast"
	"go/types"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestChanDir(t *testing.T) {
	cases := []*struct {
		v   any
		dir types.ChanDir
	}{
		{ast.SEND, types.SendOnly},
		{ast.RECV, types.RecvOnly},
		{ast.ChanDir(0), types.SendRecv},
		{reflect.SendDir, types.SendOnly},
		{reflect.RecvDir, types.RecvOnly},
		{reflect.BothDir, types.SendRecv},
		{types.SendOnly, types.SendOnly},
		{types.RecvOnly, types.RecvOnly},
		{types.SendRecv, types.SendRecv},
	}
	for _, c := range cases {
		NewWithT(t).Expect(universe.ChanDir(c.v).Unwrap()).To(Equal(c.dir))
	}
	t.Run("Invalid", func(t *testing.T) {
		t.Run("InvalidType", func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					NewWithT(t).Expect(err.(error).Error()).To(ContainSubstring("invalid dir type"))
				}
			}()
			universe.ChanDir(1)
		})
		t.Run("InvalidReflectDir", func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					NewWithT(t).Expect(err.(error).Error()).To(ContainSubstring("invalid dir type"))
				}
			}()
			universe.ChanDir(reflect.ChanDir(-1))
		})
	})
}
