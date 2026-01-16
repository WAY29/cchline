# CCometixLine Golang 重构指南

## 一、项目概述

**CCometixLine** 是一个 Claude Code 状态栏增强工具，用 Rust 编写，约 9,772 行代码。本文档为 Golang 重构提供完整的技术参考。

### 核心功能
- 从 stdin 读取 Claude Code 传入的 JSON 数据
- 收集各类状态信息（Git、目录、模型、Token 使用等）
- 生成带 ANSI 颜色的状态栏字符串输出到 stdout

---

## 二、项目结构

### 原始目录结构（精简后）

```
CCometixLine/
├── src/
│   ├── main.rs                    # 程序入口
│   ├── cli.rs                     # CLI 参数解析
│   ├── config/
│   │   ├── types.rs               # 配置数据结构
│   │   ├── loader.rs              # 配置加载
│   │   └── models.rs              # 模型配置
│   ├── core/
│   │   ├── statusline.rs          # 状态栏生成器（核心）
│   │   └── segments/              # 9 种状态段
│   │       ├── directory.rs       # 当前目录
│   │       ├── git.rs             # Git 状态
│   │       ├── model.rs           # 模型名称
│   │       ├── context_window.rs  # 上下文使用率
│   │       ├── usage.rs           # Token 使用
│   │       ├── cost.rs            # API 成本
│   │       ├── session.rs         # 会话时长
│   │       ├── output_style.rs    # 输出风格
│   │       └── update.rs          # 更新提示
│   └── ui/                        # TUI 配置界面（可选，重构时可省略）
│       └── themes/                # 9 个主题
└── Cargo.toml
```

### Golang 重构建议结构

```
ccline/
├── main.go                        # 程序入口
├── config/
│   ├── config.go                  # 配置结构和加载
│   └── theme.go                   # 主题定义（仅 Cometix）
├── segment/
│   ├── segment.go                 # Segment 接口定义
│   ├── directory.go
│   ├── git.go
│   ├── model.go
│   ├── context_window.go
│   ├── usage.go
│   ├── cost.go
│   ├── session.go
│   ├── output_style.go
│   └── update.go
├── render/
│   └── statusline.go              # 状态栏渲染器
└── go.mod
```

---

## 三、核心数据结构

### 1. 输入数据（从 stdin 读取的 JSON）

```go
// InputData 是 Claude Code 传入的 JSON 结构
type InputData struct {
    Model          ModelInfo      `json:"model"`
    Workspace      WorkspaceInfo  `json:"workspace"`
    TranscriptPath string         `json:"transcript_path"`
    Cost           CostInfo       `json:"cost"`
    OutputStyle    OutputStyleInfo `json:"output_style"`
}

type ModelInfo struct {
    ID          string `json:"id"`
    DisplayName string `json:"display_name"`
}

type WorkspaceInfo struct {
    CurrentDir string `json:"current_dir"`
}

type CostInfo struct {
    TotalCostUSD float64 `json:"total_cost_usd"`
}

type OutputStyleInfo struct {
    Name string `json:"name"`
}
```

### 2. 配置结构（简化版）

```go
// Config 主配置结构
type Config struct {
    Separator string           `toml:"separator"`
    Segments  []SegmentConfig  `toml:"segments"`
}

// SegmentConfig 段配置（简化为仅启用/禁用）
type SegmentConfig struct {
    ID      SegmentID `toml:"id"`
    Enabled bool      `toml:"enabled"`
}

// SegmentID 段类型枚举
type SegmentID string

const (
    SegmentModel         SegmentID = "model"
    SegmentDirectory     SegmentID = "directory"
    SegmentGit           SegmentID = "git"
    SegmentContextWindow SegmentID = "context_window"
    SegmentUsage         SegmentID = "usage"
    SegmentCost          SegmentID = "cost"
    SegmentSession       SegmentID = "session"
    SegmentOutputStyle   SegmentID = "output_style"
    SegmentUpdate        SegmentID = "update"
)
```

### 3. 段数据结构

```go
// SegmentData 段渲染数据
type SegmentData struct {
    Primary   string            // 主要显示内容
    Secondary string            // 次要内容（可选）
    Metadata  map[string]string // 元数据
}
```

---

## 四、核心渲染逻辑

### 1. 主流程

```
stdin (JSON) → InputData → collect_all_segments() → StatusLineGenerator.generate() → stdout
```

### 2. 渲染器核心实现

