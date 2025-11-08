package ocr

type Base64Request struct {
	Base64 string `json:"base64"`
	Trim   string `json:"trim"`
}

type FileResponse struct {
	Result  string `json:"result"`
	Version string `json:"version"`
}
