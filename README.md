# resas-cli

RESAS（地域経済分析システム）APIを操作するコマンドラインツール。
内閣府・経済産業省が提供する地域経済に関する各種統計データを取得します。

## 特徴

- シングルバイナリ（クロスプラットフォーム: Linux / macOS / Windows）
- 複数の出力形式（table / json / csv）
- 都道府県・市区町村の対話的選択（promptui）
- AIエージェント対応（`--no-input`、決定的終了コード、JSON エラー、stderr/stdout分離）
- 最小限の依存関係（cobra + promptui のみ）

## インストール

### Homebrew（macOS / Linux）

```bash
brew install planitaicojp/tap/resas
```

### Scoop（Windows）

```powershell
scoop bucket add planitaicojp https://github.com/planitaicojp/bucket
scoop install resas
```

### リリースバイナリ

[Releases](https://github.com/planitaicojp/resas-cli/releases) から対応するアーカイブをダウンロードしてください。

| OS | アーキテクチャ | ファイル |
|----|---------------|---------|
| Linux | amd64 | `resas-cli_*_linux_amd64.tar.gz` |
| Linux | arm64 | `resas-cli_*_linux_arm64.tar.gz` |
| macOS | amd64 (Intel) | `resas-cli_*_darwin_amd64.tar.gz` |
| macOS | arm64 (Apple Silicon) | `resas-cli_*_darwin_arm64.tar.gz` |
| Windows | amd64 | `resas-cli_*_windows_amd64.zip` |
| Windows | arm64 | `resas-cli_*_windows_arm64.zip` |

### ソースからビルド

```bash
go install github.com/planitaicojp/resas-cli@latest
```

または:

```bash
git clone https://github.com/planitaicojp/resas-cli.git
cd resas-cli
make build
```

## 前提条件

RESAS APIキーが必要です（無料登録）:

1. [RESAS-API](https://opendata.resas-portal.go.jp/) にアクセス
2. 利用登録を行い、APIキーを取得
3. APIキーを設定:

```bash
resas config set api_key YOUR_API_KEY
```

または環境変数で:

```bash
export RESAS_API_KEY=YOUR_API_KEY
```

## クイックスタート

```bash
# 都道府県一覧
resas area pref

# 市区町村一覧（対話的に都道府県を選択）
resas area city

# 東京都の市区町村一覧
resas area city --pref-code 13

# 東京都の人口構成
resas population composition --pref-code 13

# JSON出力
resas area pref --format json

# CSV出力
resas population composition --pref-code 13 --format csv
```

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| `resas area pref` | 都道府県一覧を表示 |
| `resas area city` | 市区町村一覧を表示 |
| `resas population composition` | 人口構成を取得 |
| `resas config show` | 現在の設定を表示 |
| `resas config set <key> <value>` | 設定値を更新 |
| `resas config path` | 設定ファイルのパスを表示 |
| `resas version` | バージョン情報を表示 |
| `resas completion <shell>` | シェル補完スクリプトを生成 |

## 設定

設定ファイル: `~/.config/resas/config.json`

```json
{
  "api_key": "YOUR_API_KEY",
  "defaults": {
    "format": "table",
    "pref_code": 13
  }
}
```

### 環境変数

| 変数 | 説明 |
|------|------|
| `RESAS_API_KEY` | APIキー |
| `RESAS_FORMAT` | デフォルト出力形式 |
| `RESAS_CONFIG_DIR` | 設定ディレクトリ |
| `RESAS_CACHE_DIR` | キャッシュディレクトリ |
| `RESAS_NO_INPUT` | 対話プロンプトを無効化（`1`で有効） |

### グローバルフラグ

| フラグ | 説明 |
|-------|------|
| `--api-key` | APIキー |
| `--format` | 出力形式: table, json, csv |
| `--no-input` | 対話プロンプトを無効化 |
| `--quiet` | 補助的な出力を抑制 |
| `--verbose` | 詳細ログを出力 |
| `--no-color` | カラー出力を無効化 |

## 終了コード

| コード | 意味 |
|-------|------|
| 0 | 成功 |
| 1 | 一般エラー |
| 2 | APIキーが未設定または無効 |
| 3 | データが見つからない |
| 4 | 入力パラメータが不正 |
| 5 | APIエラー |
| 6 | ネットワークエラー |
| 10 | ユーザーキャンセル |

## エージェント連携

AIエージェントから利用する場合:

```bash
# 対話プロンプトを無効化し、JSON出力
resas population composition --pref-code 13 --format json --no-input

# エラーもJSON構造で出力される
resas population composition --format json --no-input 2>&1
# → {"error":{"code":"validation_error","message":"...","exit_code":4}}
```

## 開発

```bash
make build      # ビルド
make test       # テスト
make lint       # リント (golangci-lint)
make coverage   # カバレッジ
make clean      # クリーンアップ
```

## API対応状況

### Phase 1（実装済み）

- 地域コード（都道府県・市区町村）
- 人口（人口構成）

### Phase 2（予定）

- 人口（ピラミッド、推移、将来推計）
- 産業（構造、企業数、付加価値額、労働生産性）
- 地域経済循環（生産・分配・支出）
- 観光（外国人訪問者、宿泊者数、観光施設）
- 雇用（有効求人倍率、就業者数）
- まちづくり（事業所立地、人口メッシュ）

## ライセンス

MIT License
