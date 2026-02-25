# treekei

简体中文 | [English](./README.md)

带行数的文件树

![Platform](https://img.shields.io/badge/platform-macOS-lightgrey)[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)



------

## 动机

现有的命令行工具如tree可以展示文件树，tokei可以分类统计代码行数，但是无法快速获得一个项目中代码量的分布。

treekei将文件树与代码行数融合在一起，让你在拿到一个新项目时，可以快速了解代码量的分布

## 使用

<img src="./README-CN.assets/image-20260217180441732.png" alt="image-20260217180441732" style="zoom:50%;" />

## 安装

### macOS

```shell
brew update
brew tap zihao-liu-qs/treekei
brew install treekei
```

### Linux

您可以从 GitHub Releases 页面下载预构建的二进制文件，并将其放置在“PATH”目录中

## 可选参数

- --all 包含隐藏文件
- --dir-only 只显示文件夹
- --lang string 只查找指定语言代码
- --level int 输出文件树的最大深度
- --sort string 树结构的排序方式，默认按代码行数，可以设置为按字母顺序

参数可通过如下命令查看

```shell
treekei --help
```

## 支持平台

### 系统

目前仅在macOS(m系列芯片)系统上进行验证，但理论上应该支持任何类unix系统；尚未对Window系统进行适配。

### 安装方式

目前仅为mac提供homebrew的安装方式，其他平台使用者可自行编译并使用 `go build -o treekei ./cmd/treekei/main.go`

### 颜色支持

如果您的终端支持颜色，treekei会有颜色支持

## LICENSE

MIT

---

如果你觉得这个工具有用，欢迎点个 ⭐，这会让更多人看到它~