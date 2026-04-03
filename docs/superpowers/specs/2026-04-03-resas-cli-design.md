# resas-cli 設計仕様書

## 概要

RESAS（地域経済分析システム）APIを操作するためのGo製CLIツール。
内閣府・経済産業省が提供する地域経済に関する各種統計データを取得する。

参照実装: `~/dev/crowdy/conoha-cli` (Cobra + サービス指向APIレイヤー)

## 設計方針

- conoha-cliの検証済みパターンを採用し、RESAS APIの特性に合わせて簡素化
- AI エージェント（Claude Code等）からの利用を前提とした設計
- 全メッセージ・ヘルプを日本語で統一
- 最小限の依存関係

---

## 1. ディレクトリ構造

```
resas-cli/
├── main.go
├── cmd/
│   ├── root.go                    # ルートコマンド + グローバルフラグ
│   ├── version.go                 # バージョン表示
│   ├── completion.go              # シェル補完
│   ├── config/                    # config show, set, path
│   ├── cmdutil/                   # 共通ヘルパー (NewClient, GetFormat, args検証)
│   ├── area/                      # 地域コード: pref, city
│   ├── population/                # 人口: composition, pyramid, change, estimate
│   ├── industry/                  # 産業: structure, company, added-value, productivity
│   ├── economy/                   # 地域経済循環: production, distribution, expenditure
│   ├── tourism/                   # 観光: foreigners, hotel, facility
│   ├── employment/                # 雇用: job-ratio, workers
│   └── municipality/              # まちづくり: business-location, mesh
│
├── internal/
│   ├── api/
│   │   ├── client.go              # HTTPクライアント (APIキーヘッダー、リトライ、デバッグ)
│   │   ├── area.go                # 地域コードAPI (prefectures, cities)
│   │   ├── population.go          # 人口API
│   │   ├── industry.go            # 産業API
│   │   ├── economy.go             # 経済循環API
│   │   ├── tourism.go             # 観光API
│   │   ├── employment.go          # 雇用API
│   │   └── municipality.go        # まちづくりAPI
│   │
│   ├── config/
│   │   ├── config.go              # config.json 読み書き
│   │   └── cache.go               # 地域コードキャッシュ (TTL: 30日)
│   │
│   ├── model/                     # APIレスポンス構造体
│   ├── output/
│   │   ├── formatter.go           # Formatterインターフェース
│   │   ├── table.go               # text/tabwriter
│   │   ├── json.go                # encoding/json
│   │   └── csv.go                 # encoding/csv
│   │
│   ├── errors/                    # エラー型 + 終了コード
│   └── prompt/                    # TTY selection (promptui)
│
└── test/                          # テストフィクスチャ
```

---

## 2. 認証・設定

### 設定ファイル

パス: `~/.config/resas/config.json`

```json
{
  "api_key": "your-resas-api-key",
  "defaults": {
    "format": "table",
    "pref_code": 13
  }
}
```

- プロファイル機能なし（単一設定）
- credentials/tokens分離なし（APIキー1つのみ）

### APIキー解決順序（優先度高→低）

1. `--api-key` フラグ
2. `RESAS_API_KEY` 環境変数
3. config.jsonの `api_key`

### キャッシュ

パス: `~/.cache/resas/`

| ファイル | 内容 | TTL |
|---------|------|-----|
| `prefectures.json` | 都道府県一覧 | 30日 |
| `cities_{prefCode}.json` | 市区町村一覧（都道府県別） | 30日 |

---

## 3. HTTPクライアント

### 基本仕様

- ヘッダー: `X-API-KEY: {api_key}`
- Base URL: `https://opendata.resas-portal.go.jp`
- タイムアウト: 30秒
- User-Agent: `planitai/resas-cli/{version}`

### リトライ

- 対象: HTTP 429（レート制限）、5xx（サーバーエラー）
- 最大: 3回
- バックオフ: `(attempt + 1) * 1秒`
- POST/PUTはリトライしない（bodyが消費済みのため）

### デバッグ

- `--verbose` フラグでリクエスト/レスポンスのログ出力
- APIキーはマスク表示

---

## 4. コマンド体系

### サブコマンド構造

| コマンド | サブコマンド | RESAS APIエンドポイント |
|---------|------------|----------------------|
| `resas area` | `pref`, `city` | `/api/v1/prefectures`, `/api/v1/cities` |
| `resas population` | `composition`, `pyramid`, `change`, `estimate` | `/api/v1/population/` 配下 |
| `resas industry` | `structure`, `company`, `added-value`, `productivity` | `/api/v1/industry/` 配下 |
| `resas economy` | `production`, `distribution`, `expenditure` | `/api/v1/regionalEconomy/` 配下 |
| `resas tourism` | `foreigners`, `hotel`, `facility` | `/api/v1/tourism/` 配下 |
| `resas employment` | `job-ratio`, `workers` | `/api/v1/employment/` 配下 |
| `resas municipality` | `business-location`, `mesh` | `/api/v1/municipality/` 配下 |
| `resas config` | `show`, `set`, `path` | （ローカル） |

### グローバルフラグ

