package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	iconfig "github.com/planitaicojp/resas-cli/internal/config"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "CLI設定を管理",
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の設定を表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := iconfig.Load()
		if err != nil {
			return err
		}

		display := *cfg
		if display.APIKey != "" {
			if len(display.APIKey) > 8 {
				display.APIKey = display.APIKey[:4] + "****" + display.APIKey[len(display.APIKey)-4:]
			} else {
				display.APIKey = "****"
			}
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(display)
	},
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "設定値を更新",
	Long: `設定値を更新します。

利用可能なキー:
  api_key      APIキー
  format       デフォルト出力形式 (table, json, csv)
  pref_code    デフォルト都道府県コード`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		cfg, err := iconfig.Load()
		if err != nil {
			return err
		}

		switch key {
		case "api_key":
			cfg.APIKey = value
		case "format":
			if value != "table" && value != "json" && value != "csv" {
				return fmt.Errorf("エラー: 無効な出力形式: %s（table, json, csv のいずれか）", value)
			}
			cfg.Defaults.Format = value
		case "pref_code":
			var code int
			if _, err := fmt.Sscanf(value, "%d", &code); err != nil || code < 1 || code > 47 {
				return fmt.Errorf("エラー: 無効な都道府県コード: %s（1〜47）", value)
			}
			cfg.Defaults.PrefCode = code
		default:
			return fmt.Errorf("エラー: 不明なキー: %s", key)
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "%s を設定しました。\n", key)
		return nil
	},
}

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "設定ファイルのパスを表示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(iconfig.Path())
	},
}

func init() {
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(setCmd)
	Cmd.AddCommand(pathCmd)
}