```go
// StatusLineGenerator 状态栏生成器
type StatusLineGenerator struct {
    config *Config
}

func (g *StatusLineGenerator) Generate(segments []SegmentResult) string {
    var output []string

    for _, seg := range segments {
        if !seg.Config.Enabled {
            continue
        }
        rendered := g.renderSegment(seg)
        if rendered != "" {
            output = append(output, rendered)
        }
    }

    return strings.Join(output, g.config.Separator)
}

func (g *StatusLineGenerator) renderSegment(seg SegmentResult) string {
    // 获取图标和颜色（从 Cometix 主题）
    theme := GetCometixTheme(seg.Config.ID)

    // 应用 ANSI 颜色
    icon := applyColor(theme.Icon, theme.IconColor)
    text := applyColor(seg.Data.Primary, theme.TextColor)

    // 如果有背景色
    if theme.BgColor != nil {
        return applyBackground(fmt.Sprintf(" %s %s ", icon, text), theme.BgColor)
    }

    return fmt.Sprintf("%s %s", icon, text)
}
```

### 3. ANSI 颜色处理
使用以下第三方库处理 ANSI 颜色：
https://github.com/fatih/color
```go
// AnsiColor 颜色类型
type AnsiColor struct {
    Type  string // "16", "256", "rgb"
    Value interface{}
}

// applyColor 应用前景色
func applyColor(text string, color *AnsiColor) string {
    if color == nil {
        return text
    }

    switch color.Type {
    case "16":
        c := color.Value.(int)
        var code int
        if c < 8 {
            code = 30 + c
        } else {
            code = 90 + (c - 8)
        }
        return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, text)
    case "256":
        c := color.Value.(int)
        return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", c, text)
    case "rgb":
        rgb := color.Value.([]int)
        return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", rgb[0], rgb[1], rgb[2], text)
    }
    return text
}

// applyBackground 应用背景色
func applyBackground(text string, color *AnsiColor) string {
    if color == nil {
        return text
    }

    switch color.Type {
    case "16":
        c := color.Value.(int)
        var code int
        if c < 8 {
            code = 40 + c
        } else {
            code = 100 + (c - 8)
        }
        return fmt.Sprintf("\x1b[%dm%s\x1b[49m", code, text)
    case "256":
        c := color.Value.(int)
        return fmt.Sprintf("\x1b[48;5;%dm%s\x1b[49m", c, text)
    case "rgb":
        rgb := color.Value.([]int)
        return fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s\x1b[49m", rgb[0], rgb[1], rgb[2], text)
    }
    return text
}
```

---

## 五、9 种 Segment 实现

### 1. Directory Segment

```go
func (s *DirectorySegment) Collect(input *InputData) SegmentData {
    dir := input.Workspace.CurrentDir
    // 只取最后一级目录名
    name := filepath.Base(dir)
    return SegmentData{Primary: name}
}
```

### 2. Git Segment

```go
func (s *GitSegment) Collect(input *InputData) SegmentData {
    dir := input.Workspace.CurrentDir

    // 获取当前分支
    branch := execGit(dir, "rev-parse", "--abbrev-ref", "HEAD")
    if branch == "" {
        return SegmentData{} // 非 Git 仓库
    }

    // 获取状态
    status := execGit(dir, "status", "--porcelain")

    // 构建显示内容
    result := branch
    if status != "" {
        result += " *" // 有未提交更改
    }

    // 获取 ahead/behind
    ahead, behind := getAheadBehind(dir)
    if ahead > 0 {
        result += fmt.Sprintf(" ↑%d", ahead)
    }
    if behind > 0 {
        result += fmt.Sprintf(" ↓%d", behind)
    }

    return SegmentData{Primary: result}
}

func execGit(dir string, args ...string) string {
    cmd := exec.Command("git", args...)
    cmd.Dir = dir
    out, err := cmd.Output()
    if err != nil {
        return ""
    }
    return strings.TrimSpace(string(out))
}
```

### 3. Model Segment

```go
func (s *ModelSegment) Collect(input *InputData) SegmentData {
    // 简化模型名称显示
    name := input.Model.DisplayName
    if name == "" {
        name = input.Model.ID
    }

    // 简化常见模型名
    name = simplifyModelName(name)

    return SegmentData{Primary: name}
}

func simplifyModelName(name string) string {
    replacements := map[string]string{
        "claude-3-5-sonnet": "Sonnet 3.5",
        "claude-3-opus":     "Opus 3",
        "claude-3-haiku":    "Haiku 3",
        "gpt-4":             "GPT-4",
    }
    for k, v := range replacements {
        if strings.Contains(strings.ToLower(name), k) {
            return v
        }
    }
    return name
}
```

