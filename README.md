# Treekei - File Tree with Line Counts

English | [简体中文](./README-CN.md)

![Platform](https://img.shields.io/badge/platform-macOS-lightgrey)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Motivation

Existing command-line tools like `tree` can display the directory structure, and `tokei` can count lines of code by language. However, there isn’t a simple way to quickly understand how code is distributed across a project.

`treekei` combines a file tree with line counts, allowing you to quickly grasp the code distribution of a project when you first explore it.

---

## Usage

<img src="./README-CN.assets/image-20260217180441732.png" alt="usage example" style="zoom:50%;" />

---

## Installation

### macOS

```shell
brew update
brew tap zihao-liu-qs/treekei
brew install treekei
```

### Linux

You can download the prebuilt binary from the GitHub Releases page and place it in a directory that is included in your `PATH`

## Options

- `--all`
  Include hidden files.
- `--dir-only`
  Display directories only.
- `--lang string`
  Only count code for the specified programming language.
- `--level int`
  Set the maximum depth of the output tree.
- `--sort string`
  Sorting method of the tree structure.
  Defaults to sorting by lines of code. Can be set to alphabetical order.

You can view all available options with:

```shell
treekei --help
```

------

## Supported Platforms

### Operating Systems

Currently tested on macOS (Apple Silicon).
In theory, it should work on any Unix-like system.
Windows is not officially supported yet, sorry.

### Installation Methods

Homebrew installation is currently provided for macOS only.
Users on other platforms can build from source:

```shell
go build -o treekei ./cmd/treekei/main.go
```

### Color Support

If your terminal supports colors, `treekei` will display colored output.

------

## License

MIT

------

If you find this tool useful, feel free to give it a ⭐ — it helps more people discover it!