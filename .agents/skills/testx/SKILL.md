---
name: testx-guideline
description: 封装 testx 断言与 BDD 测试约定; 当任务涉及编写单元测试, 断言匹配器, panic 断言或 Given/When/Then 风格测试时使用.
triggers:
  - "为...添加单元测试"
  - "为...添加bdd风格单元测试"
  - "为...的未导出方法添加单元测试"
---

# testx-guideline

- 按 `github.com/xoctopus/x/testx` 约定编写类型安全的单元测试断言
- 使用 `testx/bdd` 编写 BDD 风格测试

## 何时使用

- 需要替代手写 `if ... { t.Fatal(...) }` 的断言
- 需要断言函数 panic 及 panic 值
- 需要 Given / When / Then 结构化测试
- 需要自定义 Matcher / Checker

## 基础断言

```go
import (
	"testing"

	"github.com/xoctopus/x/testx"
)

func TestMath(t *testing.T) {
	result := 1 + 1
	testx.Expect(t, result, testx.Equal(2), testx.BeGt(0))

	testx.ExpectPanic(t, func() {
		panic("something went wrong")
	}, testx.Equal("something went wrong"))
}
```

**关键约定**:

- `Expect(t, actual, matchers...)`: actual 需同时满足全部 matcher
- `ExpectPanic[A](t, f, matchers...)`: `f` 必须 panic, 且 panic 值类型为 `A`; 可继续用 matcher 校验 panic 值
- Matcher 是泛型的, 优先用内置 matcher, 避免反射式断言

## BDD 风格

```go
import (
	"testing"

	"github.com/xoctopus/x/testx/bdd"
)

func TestCalculator(t *testing.T) {
	bdd.From(t).
		Given("initial v = 1", func(t bdd.T) {
			v := 1
			t.When("add 1", func(t bdd.T) {
				v += 1
				t.Then("should equal 2", bdd.Equal(2, v))
			})
			t.When("add 2", func(t bdd.T) {
				v += 2
				t.Then("should equal 4", bdd.Equal(4, v))
				t.Then("not equal 3", bdd.NotEqual(3, v))
			})
		})
}
```

**关键约定**:

- `bdd.From(t)` 包装 `*testing.T`
- `Given` / `When` / `Then` 会映射为子测试: `GIVEN ...` / `WHEN ...` / `THEN ...`
- `Then` 使用 `Checker`; bdd checker 基于 testx matcher
- `bdd.T` 底层是 `testing.TB`, 可直接用 `testing.TB` 的所有方法

## Matcher 速览

- 相等/比较: `Equal`, `Be`, `BeGt`, `BeLt`, `IsZero`, `Not`, ...
- 集合: `HaveLen`, `HaveCap`, `HaveKey`, `Contains`, `EquivalentSlice`, `ConsistOfSlice`
- 字符串: `HavePrefix`, `HaveSuffix`, `ContainsSubString`, `MatchRegexp`
- 类型: `BeAssignableTo`, `BeConvertibleTo`, `IsType`
- 错误: `Succeed`, `Failed`, `IsError`, `AsError`, `ErrorEqual`, `ErrorContains`, `IsCodeError`

自定义 matcher:

```go
func BeEven() testx.Matcher[int] {
	return testx.NewMatcher("BeEven", func(actual int) bool {
		return actual%2 == 0
	})
}
```

## 实现参考

- `github.com/xoctopus/x/testx`
- `github.com/xoctopus/x/testx/bdd`

## 单元测试编写范式

- 首先去看需要编写单元测试的方法或函数, 然后编写单元测试.
- 一个方法或函数, 对应一个单元测试. `Given` 对应分支条件, `When` 对应步骤.
- 单元测试是应该是方法或函数的镜像, 默认使用 `bdd` 风格, `Expect` / `ExpectPanic` 做单点断言.
- 单元测试尽量使用黑盒方案, 除非提示需要给未导出方法添加单元测试. 比如 为包 `module/xxx.go` 添加单元测试, 那么在 `module/xxx_test.go` 添加测试内容.
- 黑盒单元测试包名为 `<module_name>_test` 文件名 `<filename>_test.go`
- 白盒单元测试包名为 `<module_name>` 文件名 `<filename>_internal_test.go` (如果没有则新建)
- 如果单元测试需要构造某些`初始条件`, 尽量把他们作为独立的方法, 定义在黑盒测试包中. 方便复用.
- `bdd` 参数顺序一般是 `(actual, expect)`. 比如 `bdd.BeGt(x, y)` 那么 如果 `x>y` 断言通过.
- `Expect` 参数 `Expect(t, actual, Matcher(expect))`

### Given/When的示例

```go
package mod

func X() {
	if condA {
		// condA
		doStep1
		doStep2
	} else {
		// condB
		doStep3
	}
}
```

```go
package mod_test

import "github.com/xoctopus/x/testx/bdd"

func TestX(t *testing.T) {
	bdd.From(t).Given("condA description", func(t bdd.T) {
		// setup condA
		// ...
		t.When("step1 description", func(t bdd.T) {
			// do and assertions
		})
		t.When("step2 description", func(t bdd.T) {
			// do and assertions
		})
	})

	bdd.From(t).Given("condB description", func(t bdd.T) {
		// setup condB
		t.When("step3 description", func(t bdd.T) {
			// do and assertions
		})
	})
}

```

## 更多信息

- 断言与 matcher: [testx.spec.md](references/testx.spec.md)
- BDD 用法: [bdd.spec.md](references/bdd.spec.md)
