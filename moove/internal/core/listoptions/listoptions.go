package listoptions

type Request struct {
	Limit    int64  `json:"limit"`
	Continue string `json:"continue"`
}

type Response struct {
	Continue string      `json:"continue"`
	Items    interface{} `json:"items"`
}
