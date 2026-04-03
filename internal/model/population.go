package model

type PopulationCompositionResponse struct {
	Message any                         `json:"message"`
	Result  PopulationCompositionResult `json:"result"`
}

type PopulationCompositionResult struct {
	BoundaryYear int                        `json:"boundaryYear"`
	Data         []PopulationCompositionData `json:"data"`
}

type PopulationCompositionData struct {
	Label string                      `json:"label"`
	Data  []PopulationCompositionItem `json:"data"`
}

type PopulationCompositionItem struct {
	Year  int `json:"year"`
	Value int `json:"value"`
}
