# testx 断言规范

本文描述 `github.com/xoctopus/x/testx` 的核心断言与 matcher 用法.

## 适用范围

- 标准 `testing` 包单元测试
- 需要类型安全 matcher 组合断言
- 需要校验 panic

## 主入口

### 断言接口

- `Expect[A](t, actual A, matchers ...Matcher[A])`
  - 对 `actual` 依次应用全部 matcher
  - 任一失败则 fail test
- `ExpectPanic[A](t, f func(), matchers ...Matcher[A])`
  - 要求 `f` panic
  - panic 值必须可断言为类型 `A`
  - 若提供 matcher, 继续校验 panic 值
  - 无 matcher 时只要求 panic 且值非 nil

### 自定义 Matcher

- `NewMatcher[Actual](name, func(Actual) bool)`
- `NewComparedMatcher[Actual, Expect](name, func(Actual, Expect) bool)`
- `Not(matcher)` 取反

## Matcher 分类

### 相等与比较

| Matcher                             | 说明           |
|-------------------------------------|----------------|
| `Equal(expect)`                     | 深度相等       |
| `NotEqual(expect)`                  | 深度不等       |
| `Be(expect)` / `NotBe(expect)`      | 同一值语义比较 |
| `BeGt` / `BeGte` / `BeLt` / `BeLte` | 有序比较       |
| `IsZero` / `IsNotZero`              | 零值判断       |
| `BeNil` / `NotBeNil`                | nil 判断       |
| `BeTrue` / `BeFalse`                | bool           |

### 集合

| Matcher                   | 说明                |
|---------------------------|---------------------|
| `HaveLen(n)`              | len                 |
| `HaveCap(n)`              | cap                 |
| `HaveKey(k)`              | map 含 key          |
| `Contains(v)`             | slice 含元素        |
| `EquivalentSlice(expect)` | 元素等价 (顺序无关) |
| `ConsistOfSlice(expect)`  | 元素组成一致        |

### 字符串

| Matcher                     | 说明   |
|-----------------------------|--------|
| `HavePrefix` / `HaveSuffix` | 前后缀 |
| `ContainsSubString`         | 子串   |
| `MatchRegexp`               | 正则   |

### 类型

| Matcher                | 说明       |
|------------------------|------------|
| `IsType[T]()`          | 精确类型   |
| `BeAssignableTo[T]()`  | 可赋值给 T |
| `BeConvertibleTo[T]()` | 可转换到 T |

### 错误

| Matcher                   | 说明              |
|---------------------------|-------------------|
| `Succeed()`               | err == nil        |
| `Failed()`                | err != nil        |
| `IsError(expect)`         | errors.Is         |
| `AsError(target)`         | errors.As         |
| `AsErrorType[T]()`        | 可 As 到类型 T    |
| `ErrorEqual(msg)`         | 错误消息全等      |
| `ErrorContains(sub)`      | 错误消息包含      |
| `IsCodeError[Code](code)` | codex.Code 错误码 |

## 示例

```go
func TestExample(t *testing.T) {
	testx.Expect(t, []int{1, 2}, testx.HaveLen(2), testx.Contains(1))
	testx.Expect(t, "hello", testx.HavePrefix("he"))
	testx.Expect(t, err, testx.Succeed())

	testx.ExpectPanic[string](t, func() {
		panic("boom")
	}, testx.Equal("boom"))
	
	panicked := fmt.Errorf("boom")
	testx.ExpectPanic[error](t, func() {
		panic(panicked)
	}, testx.IsError(panicked))
}
```

## 与 BDD 的关系

`github.com/xoctopus/x/testx/bdd` 将上述 matcher 包装为 `Checker`, 供 `Then` 使用. 详见 [bdd.spec.md](bdd.spec.md).
