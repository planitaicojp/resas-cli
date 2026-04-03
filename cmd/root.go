package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/resas-cli/cmd/area"
	configcmd "github.com/planitaicojp/resas-cli/cmd/config"
	"github.com/planitaicojp/resas-cli/cmd/population"
	"github.com/planitaicojp/resas-cli/internal/api"
	cerrors "github.com/planitaicojp/resas-cli/internal/errors"
	"github.com/planitaicojp/resas-cli/internal/prompt"
)

var version = "dev"

var (
	flagAPIKey  string
	flagFormat  string
	flagNoInput bool
	flagQuiet   bool
	flagVerbose bool
	flagNoColor bool
)

var rootCmd = &cobra.Command{
	Use:   "resas",
	Short: "RESAS API CLIツール",
	Long:  "RESAS（地域経済分析システム）APIを操作するコマンドラインツール",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if flagVerbose {
			api.Verbose = true
		}
		if flagNoInput {
			prompt.NoInputFlag = true
			os.Setenv("RESAS_NO_INPUT", "1")
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "APIキー")
	rootCmd.PersistentFlags().StringVar(&flagFormat, "format", "", "出力形式: table, json, csv（デフォルト: table）")
	rootCmd.PersistentFlags().BoolVar(&flagNoInput, "no-input", false, "対話プロンプトを無効化")
	rootCmd.PersistentFlags().BoolVar(&flagQuiet, "quiet", false, "補助的な出力を抑制")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "詳細ログを出力")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "カラー出力を無効化")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)

	area.SetParentAccessors(GetAPIKeyFlag, GetFormatFlag)
	rootCmd.AddCommand(area.Cmd)

	rootCmd.AddCommand(configcmd.Cmd)

	population.SetParentAccessors(GetAPIKeyFlag, GetFormatFlag)
	rootCmd.AddCommand(population.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if flagFormat == "json" {
			errJSON := cerrors.ToJSON(err)
			enc := json.NewEncoder(os.Stderr)
			enc.SetIndent("", "  ")
			enc.Encode(errJSON)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(cerrors.GetExitCode(err))
	}
}

func GetAPIKeyFlag() string { return flagAPIKey }
func GetFormatFlag() string { return flagFormat }
func IsQuiet() bool         { return flagQuiet }