| フラグ | 型 | 説明 |
|-------|-----|------|
| `--api-key` | STRING | APIキー |
| `--format` | STRING | 出力形式: table, json, csv（デフォルト: table） |
| `--no-input` | BOOL | 対話プロンプトを無効化 |
| `--quiet` | BOOL | 補助的な出力を抑制 |
| `--verbose` | BOOL | 詳細ログを出力 |
| `--no-color` | BOOL | カラー出力を無効化 |

### データコマンド共通フラグ

| フラグ | 型 | 説明 |
|-------|-----|------|
| `--pref-code` | INT | 都道府県コード（例: 13） |
| `--city-code` | INT | 市区町村コード（例: 13101、一部コマンドのみ） |

### TTY Selection

`--pref-code` 未指定かつTTY接続時、promptuiで対話的に選択:

```
$ resas population composition
? 都道府県を選択してください:
  ▸ 01: 北海道
    02: 青森県
    03: 岩手県
    ...
    47: 沖縄県
```

- 選択後、`--city-code` が必要なコマンドでは市区町村selectionが続行
- `--no-input` モードで必須コード未指定時 → エラー + exit code 4

---

## 5. 出力形式

### フォーマット

**table**（デフォルト）:
```
年度    総人口      年少人口    生産年齢人口  老年人口
2020    13,960,000  1,520,000   8,940,000    3,500,000
```

**json**:
```json
[
  {
    "year": 2020,
    "total": 13960000,
    "young": 1520000,
    "working": 8940000,
    "elderly": 3500000
  }
]
```

**csv**:
```
年度,総人口,年少人口,生産年齢人口,老年人口
2020,13960000,1520000,8940000,3500000
```

### フォーマット決定順序（優先度高→低）

1. `--format` フラグ
2. `RESAS_FORMAT` 環境変数
3. config.json `defaults.format`
4. `table`（ハードコードデフォルト）

### stderr/stdout分離

- stdout → データ出力のみ（パイプ可能）
- stderr → プロンプト、進行メッセージ、エラーメッセージ

---

## 6. エラー処理

### 終了コード

| コード | 定数 | 意味 |
|-------|------|------|
| 0 | `ExitOK` | 成功 |
| 1 | `ExitGeneral` | 一般エラー |
| 2 | `ExitAuth` | APIキーが未設定または無効 |
| 3 | `ExitNotFound` | 指定されたデータが見つからない |
| 4 | `ExitValidation` | 入力パラメータが不正 |
| 5 | `ExitAPI` | APIエラー（レート制限等） |
| 6 | `ExitNetwork` | ネットワークエラー |
| 10 | `ExitCancelled` | ユーザーがキャンセル |

### JSON エラー出力

`--format json` 時、エラーもJSON構造で出力:

```json
{
  "error": {
    "code": "auth_error",
    "message": "APIキーが設定されていません。resas config set api_key <KEY> で設定してください。",
    "exit_code": 2
  }
}
```

### 通常エラー出力（stderr）

```
エラー: APIキーが設定されていません。resas config set api_key <KEY> で設定してください。
```

---

## 7. 依存関係

### 直接依存（2パッケージ）

| パッケージ | 用途 |
|-----------|------|
| `github.com/spf13/cobra` | CLIフレームワーク |
| `github.com/manifoldco/promptui` | TTY selection |

### 標準ライブラリ活用

| パッケージ | 用途 |
|-----------|------|
| `net/http` | HTTPクライアント |
| `encoding/json` | JSON入出力、config.json |
| `encoding/csv` | CSV出力 |
| `text/tabwriter` | table出力 |
| `os`, `path/filepath` | 環境変数、ファイルI/O |
| `time` | キャッシュTTL |

### ビルド

- Go 1.24+
- Makefile: build, test, lint, install
- .golangci.yml: linter設定
- クロスコンパイル: linux/darwin/windows (amd64, arm64)
- バージョン: ldflags注入 (`-X main.version=...`)

---

## 8. エージェント親和性

| 機能 | 詳細 |
|------|------|
| `--no-input` | 対話プロンプト無効、未指定の必須パラメータはエラー |
| 決定的終了コード | エラー種別ごとに固定の終了コード |
| JSON エラー | `--format json` 時、エラーも構造化JSON |
| stderr/stdout分離 | データはstdout、それ以外はstderr |
| `--quiet` | 進行メッセージを抑制 |

---

## conoha-cliからの変更点まとめ

| 項目 | conoha-cli | resas-cli |
|------|-----------|-----------|
| 認証 | OpenStack Keystone v3 + トークン | APIキー（X-API-KEY ヘッダー） |
| 設定ファイル | config.yaml + credentials.yaml + tokens.yaml | config.json のみ |
| プロファイル | マルチプロファイル | なし（単一設定） |
| 出力形式 | table, json, yaml, csv | table, json, csv |
| 設定形式 | YAML | JSON |
| TTY selection | サーバー/フレーバー選択 | 都道府県/市区町村選択 |
| キャッシュ | なし | 地域コードキャッシュ（TTL 30日） |
| 言語 | 英語 | 日本語 |
| YAML依存 | gopkg.in/yaml.v3 | なし |