### 4. ContextWindow Segment

```go
func (s *ContextWindowSegment) Collect(input *InputData) SegmentData {
    if input.TranscriptPath == "" {
        return SegmentData{}
    }

    // 解析 transcript 文件计算 token 使用
    totalTokens := parseTranscriptTokens(input.TranscriptPath)

    // 获取模型的上下文限制
    contextLimit := getModelContextLimit(input.Model.ID)

    // 计算使用率
    percentage := float64(totalTokens) / float64(contextLimit) * 100

    return SegmentData{
        Primary: fmt.Sprintf("%.0f%%", percentage),
        Metadata: map[string]string{
            "tokens": fmt.Sprintf("%d", totalTokens),
            "limit":  fmt.Sprintf("%d", contextLimit),
        },
    }
}
```

### 5. Usage Segment

```go
func (s *UsageSegment) Collect(input *InputData) SegmentData {
    if input.TranscriptPath == "" {
        return SegmentData{}
    }

    // 解析 transcript 统计 token
    inputTokens, outputTokens := parseTranscriptUsage(input.TranscriptPath)

    return SegmentData{
        Primary: fmt.Sprintf("↓%s ↑%s",
            formatTokenCount(inputTokens),
            formatTokenCount(outputTokens)),
    }
}

func formatTokenCount(count int) string {
    if count >= 1000000 {
        return fmt.Sprintf("%.1fM", float64(count)/1000000)
    }
    if count >= 1000 {
        return fmt.Sprintf("%.1fK", float64(count)/1000)
    }
    return fmt.Sprintf("%d", count)
}
```

### 6. Cost Segment

```go
func (s *CostSegment) Collect(input *InputData) SegmentData {
    cost := input.Cost.TotalCostUSD
    if cost == 0 {
        return SegmentData{}
    }
    return SegmentData{
        Primary: fmt.Sprintf("$%.2f", cost),
    }
}
```

### 7. Session Segment

```go
func (s *SessionSegment) Collect(input *InputData) SegmentData {
    if input.TranscriptPath == "" {
        return SegmentData{}
    }

    // 从 transcript 获取会话开始时间
    startTime := getSessionStartTime(input.TranscriptPath)
    if startTime.IsZero() {
        return SegmentData{}
    }

    duration := time.Since(startTime)
    return SegmentData{
        Primary: formatDuration(duration),
    }
}

func formatDuration(d time.Duration) string {
    hours := int(d.Hours())
    minutes := int(d.Minutes()) % 60

    if hours > 0 {
        return fmt.Sprintf("%dh%dm", hours, minutes)
    }
    return fmt.Sprintf("%dm", minutes)
}
```

### 8. OutputStyle Segment

```go
func (s *OutputStyleSegment) Collect(input *InputData) SegmentData {
    style := input.OutputStyle.Name
    if style == "" {
        return SegmentData{}
    }
    return SegmentData{Primary: style}
}
```

### 9. Update Segment

```go
func (s *UpdateSegment) Collect(input *InputData) SegmentData {
    // 检查是否有新版本（可选实现）
    // 建议使用缓存避免频繁检查
    return SegmentData{}
}
```

---

## 六、Cometix 主题定义（唯一保留的主题）

```go
// CometixTheme Cometix 主题配置
type CometixTheme struct {
    Icon      string
    IconColor *AnsiColor
    TextColor *AnsiColor
    BgColor   *AnsiColor
    Bold      bool
}

// GetCometixTheme 获取 Cometix 主题配置
func GetCometixTheme(id SegmentID) CometixTheme {
    themes := map[SegmentID]CometixTheme{
        SegmentModel: {
            Icon:      "🤖",  // Nerd Font: \ue26d
            TextColor: &AnsiColor{Type: "16", Value: 14}, // Light Cyan
            Bold:      true,
        },
        SegmentDirectory: {
            Icon:      "📁",  // Nerd Font: \uf024b
            IconColor: &AnsiColor{Type: "16", Value: 11}, // Light Yellow
            TextColor: &AnsiColor{Type: "16", Value: 10}, // Light Green
            Bold:      true,
        },
        SegmentGit: {
            Icon:      "🌿",  // Nerd Font: \uf02a2
            TextColor: &AnsiColor{Type: "16", Value: 12}, // Light Blue
            Bold:      true,
        },
        SegmentContextWindow: {
            Icon:      "⚡️", // Nerd Font: \uf49b
            TextColor: &AnsiColor{Type: "16", Value: 13}, // Light Magenta
            Bold:      true,
        },
        SegmentUsage: {
            Icon:      "📊",  // Nerd Font: \uf0a9e
            TextColor: &AnsiColor{Type: "16", Value: 14}, // Light Cyan
        },
        SegmentCost: {
            Icon:      "💰",  // Nerd Font: \ueec1
            TextColor: &AnsiColor{Type: "16", Value: 3},  // Yellow
            Bold:      true,
        },
        SegmentSession: {
            Icon:      "⏱️",  // Nerd Font: \uf19bb
            TextColor: &AnsiColor{Type: "16", Value: 2},  // Green
            Bold:      true,
        },
        SegmentOutputStyle: {
            Icon:      "🎯",  // Nerd Font: \uf12f5
            TextColor: &AnsiColor{Type: "16", Value: 6},  // Cyan
            Bold:      true,
        },
        SegmentUpdate: {
            Icon:      "🔄",
            TextColor: &AnsiColor{Type: "16", Value: 11}, // Light Yellow
        },
    }

    if theme, ok := themes[id]; ok {
        return theme
    }
    return CometixTheme{}
}
```

