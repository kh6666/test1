# Go Testing Rules

## 基本原则

- 每个导出函数必须有对应测试，覆盖率目标 ≥ 80%
- 测试必须是幂等的，可以任意顺序、并行执行
- 测试不依赖外部环境（网络、文件系统、时钟），需隔离或 mock
- CI 中所有测试必须通过，禁止合入带有 `t.Skip` 占位的测试

## 文件组织

```
pkg/
├── server.go
├── server_test.go        # 白盒测试（同包）
└── server_external_test.go  # 黑盒测试（package xxx_test）
```

- 单元测试与源文件同目录
- 集成测试放 `tests/integration/`，需 `//go:build integration` 构建标签

## 命名规范

```go
// 函数测试
func TestFunctionName_Scenario(t *testing.T) {}

// 方法测试
func TestTypeName_MethodName_Scenario(t *testing.T) {}

// 基准测试
func BenchmarkFunctionName(b *testing.B) {}

// 模糊测试
func FuzzFunctionName(f *testing.F) {}
```

场景描述用下划线分隔，尽量用自然语言：
`TestParse_EmptyInput`、`TestServer_Start_PortInUse`

## 表驱动测试

所有多用例场景必须用表驱动，禁止复制粘贴多个 `TestXxx` 函数：

```go
func TestDivide(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name    string
        a, b    float64
        want    float64
        wantErr bool
    }{
        {"normal", 10, 2, 5, false},
        {"divide by zero", 10, 0, 0, true},
        {"negative", -6, 3, -2, false},
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            got, err := Divide(tt.a, tt.b)
            if (err != nil) != tt.wantErr {
                t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 断言与错误报告

- 优先用标准库 `t.Errorf` / `t.Fatalf`，不强制引入断言库
- 若引入第三方断言库，统一用 `github.com/stretchr/testify`，禁止混用多个库
- 错误信息必须说明期望值和实际值：

```go
// 错误 ✗
t.Error("wrong result")

// 正确 ✓
t.Errorf("Parse(%q) = %v, want %v", input, got, want)
```

- 辅助函数调用 `t.Helper()` 使错误定位到调用处：

```go
func assertEqual(t *testing.T, got, want int) {
    t.Helper()
    if got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}
```

## Mock 与依赖隔离

- 通过接口注入依赖，不 mock 具体类型
- 使用 `gomock` 生成 mock（`go generate` 驱动），mock 文件放 `mocks/` 目录
- 时间依赖注入 `clock` 接口，禁止在业务代码中直接调用 `time.Now()`
- HTTP 外部调用用 `httptest.NewServer` 拦截，不发真实请求

```go
// 接口定义
type Clock interface {
    Now() time.Time
}

// 测试用 mock
type fixedClock struct{ t time.Time }
func (c fixedClock) Now() time.Time { return c.t }
```

## 数据库测试

- 集成测试使用独立测试数据库，通过环境变量 `TEST_DB_DSN` 注入连接串
- 每个测试用例在事务中执行，结束后回滚：

```go
func TestUserRepo_Create(t *testing.T) {
    db := testDB(t)
    tx, _ := db.Begin()
    t.Cleanup(func() { tx.Rollback() })
    // ...
}
```

- 禁止测试间共享可变数据库状态

## 并行测试

- 纯单元测试默认调用 `t.Parallel()`
- 有共享资源（端口、文件、数据库）的测试不并行
- 使用 `t.Setenv` 修改环境变量（自动还原），不直接调用 `os.Setenv`

## 基准测试

```go
func BenchmarkEncode(b *testing.B) {
    data := generateTestData()
    b.ResetTimer()          // 排除准备时间
    b.ReportAllocs()        // 报告内存分配

    for i := 0; i < b.N; i++ {
        Encode(data)
    }
}
```

- 运行：`go test -bench=. -benchmem -count=5`
- 性能回归需在 PR 中附 benchmark 对比数据

## 常用命令

```bash
go test ./...                          # 运行所有测试
go test -race ./...                    # 开启竞态检测
go test -run TestFoo ./pkg/...         # 只运行匹配的测试
go test -count=1 ./...                 # 禁用缓存强制重跑
go test -cover -coverprofile=c.out ./... && go tool cover -html=c.out  # 覆盖率报告
go test -tags integration ./tests/...  # 运行集成测试
```
