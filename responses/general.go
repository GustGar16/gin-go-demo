package responses

type GeneralResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"mesaage"`
	Data    map[string]interface{} `json:"data"`
}