### Cometix 配色方案

| Segment | ANSI 颜色 | 颜色名称 | 默认启用 |
|---------|----------|---------|---------|
| Model | 14 | Light Cyan | ✅ |
| Directory | 11/10 | Light Yellow/Green | ✅ |
| Git | 12 | Light Blue | ✅ |
| ContextWindow | 13 | Light Magenta | ✅ |
| Usage | 14 | Light Cyan | ❌ |
| Cost | 3 | Yellow | ❌ |
| Session | 2 | Green | ❌ |
| OutputStyle | 6 | Cyan | ❌ |
| Update | 11 | Light Yellow | ❌ |

---

## 七、简化配置文件格式

### 配置文件位置
`~/.claude/ccline/config.toml`

### 简化配置示例

```toml
# CCometixLine 配置文件
separator = " | "

# 启用/禁用各个组件
[segments]
model = true
directory = true
git = true
context_window = true
usage = false
cost = false
session = false
output_style = false
update = false
```

### 配置加载代码

```go
type SimpleConfig struct {
    Separator string         `toml:"separator"`
    Segments  SegmentToggles `toml:"segments"`
}

type SegmentToggles struct {
    Model         bool `toml:"model"`
    Directory     bool `toml:"directory"`
    Git           bool `toml:"git"`
    ContextWindow bool `toml:"context_window"`
    Usage         bool `toml:"usage"`
    Cost          bool `toml:"cost"`
    Session       bool `toml:"session"`
    OutputStyle   bool `toml:"output_style"`
    Update        bool `toml:"update"`
}

func LoadConfig() (*SimpleConfig, error) {
    configPath := filepath.Join(os.Getenv("HOME"), ".claude", "ccline", "config.toml")

    // 默认配置
    config := &SimpleConfig{
        Separator: " | ",
        Segments: SegmentToggles{
            Model:         true,
            Directory:     true,
            Git:           true,
            ContextWindow: true,
            Usage:         false,
            Cost:          false,
            Session:       false,
            OutputStyle:   false,
            Update:        false,
        },
    }

    data, err := os.ReadFile(configPath)
    if err != nil {
        return config, nil // 使用默认配置
    }

    if err := toml.Unmarshal(data, config); err != nil {
        return nil, err
    }

    return config, nil
}
```

---

## 八、主程序入口

```go
package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
)

func main() {
    // 加载配置
    config, err := LoadConfig()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
        os.Exit(1)
    }

    // 检查 stdin 是否有数据
    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) != 0 {
        // 无 stdin 数据，显示帮助或退出
        fmt.Println("CCometixLine - Claude Code Status Line")
        fmt.Println("Usage: claude-code | ccline")
        return
    }

    // 读取 stdin JSON
    reader := bufio.NewReader(os.Stdin)
    var input InputData
    if err := json.NewDecoder(reader).Decode(&input); err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
        os.Exit(1)
    }

    // 收集所有段数据
    segments := collectAllSegments(config, &input)

    // 生成状态栏
    generator := NewStatusLineGenerator(config)
    output := generator.Generate(segments)

    // 输出到 stdout
    fmt.Print(output)
}

func collectAllSegments(config *SimpleConfig, input *InputData) []SegmentResult {
    var results []SegmentResult

    // 按顺序收集各段
    collectors := []struct {
        id      SegmentID
        enabled bool
        collect func(*InputData) SegmentData
    }{
        {SegmentModel, config.Segments.Model, (&ModelSegment{}).Collect},
        {SegmentDirectory, config.Segments.Directory, (&DirectorySegment{}).Collect},
        {SegmentGit, config.Segments.Git, (&GitSegment{}).Collect},
        {SegmentContextWindow, config.Segments.ContextWindow, (&ContextWindowSegment{}).Collect},
        {SegmentUsage, config.Segments.Usage, (&UsageSegment{}).Collect},
        {SegmentCost, config.Segments.Cost, (&CostSegment{}).Collect},
        {SegmentSession, config.Segments.Session, (&SessionSegment{}).Collect},
        {SegmentOutputStyle, config.Segments.OutputStyle, (&OutputStyleSegment{}).Collect},
        {SegmentUpdate, config.Segments.Update, (&UpdateSegment{}).Collect},
    }

    for _, c := range collectors {
        if c.enabled {
            data := c.collect(input)
            if data.Primary != "" {
                results = append(results, SegmentResult{
                    ID:   c.id,
                    Data: data,
                })
            }
        }
    }

    return results
}
```

