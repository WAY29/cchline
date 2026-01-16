# CCHLine

Claude Code 状态栏增强工具，用 Go 编写。

> 本项目参考了 [CCometixLine](https://github.com/Haleclipse/CCometixLine) 的设计和代码结构，使用 Go 语言进行了重写。

## 功能

- 显示模型名称、当前目录、Git 状态、上下文使用率等信息
- 支持 9 种状态段（Segment），可自由启用/禁用
- 支持两套主题：`default`（Emoji）和 `nerd_font`（Nerd Font 图标）
- 交互式配置界面

## 安装

```bash
# 克隆项目
git clone https://github.com/WAY29/cchline.git
cd cchline

# 构建
go build -ldflags="-s -w" -o cchline

# 安装到 PATH（可选）
sudo mv cchline /usr/local/bin/
```

## 配置 Claude Code

在 `~/.claude/settings.json` 中添加：

```json
{
  "statusLineCommand": "/path/to/cchline"
}
```

## 使用

### 状态栏模式

```bash
# 由 Claude Code 自动调用
claude-code | cchline
```

### 交互式配置

```bash
cchline -c
```

操作说明：
- `↑` `↓` / `j` `k` - 上下移动
- `Space` / `Enter` - 切换选项
- `s` - 保存配置
- `Esc` - 退出

## 配置文件

路径：`~/.claude/cchline/config.toml`

```toml
theme = "nerd_font"  # "default" 或 "nerd_font"
separator = " | "

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

# CCH Segments
cch_model = false
cch_provider = false
cch_cost = false
cch_requests = false
cch_limits = false
```

## CCH 配置

CCHLine 支持连接 CCH (Claude Code Hub) 服务，显示额外的状态信息。

### 配置方式

**方式一：TUI 配置界面（推荐）**

```bash
cchline -c
```

在界面底部的 **CCH SETTINGS** 区域：
- 选择 `CCH URL`，按 `Enter` 输入服务器地址
- 选择 `API Key`，按 `Enter` 输入 API 密钥（输入时显示为 `****`）
- 按 `Esc` 保存并退出

**方式二：命令行参数**

```bash
cchline -u "https://your-cch-server.com" -k "your-api-key"
```

**方式三：直接编辑配置文件**

在 `~/.claude/cchline/config.toml` 中添加：

```toml
cch_url = "https://your-cch-server.com"
cch_api_key = "your-api-key"
```

### CCH Segments

配置 CCH 后，可启用以下状态段：

| Segment | 说明 |
|---------|------|
| CCH Model | CCH 上最后使用的模型 |
| CCH Provider | CCH 提供商名称 |
| CCH Cost | 今日成本/每日配额 |
| CCH Requests | 今日请求数 |
| CCH Limits | 5小时/周/月限额 |

## 状态段说明

| Segment | 说明 | 默认 |
|---------|------|------|
| Model | 当前使用的模型名称 | ✅ |
| Directory | 当前工作目录 | ✅ |
| Git | Git 分支和状态 | ✅ |
| Context Window | 上下文窗口使用率 | ✅ |
| Usage | Token 使用量 | ❌ |
| Cost | API 费用 | ❌ |
| Session | 会话时长 | ❌ |
| Output Style | 输出风格 | ❌ |
| Update | 更新提示 | ❌ |

## 主题

### default（Emoji）

```
🤖 Sonnet 3.5 | 📁 myapp | 🌿 main | ⚡️ 15.6% · 31.1k tokens
```

### nerd_font（Nerd Font）

```
 Sonnet 3.5 | 󰉋 myapp | 󰊢 main |  15.6% · 31.1k tokens
```

> 使用 `nerd_font` 主题需要终端安装 [Nerd Font](https://www.nerdfonts.com/) 字体。

## 依赖

- [BurntSushi/toml](https://github.com/BurntSushi/toml) - TOML 解析
- [fatih/color](https://github.com/fatih/color) - ANSI 颜色
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - TUI 框架
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - TUI 样式

## License

MIT
