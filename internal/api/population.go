package api

import (
	"fmt"

	"github.com/planitaicojp/resas-cli/internal/model"
)

type PopulationAPI struct {
	Client *Client
}

func NewPopulationAPI(c *Client) *PopulationAPI {
	return &PopulationAPI{Client: c}
}

func (a *PopulationAPI) GetComposition(prefCode int, cityCode string) (*model.PopulationCompositionResult, error) {
	path := fmt.Sprintf("/api/v1/population/composition/perYear?prefCode=%d&cityCode=%s", prefCode, cityCode)
	var resp model.PopulationCompositionResponse
	if err := a.Client.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp.Result, nil
}
