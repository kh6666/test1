# Go Code Style Guide

## 命名规范

- 包名使用小写单词，不用下划线或驼峰：`package httputil`
- 导出标识符用 PascalCase，非导出用 camelCase
- 接口名通常以 `-er` 结尾：`Reader`、`Stringer`、`Handler`
- 缩写词保持大小写一致：`URL`、`ID`、`HTTP`，不写成 `Url`、`Id`、`Http`
- 变量名尽量短，作用域越小越短；循环变量用 `i`、`j`，不用 `index`

## 错误处理

- 错误变量命名为 `err`，自定义错误类型以 `Error` 结尾
- 立即处理错误，不跳过：

```go
// 错误 ✗
result, _ := doSomething()

// 正确 ✓
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething: %w", err)
}
```

- 使用 `%w` 包装错误以保留调用链，方便 `errors.Is` / `errors.As` 解包
- 不要在库函数中使用 `log.Fatal` / `panic`，将错误返回给调用方

## 代码格式

- 始终用 `gofmt` / `goimports` 格式化，提交前必须通过
- 每行不超过 120 个字符
- import 分三组，用空行隔开：标准库 → 第三方库 → 内部包

```go
import (
    "context"
    "fmt"

    "github.com/gin-gonic/gin"

    "example.com/myproject/internal/model"
)
```

## 函数与方法

- 函数参数超过 3 个时，改用结构体传参
- 返回值超过 3 个时，改用结构体返回
- 避免裸返回（naked return），始终显式写返回值
- receiver 名用类型名首字母缩写，保持一致：

```go
func (s *Server) Start() error { ... }
func (s *Server) Stop() error  { ... }
```

## 并发

- goroutine 必须有明确的退出路径，用 `context` 传递取消信号
- 共享状态用 `sync.Mutex` 保护，或改用 channel 通信
- 避免在循环中直接启动 goroutine 捕获循环变量（Go 1.22 之前需显式传参）

```go
// Go 1.21 及以下 ✗
for _, v := range items {
    go func() { process(v) }()
}

// 正确 ✓
for _, v := range items {
    v := v
    go func() { process(v) }()
}
```

## 注释

- 导出的类型、函数必须有文档注释，以标识符名称开头：

```go
// Server handles incoming HTTP requests.
type Server struct { ... }

// Start begins listening on the configured address.
func (s *Server) Start() error { ... }
```

- 非导出代码只在 WHY 不明显时才写注释，不解释 WHAT

## 测试

- 测试文件与被测文件同包（白盒）或加 `_test` 后缀（黑盒）
- 测试函数命名：`TestFunctionName_Scenario`
- 用表驱动测试覆盖多个用例：

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

- 使用 `t.Helper()` 在辅助函数中定位真实失败行

## 其他

- 优先用 `var` 声明零值变量，用 `:=` 声明有初始值的变量
- slice / map 初始化时预分配容量：`make([]T, 0, n)`
- 不用 `init()` 做有副作用的初始化，改为显式调用
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
