# Git Workflow Rules

## 分支命名

```
feat/short-description      # 新功能
fix/issue-123               # Bug 修复（关联 Issue 编号）
refactor/module-name        # 重构
chore/dependency-update     # 非功能性变更
```

## Commit Message 格式

遵循 Conventional Commits 规范：

```
<type>(<scope>): <subject>

[可选 body]

[可选 footer]
```

**type 取值：**
- `feat` — 新功能
- `fix` — Bug 修复
- `refactor` — 重构（不修改行为）
- `test` — 增加或修改测试
- `chore` — 构建、依赖、配置等非业务变更
- `docs` — 文档更新

**规则：**
- subject 用中文，不超过 50 字
- 使用祈使句（"添加"而非"添加了"）
- 不加句号结尾

**示例：**
```
feat(auth): 添加 JWT 刷新 token 接口

fix(cache): 修复 LRUCache 并发读写导致的 panic

chore(deps): 升级 gin 到 v1.10.0
```

## Pull Request

**标题**：与 commit message 格式一致，不超过 70 字符

**Body 模板：**
```markdown
## 改动说明
- 简要描述做了什么，为什么这样做

## 测试方案
- [ ] 单元测试通过：`go test ./...`
- [ ] 手动验证场景：...

## 关联 Issue
Closes #123
```

## 合并策略

- 默认使用 Squash Merge，保持主干历史整洁
- 合并前必须通过 CI 和至少 1 人 review
- 禁止直接 push 到 `main` 分支
