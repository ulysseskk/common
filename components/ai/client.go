package client

import (
	"context"
	"net/http"
)

var (
	defaultClient Client
)

func GetDefaultClient() Client {
	return defaultClient
}

type Config interface {
	GetPassword() string
	GetModel() string
	GetBaseURL() string
	GetProxyEndpoint() string
	GetEndpointName() string
	GetEngine() string
	GetTemperature() float32
	GetProviderRegion() string
	GetTopP() float32
	GetTopK() int32
	GetMaxTokens() int
	GetProviderId() string
	GetCompartmentId() string
	GetOrganizationId() string
	GetCustomHeaders() []http.Header
}

type AIConfiguration struct {
	Providers       []AIProvider `json:"providers"`
	DefaultProvider string       `json:"defaultprovider"`
}

type AIProvider struct {
	Name           string        `json:"name" yaml:"name"`
	Model          string        `json:"model" yaml:"model"`
	Password       string        `json:"password" yaml:"password,omitempty"`
	BaseURL        string        `json:"baseurl" yaml:"baseurl,omitempty"`
	ProxyEndpoint  string        `json:"proxyEndpoint" yaml:"proxyEndpoint,omitempty"`
	ProxyPort      string        `json:"proxyPort" yaml:"proxyPort,omitempty"`
	EndpointName   string        `json:"endpointname" yaml:"endpointname,omitempty"`
	Engine         string        `json:"engine" yaml:"engine,omitempty"`
	Temperature    float32       `json:"temperature" yaml:"temperature,omitempty"`
	ProviderRegion string        `json:"providerregion" yaml:"providerregion,omitempty"`
	ProviderId     string        `json:"providerid" yaml:"providerid,omitempty"`
	CompartmentId  string        `json:"compartmentid" yaml:"compartmentid,omitempty"`
	TopP           float32       `json:"topp" yaml:"topp,omitempty"`
	TopK           int32         `json:"topk" yaml:"topk,omitempty"`
	MaxTokens      int           `json:"maxtokens" yaml:"maxtokens,omitempty"`
	OrganizationId string        `json:"organizationid" yaml:"organizationid,omitempty"`
	CustomHeaders  []http.Header `json:"customHeaders"`
}

func (p *AIProvider) GetBaseURL() string {
	return p.BaseURL
}

func (p *AIProvider) GetProxyEndpoint() string {
	return p.ProxyEndpoint
}

func (p *AIProvider) GetEndpointName() string {
	return p.EndpointName
}

func (p *AIProvider) GetTopP() float32 {
	return p.TopP
}

func (p *AIProvider) GetTopK() int32 {
	return p.TopK
}

func (p *AIProvider) GetMaxTokens() int {
	return p.MaxTokens
}

func (p *AIProvider) GetPassword() string {
	return p.Password
}

func (p *AIProvider) GetModel() string {
	return p.Model
}

func (p *AIProvider) GetEngine() string {
	return p.Engine
}
func (p *AIProvider) GetTemperature() float32 {
	return p.Temperature
}

func (p *AIProvider) GetProviderRegion() string {
	return p.ProviderRegion
}

func (p *AIProvider) GetProviderId() string {
	return p.ProviderId
}

func (p *AIProvider) GetCompartmentId() string {
	return p.CompartmentId
}

func (p *AIProvider) GetOrganizationId() string {
	return p.OrganizationId
}

func (p *AIProvider) GetCustomHeaders() []http.Header {
	return p.CustomHeaders
}

type Client interface {
	Configure(config Config) error
	GetCompletion(ctx context.Context, prompt string) (string, error)
	GetName() string
	Close()
}

type nopCloser struct{}

func (nopCloser) Close() {}
