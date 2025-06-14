package v1

type EchoToolRequest struct {
	Message string `json:"message"`
}
type AddToolRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}
type HttpToolRequest struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Body   string `json:"body"`
}
