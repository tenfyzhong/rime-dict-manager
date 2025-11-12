# Rime Dict Manager

一个用于管理 [Rime](https://rime.im/) 用户词典的命令行工具, 特别适用于五笔86极点码用户词典 (`wubi86_jidian_user.dict.yaml`).

## 功能特性

- **词条管理**: 轻松添加, 更新, 查询和删除词条.
- **权重调整**: 快速设置词条的权重.
- **美观列表**: 以清晰, 对齐的格式列出所有词典条目, 并按组显示.
- **自动编码**: 为新词条自动生成五笔编码 (需要主词典文件).
- **自动部署**: 在修改词典后可自动触发 Rime 的重新部署.
- **灵活配置**: 通过命令行标志轻松配置词典文件路径和部署命令.

## 安装

确保你已经安装了 Go (版本 1.21 或更高), 然后运行以下命令:

```bash
go install github.com/tenfyzhong/rime-dict-manager@latest
```

## 配置

该工具通过以下顺序确定要使用的用户词典文件:

1. **`--file` / `-f` 标志**: 命令行中指定的最高优先级路径.
2. **默认路径**: `~/Library/Rime/wubi86_jidian_user.dict.yaml` (macOS).

你可以通过全局标志来自定义文件路径和行为:

- `--file, -f`: 指定用户词典文件的路径.
- `--main-dict`: 指定用于生成五笔编码的主词典文件路径 (默认为 `~/Library/Rime/wubi86_jidian.dict.yaml`).
- `--deploy-cmd`: 指定 Rime 重新部署时要执行的命令.
- `--no-deploy`: 禁用在操作后自动重新部署 Rime.

## 使用方法

### `list` - 列出所有词条

以美观, 对齐的格式打印出词典中的所有内容.

```bash
rime-dict-manager list
```

**输出示例:**

```
词典文件: /Users/me/Library/Rime/wubi86_jidian_user.dict.yaml

词语 (Word)              编码 (Code)         权重 (Weight)
-------------------------------------------------------

************************ 工作 *************************
用例                     etwg                200
解耦                     qedi                100
...
```

### `add` - 添加或更新词条

添加一个新词条. 如果词条已存在, 则会更新它.

```bash
rime-dict-manager add <词语> [标志]
```

**标志:**

- `--code, -c`: 手动指定五笔编码. 如果未提供, 将尝试自动生成.
- `--weight, -w`: 指定词条权重 (默认为 `100`).
- `--group, -g`: 指定词条所属的分组 (默认为 `个人`).

**示例:**

```bash
# 手动指定编码和权重
rime-dict-manager add 幂等 --code pjtf --weight 150

# 自动生成编码 (需要主词典)
rime-dict-manager add 区块链

# 添加到指定分组
rime-dict-manager add 哈希 --group 工作
```

### `query` - 查询词条

在词典中查找一个词条并显示其详细信息.

```bash
rime-dict-manager query <词语>
```

**示例:**

```bash
$ rime-dict-manager query 用例
Found entries for '用例':
- Word:   用例
  Code:   etwg
  Weight: 200
  Group:  工作
---
```

### `delete` - 删除词条

从词典中删除一个指定的词条.

```bash
rime-dict-manager delete <词语>
```

**示例:**

```bash
rime-dict-manager delete 触达
```

### `set-weight` - 设置权重

修改一个现有词条的权重.

```bash
rime-dict-manager set-weight <词语> <新权重>
```

**示例:**

```bash
rime-dict-manager set-weight 用例 15000
```

## 从源码构建

```bash
git clone https://github.com/tenfyzhong/rime-dict-manager.git
cd rime-dict-manager
go build .
```

之后, 你可以直接使用 `./rime-dict-manager` 命令.
