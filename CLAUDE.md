# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## My Project

基于 Go + React 的全栈 Web 应用，提供 RESTful API 服务与前端交互界面。

## Tech Stack

- Backend: Go
- Frontend: React + TypeScript
- Database: PostgreSQL

## Common Commands

```bash
go run ./...         # 启动服务
go test ./...        # 运行测试
go build -o app .    # 构建
npm run dev          # 启动前端开发服务器
npm run build        # 构建前端生产版本
```

## Code Conventions

- 导出标识符使用 PascalCase，非导出使用 camelCase
- 错误立即处理，不用 `_` 忽略
- 遵循 `.claude/rules/code-style.md` 中的 Go 规范

## Architecture Overview

```
src/
├── api/        # 路由层
├── services/   # 业务逻辑层
└── models/     # 数据模型层
```
