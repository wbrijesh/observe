package schema

type Response struct {
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
}
