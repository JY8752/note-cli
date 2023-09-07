<p align="center">
  <img src="./assets/logo.jpeg" width="500"/>
</p>

<p align="center">
Note CLI tool!!
</p>

<p align='center'>
<a href="https://github.com/JY8752/note-cli/releases/latest"><img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/JY8752/note-cli?style=flat"></a>
<a href="https://github.com/JY8752/note-cli/releases/latest"><img alt="GitHub all releases" src="https://img.shields.io/github/downloads/JY8752/note-cli/total?style=flat"></a>
<a href="./LICENSE"><img src="https://img.shields.io/github/license/JY8752/note-cli?style=flat" /></a>
<!-- <a href="https://github.com/JY8752/note-cli/actions/workflows/ci.yml"><img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/JY8752/note-cli/ci.yml?branch=main&logo=github&style=flat" /></a>
<a href="https://codeclimate.com/github/JY8752/note-cli/maintainability"><img alt="Code Climate maintainability" src="https://img.shields.io/codeclimate/maintainability/JY8752/note-cli?logo=codeclimate&style=flat" /></a> -->
<a href="https://goreportcard.com/report/github.com/JY8752/note-cli"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/JY8752/note-cli" /></a>
<a href="https://codecov.io/github/JY8752/note-cli"><img alt="Codecov" src="https://img.shields.io/codecov/c/github/JY8752/note-cli?logo=codecov" /></a>
</p>

<p align="center">
<a href="./README.md">English</a> | 日本語
</p>

# note-cli

**note-cli**は記事投稿サイトである[note](https://note.com/)の記事作成、管理をサービス上ではなく自分のローカル環境で行うためのCLIツールです。note-cliを作成したモチベーションは以下の通りです。

- webブラウザで開いたエディタではなく自分の好きなエディタで執筆したい。(例えば、VSCode)
- markdown形式で執筆した記事をGitHubなどでバージョン管理したい。
- 記事画像を用意するのがめんどくさい。(noteが用意してくれている無料の画像はあまり使いたくない)

## note-cliでできること

- 記事をすぐに執筆できるようにmarkdownファイルを含む記事ディレクトリをコマンドで作成します。
- 記事にアップロードするための画像をコマンドで作成することができます。

## インストール

```
go install github.com/JY8752/note-cli@latest
```

```
% note-cli -h

note-cli is a CLI command tool for creating, writing, and managing note articles

Usage:
  note-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create a new article directory.
  help        Help about any command

Flags:
  -h, --help     help for note-cli
  -t, --toggle   Help message for toggle

Use "note-cli [command] --help" for more information about a command.
```

## はじめに

1. 記事を管理するディレクトリを作成し移動します。

```
mkdir note-cli-demo
cd note-cli-demo
```

2. 記事ディレクトリを作成します。

```
% note-cli create article

Create directory. a6b420c6-9bb2-4060-869c-20c171fc9827
Create file. a6b420c6-9bb2-4060-869c-20c171fc9827.md
```

```
.
└── a6b420c6-9bb2-4060-869c-20c171fc9827
    └── a6b420c6-9bb2-4060-869c-20c171fc9827.md
```

3. 記事画像を生成する。

```
```