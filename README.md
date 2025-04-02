# Go Find

一个使用 Go 语言实现的类似于 Linux find 命令的文件搜索工具。

## 功能特点

- 支持在指定目录下搜索文件和目录
- 支持基本的正则表达式匹配
- 支持命令行参数配置
- 支持并发文件查找，提高搜索效率
- 支持终端高亮显示搜索结果

## 安装

### 方法一：使用 go install

```bash
# 直接从 GitHub 安装
go install github.com/Lyrichu/go_find/cmd/go_find@latest

# 安装完成后，确保 $GOPATH/bin 在你的系统 PATH 中
# 现在可以直接使用 go_find 命令
go_find -help
```

### 方法二：从源码编译

```bash
# 克隆项目
git clone https://github.com/Lyrichu/go_find
cd go_find

# 安装依赖
go mod tidy

# 编译
go build -o go_find cmd/go_find/main.go
```

## 使用方法

```bash
# 基本用法
./go_find [options] [path]

# 示例：在当前目录下搜索所有 .go 文件
./go_find -name "*.go" .

# 示例：在指定目录下搜索所有名称中含有 "test" 的目录
./go_find -name test -type d /path/to/search

# 示例：同时搜索 .txt 和 .jpg 文件（使用正则表达式或运算符）
./go_find -name ".*\.txt$|.*\.jpg$" .

```

## 选项说明

- `-name pattern`: 使用正则表达式匹配文件名
- `-type [f|d]`: 指定搜索类型（f: 文件，d: 目录）
- `-help`: 显示帮助信息

## 开发说明

项目使用 Go 1.24 开发，遵循标准的 Go 项目结构：

- `cmd/`: 包含主程序入口
- `internal/`: 包含核心查找逻辑
- `internal/finder/`: 实现文件查找功能
- `internal/matcher/`: 实现正则匹配功能

## 许可证

MIT License