# CCHLine 项目架构文档

## 项目概述

CCHLine 是一个 Claude Code 状态栏增强工具，从 stdin 读取 JSON 数据，生成带 ANSI 颜色的状态栏输出到 stdout。

> 本项目参考了 [CCometixLine](https://github.com/Haleclipse/CCometixLine) 的设计和代码结构，使用 Go 语言进行了重写。

## CCH 配置

CCHLine 支持连接 CCH 服务，配置方式：

1. **TUI 配置**：`cchline -c`，在底部 CCH SETTINGS 区域设置
2. **命令行**：`cchline -u "URL" -k "API_KEY"`
3. **配置文件**：`~/.claude/cchline/config.toml` 中设置 `cch_url` 和 `cch_api_key`

## 核心架构

### 数据流

1. **输入阶段**：`main.go` 从 stdin 读取 JSON，解析为 `config.InputData` 结构体
2. **收集阶段**：`main.go` 中的 `collectAllSegments()` 遍历所有启用的 Segment，调用各自的 `Collect()` 方法，返回 `[]segment.SegmentResult`
3. **渲染阶段**：`render/statusline.go` 中的 `StatusLineGenerator.Generate()` 遍历结果，根据主题配置应用图标和颜色，用分隔符连接
4. **输出阶段**：将带 ANSI 颜色码的字符串输出到 stdout

### 模块职责

| 文件 | 职责 |
|------|------|
| `main.go` | 程序入口，命令行参数解析（`-c` 进入配置模式），`collectAllSegments()` 协调各 Segment |
| `config/config.go` | 定义 `InputData`（stdin JSON 结构）和 `SimpleConfig`（TOML 配置结构），`LoadConfig()` 加载配置 |
| `config/theme.go` | 定义 `SegmentID`、`ThemeMode`、`SegmentTheme`，`GetSegmentTheme()` 获取主题，`ApplyColor()` 应用颜色 |
| `segment/segment.go` | 定义 `SegmentData`、`SegmentResult` 结构体和 `Segment` 接口，`formatTokenCount()` 辅助函数 |
| `segment/*.go` | 各 Segment 实现，每个文件包含一个 `XxxSegment` 结构体和 `Collect()` 方法 |
| `render/statusline.go` | `StatusLineGenerator` 将 `[]SegmentResult` 渲染为最终字符串 |
| `tui/tui.go` | Bubbletea TUI 应用，交互式配置界面，`Run()` 启动 TUI |

## 关键类型

### 输入数据（定义在 `config/config.go`）

- `InputData` - stdin JSON 结构，包含 Model、Workspace、TranscriptPath、Cost、OutputStyle
- `SimpleConfig` - TOML 配置，包含 Theme、Separator、Segments 开关

### Segment 系统（定义在 `segment/segment.go`）

- `SegmentData` - Segment 输出数据，包含 Primary（主显示内容）、Secondary、Metadata
- `SegmentResult` - 包含 SegmentID 和 SegmentData
- `Segment` 接口 - 定义 `Collect(input *InputData) SegmentData` 方法

### 主题系统（定义在 `config/theme.go`）

- `SegmentID` - 9 种 Segment 的枚举常量
- `ThemeMode` - `default`（Emoji）和 `nerd_font`（Nerd Font）两种模式
- `SegmentTheme` - 包含 Icon、IconColor、TextColor、BgColor、Bold
- 图标映射 - `defaultIcons` 和 `nerdFontIcons` 两个 map

## 14 种 Segment

| Segment | 文件 | 数据来源 | 输出示例 |
|---------|------|----------|----------|
| Model | `segment/model.go` | `input.Model.DisplayName` | `Sonnet 3.5` |
| Directory | `segment/directory.go` | `input.Workspace.CurrentDir` | `myapp` |
| Git | `segment/git.go` | 执行 git 命令 | `main *` |
| ContextWindow | `segment/context_window.go` | 解析 transcript 文件 | `15.6% · 31.1k tokens` |
| Usage | `segment/usage.go` | 解析 transcript 文件 | `↓31.1K ↑5.2K` |
| Cost | `segment/cost.go` | `input.Cost.TotalCostUSD` | `$0.15` |
| Session | `segment/session.go` | transcript 文件时间戳 | `1h23m` |
| OutputStyle | `segment/output_style.go` | `input.OutputStyle.Name` | `default` |
| Update | `segment/update.go` | 检查新版本 | 空或提示 |
| CCH Model | `segment/cch_model.go` | CCH API | `claude-3-opus` |
| CCH Provider | `segment/cch_provider.go` | CCH API | `anthropic` |
| CCH Cost | `segment/cch_cost.go` | CCH API | `$1.50/$10` |
| CCH Requests | `segment/cch_requests.go` | CCH API | `123 reqs` |
| CCH Limits | `segment/cch_limits.go` | CCH API | `5h:$0 W:$50 M:$100` |

## 配置文件

路径：`~/.claude/cchline/config.toml`

配置项：
- `theme` - 主题模式，`default` 或 `nerd_font`
- `separator` - 分隔符，默认 ` | `
- `cch_url` - CCH 服务器地址
- `cch_api_key` - CCH API 密钥
- `[segments]` - 各 Segment 的 bool 开关

## TUI 配置界面

文件：`tui/tui.go`

- 使用 Bubbletea 框架
- 通过 `cchline -c` 启动
- `Model` 结构体管理状态
- `Update()` 处理按键事件
- `View()` 渲染界面
- `saveConfig()` 保存到 TOML 文件
- 支持文本输入（CCH URL、API Key），使用 `bubbles/textinput` 组件

## 扩展指南

### 添加新 Segment

1. `config/theme.go` - 添加 `SegmentID` 常量，在 `defaultIcons` 和 `nerdFontIcons` 添加图标，在 `segmentColors` 添加颜色
2. `segment/` - 创建新文件，定义 `XxxSegment` 结构体，实现 `Collect()` 方法
3. `config/config.go` - 在 `SegmentToggles` 添加 bool 字段
4. `main.go` - 在 `collectAllSegments()` 的 collectors 切片中注册
5. `tui/tui.go` - 在 `NewModel()` 的 items 切片中添加菜单项

### 添加新主题

1. `config/theme.go` - 添加 `ThemeMode` 常量
2. `config/theme.go` - 创建新的图标映射 map
3. `config/theme.go` - 在 `GetSegmentTheme()` 的 switch 中添加 case

## 依赖

- `github.com/BurntSushi/toml` - TOML 解析
- `github.com/charmbracelet/lipgloss` / `github.com/muesli/termenv` - ANSI 颜色
- `github.com/charmbracelet/bubbletea` - TUI 框架
- `github.com/charmbracelet/lipgloss` - TUI 样式
- `github.com/charmbracelet/bubbles/textinput` - TUI 文本输入组件

## 构建

```bash
go build -ldflags="-s -w" -o cchline
```
