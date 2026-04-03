package api

import (
	"fmt"

	"github.com/planitaicojp/resas-cli/internal/model"
)

type AreaAPI struct {
	Client *Client
}

func NewAreaAPI(c *Client) *AreaAPI {
	return &AreaAPI{Client: c}
}

func (a *AreaAPI) GetPrefectures() ([]model.Prefecture, error) {
	var resp model.PrefecturesResponse
	if err := a.Client.Get("/api/v1/prefectures", &resp); err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (a *AreaAPI) GetCities(prefCode int) ([]model.City, error) {
	path := fmt.Sprintf("/api/v1/cities?prefCode=%d", prefCode)
	var resp model.CitiesResponse
	if err := a.Client.Get(path, &resp); err != nil {
		return nil, err
	}
	return resp.Result, nil
}
