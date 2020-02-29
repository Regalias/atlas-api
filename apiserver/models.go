package apiserver

type genericResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}