---

## 九、Transcript 文件解析

Transcript 文件是 JSONL 格式，每行一个 JSON 对象：

```go
type TranscriptEntry struct {
    Role    string       `json:"role"`
    Content string       `json:"content"`
    Usage   *UsageInfo   `json:"usage,omitempty"`
}

type UsageInfo struct {
    // Anthropic 格式
    InputTokens              int `json:"input_tokens"`
    OutputTokens             int `json:"output_tokens"`
    CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
    CacheReadInputTokens     int `json:"cache_read_input_tokens"`

    // OpenAI 格式（兼容）
    PromptTokens     int `json:"prompt_tokens"`
    CompletionTokens int `json:"completion_tokens"`
}

func parseTranscriptTokens(path string) int {
    file, err := os.Open(path)
    if err != nil {
        return 0
    }
    defer file.Close()

    var totalTokens int
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        var entry TranscriptEntry
        if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
            continue
        }

        if entry.Usage != nil {
            // 优先使用 Anthropic 格式
            if entry.Usage.InputTokens > 0 {
                totalTokens += entry.Usage.InputTokens + entry.Usage.OutputTokens
            } else {
                // 回退到 OpenAI 格式
                totalTokens += entry.Usage.PromptTokens + entry.Usage.CompletionTokens
            }
        }
    }

    return totalTokens
}
```

---

## 十、模型上下文限制

```go
var modelContextLimits = map[string]int{
    "claude-3-opus":         200000,
    "claude-3-sonnet":       200000,
    "claude-3-haiku":        200000,
    "claude-3-5-sonnet":     200000,
    "claude-3-5-haiku":      200000,
    "gpt-4":                 128000,
    "gpt-4-turbo":           128000,
    "gpt-4o":                128000,
}

func getModelContextLimit(modelID string) int {
    modelID = strings.ToLower(modelID)

    for pattern, limit := range modelContextLimits {
        if strings.Contains(modelID, pattern) {
            return limit
        }
    }

    return 200000 // 默认值
}
```

---

## 十一、重构要点总结

### 保留的功能
1. ✅ 9 种 Segment（可通过配置启用/禁用）
2. ✅ Cometix 主题（唯一主题，硬编码）
3. ✅ ANSI 颜色渲染（16/256/RGB）
4. ✅ Git 状态检测
5. ✅ Transcript 解析
6. ✅ 简化的 TOML 配置

### 移除的功能
1. ❌ TUI 配置界面
2. ❌ 其他 8 个主题
3. ❌ 颜色/图标自定义
4. ❌ Powerline 箭头分隔符
5. ❌ Claude Code 补丁功能
6. ❌ 自动更新检查
7. ❌ NPM 发布相关

### Golang 依赖建议
```go
require (
    github.com/BurntSushi/toml v1.3.0  // TOML 解析
)
```

### 构建命令
```bash
go build -o ccline .
```

---

## 十二、测试数据

### 示例输入 JSON

```json
{
  "model": {
    "id": "claude-3-5-sonnet-20241022",
    "display_name": "Claude 3.5 Sonnet"
  },
  "workspace": {
    "current_dir": "/Users/user/projects/myapp"
  },
  "transcript_path": "/tmp/claude-transcript-12345.jsonl",
  "cost": {
    "total_cost_usd": 0.15
  },
  "output_style": {
    "name": "concise"
  }
}
```

### 预期输出

```
🤖 Sonnet 3.5 | 📁 myapp | 🌿 main | ⚡️ 45%
```

（带 ANSI 颜色代码）
