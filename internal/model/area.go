package model

type Prefecture struct {
	PrefCode int    `json:"prefCode"`
	PrefName string `json:"prefName"`
}

type PrefecturesResponse struct {
	Message any          `json:"message"`
	Result  []Prefecture `json:"result"`
}

type City struct {
	PrefCode    int    `json:"prefCode"`
	CityCode    string `json:"cityCode"`
	CityName    string `json:"cityName"`
	BigCityFlag string `json:"bigCityFlag"`
}

type CitiesResponse struct {
	Message any    `json:"message"`
	Result  []City `json:"result"`
}
