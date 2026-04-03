package area

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/resas-cli/cmd/cmdutil"
	"github.com/planitaicojp/resas-cli/internal/api"
	"github.com/planitaicojp/resas-cli/internal/output"
	"github.com/planitaicojp/resas-cli/internal/prompt"
)

var parentGetAPIKey func() string
var parentGetFormat func() string

func SetParentAccessors(getAPIKey, getFormat func() string) {
	parentGetAPIKey = getAPIKey
	parentGetFormat = getFormat
}

var Cmd = &cobra.Command{
	Use:   "area",
	Short: "地域コード（都道府県・市区町村）を表示",
}

var prefCmd = &cobra.Command{
	Use:   "pref",
	Short: "都道府県一覧を表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(parentGetAPIKey())
		if err != nil {
			return err
		}

		areaAPI := api.NewAreaAPI(client)
		prefs, err := areaAPI.GetPrefectures()
		if err != nil {
			return err
		}

		type row struct {
			PrefCode int    `json:"prefCode"`
			PrefName string `json:"prefName"`
		}
		rows := make([]row, len(prefs))
		for i, p := range prefs {
			rows[i] = row{PrefCode: p.PrefCode, PrefName: p.PrefName}
		}

		format := cmdutil.GetFormat(parentGetFormat())
		return output.New(format).Format(os.Stdout, rows)
	},
}

var cityPrefCode int

var cityCmd = &cobra.Command{
	Use:   "city",
	Short: "市区町村一覧を表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(parentGetAPIKey())
		if err != nil {
			return err
		}

		if cityPrefCode == 0 {
			pref, err := SelectPrefecture(client)
			if err != nil {
				return err
			}
			cityPrefCode = pref
		}

		areaAPI := api.NewAreaAPI(client)
		cities, err := areaAPI.GetCities(cityPrefCode)
		if err != nil {
			return err
		}

		type row struct {
			CityCode    string `json:"cityCode"`
			CityName    string `json:"cityName"`
			BigCityFlag string `json:"bigCityFlag"`
		}
		rows := make([]row, len(cities))
		for i, c := range cities {
			rows[i] = row{CityCode: c.CityCode, CityName: c.CityName, BigCityFlag: c.BigCityFlag}
		}

		format := cmdutil.GetFormat(parentGetFormat())
		return output.New(format).Format(os.Stdout, rows)
	},
}

// SelectPrefecture is exported so other commands can reuse it
func SelectPrefecture(client *api.Client) (int, error) {
	areaAPI := api.NewAreaAPI(client)
	prefs, err := areaAPI.GetPrefectures()
	if err != nil {
		return 0, err
	}

	items := make([]prompt.SelectItem, len(prefs))
	for i, p := range prefs {
		items[i] = prompt.SelectItem{
			Label: fmt.Sprintf("%02d: %s", p.PrefCode, p.PrefName),
			Value: fmt.Sprintf("%d", p.PrefCode),
		}
	}

	val, err := prompt.Select("都道府県を選択してください", items)
	if err != nil {
		return 0, err
	}

	var code int
	fmt.Sscanf(val, "%d", &code)
	return code, nil
}

func init() {
	cityCmd.Flags().IntVar(&cityPrefCode, "pref-code", 0, "都道府県コード")
	Cmd.AddCommand(prefCmd)
	Cmd.AddCommand(cityCmd)
}
