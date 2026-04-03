package population

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/resas-cli/cmd/area"
	"github.com/planitaicojp/resas-cli/cmd/cmdutil"
	"github.com/planitaicojp/resas-cli/internal/api"
	"github.com/planitaicojp/resas-cli/internal/output"
)

var parentGetAPIKey func() string
var parentGetFormat func() string

func SetParentAccessors(getAPIKey, getFormat func() string) {
	parentGetAPIKey = getAPIKey
	parentGetFormat = getFormat
}

var Cmd = &cobra.Command{
	Use:   "population",
	Short: "人口データを取得",
}

var prefCode int
var cityCode string

var compositionCmd = &cobra.Command{
	Use:   "composition",
	Short: "人口構成を取得",
	Long:  "指定した都道府県（市区町村）の人口構成（総人口、年少人口、生産年齢人口、老年人口）を取得します。",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(parentGetAPIKey())
		if err != nil {
			return err
		}

		if prefCode == 0 {
			code, err := area.SelectPrefecture(client)
			if err != nil {
				return err
			}
			prefCode = code
		}

		if cityCode == "" {
			cityCode = "-"
		}

		pop := api.NewPopulationAPI(client)
		result, err := pop.GetComposition(prefCode, cityCode)
		if err != nil {
			return err
		}

		format := cmdutil.GetFormat(parentGetFormat())

		if format == "json" {
			return output.New("json").Format(os.Stdout, result)
		}

		if len(result.Data) == 0 {
			fmt.Fprintln(os.Stderr, "データがありません。")
			return nil
		}

		type row struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		}
		rows := make([]row, len(result.Data[0].Data))
		for i, item := range result.Data[0].Data {
			rows[i] = row{Year: item.Year, Value: item.Value}
		}
		return output.New(format).Format(os.Stdout, rows)
	},
}

func init() {
	compositionCmd.Flags().IntVar(&prefCode, "pref-code", 0, "都道府県コード（例: 13）")
	compositionCmd.Flags().StringVar(&cityCode, "city-code", "", "市区町村コード（例: 13101、省略時は都道府県全体）")
	Cmd.AddCommand(compositionCmd)
}
