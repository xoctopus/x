# BDD 测试规范

本文描述 `github.com/xoctopus/x/testx/bdd` 的 Given / When / Then 用法.

## 适用范围

- 需要自然语言结构的行为测试
- 需要将场景拆成可读子测试
- 断言仍基于 `testx` matcher

## 主入口

- `bdd.From(t *testing.T) bdd.T`
- `bdd.Given(setup) func(*testing.T)`: 直接生成标准测试函数

### DSL

- `Given(summary, func(t bdd.T))`: 前置条件
- `When(summary, func(t bdd.T))`: 触发行为
- `Then(summary, checkers ...Checker)`: 期望结果

子测试命名:

- `GIVEN  <summary>`
- `WHEN  <summary>`
- `THEN  <summary>`

## Checker

`Then` 接收 `Checker`. bdd 提供与 testx 对应的内置 checker:

```go
bdd.Equal(v, 2)
bdd.NotEqual(v, 3)
bdd.HaveLen(s, n)
bdd.Succeed(err)
bdd.Failed(err)
```

也可用:

- `AsChecker(matcher, actual)`
- `AsNegativeChecker(matcher, actual)` / `NegativeChecker(c)`

## 示例

```go
func TestSimple(t *testing.T) {
	bdd.From(t).
		Given("initial v = 1", func(t bdd.T) {
			// 设置初始条件
			v := 1

			t.When("add 1", func(t bdd.T) {
				// 执行触发行为
				v += 1
				// 断言期望结果
				t.Then("should equal 2", bdd.Equal(v, 2))
			})

			t.When("add 2", func(t bdd.T) {
				v += 2
				t.Then("should equal 4", bdd.Equal(v, 4))
				t.Then("should greater than 2", bdd.BeGt(v, 2))
				t.Then("not equal 3", bdd.NotEqual(v, 3))
			})
		})

	v := 0
	bdd.Given(func(t bdd.T) {
		v = 1

		t.Then("the string should equal '1'", bdd.Equal(strconv.Itoa(v), "1"))
		t.Then("assert nothing")
	})(t)
}
```

## 约定

- 一个 `Given` 下可有多个 `When`
- 一个 `When` 下可有多个 `Then`
- `Then` 中可组合多个 `Checker`
- 尽量一个 `Then` 描述一个结果语义
- 需要 `panic` 断言时, 可在 `When` 内 `recover`, 再用某个 `Checker` 校验

## 参考

- `github.com/xoctopus/x/testx/bdd`
- matcher 细节见 [testx.spec.md](testx.spec.md)
