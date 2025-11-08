package model

type OpenSearchMeta struct {
	Code       int      `json:"code"`
	HasError   bool     `json:"has_error"`
	Error      string   `json:"error"`
	HasWarning bool     `json:"has_warning"`
	Warnings   []string `json:"warnings"`
}
