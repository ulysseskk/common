package ocr

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	Host string `json:"host"`
}

type OCRClient struct {
	restyClient *resty.Client
}

func NewOCRClient(conf *Config) *OCRClient {
	restyClient := resty.New()
	restyClient.SetHostURL(conf.Host).SetHeader("content-type", "application/json")
	return &OCRClient{restyClient}
}

func (client *OCRClient) GetByBase64(ctx context.Context, fileBase64 string) (string, error) {
	resp, err := client.restyClient.R().SetContext(ctx).SetBody(&Base64Request{
		Base64: fileBase64,
		Trim:   "\n",
	}).Post("/base64")
	if err != nil {
		return "", err
	}
	rawResponse := resp.String()
	jsonResult := &FileResponse{}
	if err := json.Unmarshal([]byte(rawResponse), jsonResult); err != nil {
		return "", err
	}
	return jsonResult.Result, nil
}
