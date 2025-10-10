# 仓库指南

本仓库包含基于 Vite 与 Tailwind CSS 的 Vue 3 单页应用，位于 `front/`；以及静态 HTML 原型，位于 `demo/`。

## 项目结构与模块组织
- `front/src/` 应用源码。入口：`front/src/main.js`；根组件：`front/src/App.vue`。
- 视图位于 `front/src/views/`（如 `Home.vue`、`Profile.vue`、`TaskManager.vue`）。推荐一文件一功能。
- 样式：`front/src/style.css` 使用 Tailwind；配置见 `front/tailwind.config.js` 与 `front/postcss.config.js`。
- 工具配置：`front/vite.config.js`。
- `demo/` 为独立 HTML 演示，与 SPA 改动保持隔离。

## 构建、测试与开发命令
- 安装：`cd front && npm install` — 安装依赖。
- 开发：`npm run dev` — 启动带 HMR 的 Vite 开发服务器。
- 构建：`npm run build` — 产物输出到 `front/dist/`。
- 预览：`npm run preview` — 本地服务已构建应用。

## 代码风格与命名约定
- 语言：ES Modules 与 Vue SFC。
- 缩进：2 空格；建议行宽 100–120 字符。
- 文件名：Vue 组件用 PascalCase（如 `TeamTasks.vue`）；JS 模块用 kebab-case。
- 标识符：变量/函数用 camelCase；组件用 PascalCase；环境键用 UPPER_SNAKE_CASE。
- CSS：优先使用 Tailwind 工具类；共用样式提炼为组件或 CSS layer。

## 测试指南
- 当前未配置正式测试。变更时可引入 Vitest 编写轻量单测，并通过 `npm run dev` / `npm run preview` 手动验证。
- 测试文件命名为 `<file>.spec.{js,ts}`，与源文件同级或放在 `__tests__/` 目录。

## 提交与 PR 指南
- 提交信息：使用祈使语与范围前缀（如 `feat(views): add StudyRoom layout`），相关改动合并为一次提交。
- PR：提供清晰描述、关联的 issue（如 `Closes #123`）、UI 前后截图、复现/测试步骤。
- 发起评审前请确保本地构建通过（`npm run build`）。

## 安全与配置提示
- 机密放在 `.env`；需在客户端可用的变量使用 `VITE_` 前缀。切勿提交机密。
- 避免引入全局可变状态；优先使用 Vue 可组合函数或 props/emit 传递。

## 面向代理的说明
- 作用域：本指南适用于 `front/`；保持 `demo/` 独立。
- 优先一致性：在 `front/src/views/` 中沿用既有模式，再引入新模式。
