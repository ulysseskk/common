package clients

import (
	"context"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/ulysseskk/common/model/errors"
	"github.com/ulysseskk/common/model/rest"
	"net/http"
)

type ClientHelper struct {
	CodeSuccess int
}

func (helper *ClientHelper) MethodGet(ctx context.Context, restyClient *resty.Client, url string, pathParams, query map[string]string, responseData interface{}) (rest.Meta, error) {
	if restyClient == nil {
		return rest.Meta{}, errors.NewError().WithCode(rest.ClientError).WithMessage("Client has not been initialized")
	}
	resp, err := constructRequest(ctx, restyClient, pathParams, query, nil).
		Get(url)
	raw, err := validateResponse(resp, url, helper.CodeSuccess, err)
	if err != nil {
		return raw.Meta, err
	}
	err = dealWithResponse(raw, responseData)
	if err != nil {
		return raw.Meta, err
	}
	return raw.Meta, nil
}

func (helper *ClientHelper) MethodPut(ctx context.Context, restyClient *resty.Client, url string, pathParams, query map[string]string, body interface{}, responseData interface{}) (rest.Meta, error) {
	if restyClient == nil {
		return rest.Meta{}, errors.NewError().WithCode(rest.ClientError).WithMessage("Client has not been initialized")

	}
	resp, err := constructRequest(ctx, restyClient, pathParams, query, body).
		Put(url)
	raw, err := validateResponse(resp, url, helper.CodeSuccess, err)
	if err != nil {
		return raw.Meta, err
	}
	err = dealWithResponse(raw, responseData)
	if err != nil {
		return raw.Meta, err
	}
	return raw.Meta, nil
}

func (helper *ClientHelper) MethodPost(ctx context.Context, restyClient *resty.Client, url string, pathParams, query map[string]string, body interface{}, responseData interface{}) (rest.Meta, error) {
	if restyClient == nil {
		return rest.Meta{}, errors.NewError().WithCode(rest.ClientError).WithMessage("Client has not been initialized")

	}
	resp, err := constructRequest(ctx, restyClient, pathParams, query, body).
		Post(url)
	raw, err := validateResponse(resp, url, helper.CodeSuccess, err)
	if err != nil {
		return raw.Meta, err
	}
	err = dealWithResponse(raw, responseData)
	if err != nil {
		return raw.Meta, err
	}
	return raw.Meta, nil
}

func (helper *ClientHelper) MethodPatch(ctx context.Context, restyClient *resty.Client, url string, pathParams, query map[string]string, body interface{}, responseData interface{}) (rest.Meta, error) {
	if restyClient == nil {
		return rest.Meta{}, errors.NewError().WithCode(rest.ClientError).WithMessage("Client has not been initialized")

	}
	resp, err := constructRequest(ctx, restyClient, pathParams, query, body).
		Patch(url)
	raw, err := validateResponse(resp, url, helper.CodeSuccess, err)
	if err != nil {
		return raw.Meta, err
	}
	err = dealWithResponse(raw, responseData)
	if err != nil {
		return raw.Meta, err
	}
	return raw.Meta, nil
}

func (helper *ClientHelper) MethodDelete(ctx context.Context, restyClient *resty.Client, url string, pathParams, query map[string]string, body interface{}, responseData interface{}) (rest.Meta, error) {
	if restyClient == nil {
		return rest.Meta{}, errors.NewError().WithCode(rest.ClientError).WithMessage("Client has not been initialized")

	}
	resp, err := constructRequest(ctx, restyClient, pathParams, query, body).
		Delete(url)
	raw, err := validateResponse(resp, url, helper.CodeSuccess, err)
	if err != nil {
		return raw.Meta, err
	}
	err = dealWithResponse(raw, responseData)
	if err != nil {
		return raw.Meta, err
	}
	return raw.Meta, nil
}

type restyClientMethod func(client *resty.Request) *resty.Request

func (helper *ClientHelper) MethodPostCustom(ctx context.Context, restyClient *resty.Client, url string, pathParams, query map[string]string, body interface{}, responseData interface{}, methods ...restyClientMethod) (rest.Meta, error) {
	if restyClient == nil {
		return rest.Meta{}, errors.NewError().WithCode(rest.ClientError).WithMessage("Client has not been initialized")

	}
	resp, err := constructRequest(ctx, restyClient, pathParams, query, body, methods...).
		Post(url)
	raw, err := validateResponse(resp, url, helper.CodeSuccess, err)
	if err != nil {
		return raw.Meta, err
	}
	err = dealWithResponse(raw, responseData)
	if err != nil {
		return raw.Meta, err
	}
	return raw.Meta, nil
}

func constructRequest(ctx context.Context, restyClient *resty.Client, pathParams, query map[string]string, body interface{}, methods ...restyClientMethod) *resty.Request {

	req := restyClient.R().
		SetContext(ctx).
		SetResult(&rest.Response{})
	for _, method := range methods {
		req = method(req)
	}
	if body != nil {
		req = req.SetBody(body)
	}
	if query != nil {
		req = req.SetQueryParams(query)
	}
	if pathParams != nil {
		req = req.SetPathParams(pathParams)
	}
	return req
}

func dealWithResponse(raw *rest.Response, responseData interface{}) error {
	if responseData == nil {
		return nil
	}
	json, err := jsoniter.Marshal(raw.Data)
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(json, responseData)
	if err != nil {
		return err
	}
	return nil
}

func validateResponse(resp *resty.Response, url string, restCodeSuccess int, requestErr error) (*rest.Response, error) {
	var raw *rest.Response
	if requestErr != nil {
		raw := &rest.Response{
			Meta: rest.Meta{
				Code: rest.ClientError,
			},
		}
		return raw, requestErr
	}
	if resp.StatusCode() != http.StatusOK {
		raw = &rest.Response{
			Meta: rest.Meta{},
		}
		raw.Meta.Code = resp.StatusCode()
		raw.Meta.Message = string(resp.Body())
		return raw, errors.NewError().WithCode(raw.Meta.Code).WithMessage(raw.Meta.Message)
	}
	raw = resp.Result().(*rest.Response)
	if raw.Meta.Code != restCodeSuccess {
		return raw, errors.NewError().WithCode(raw.Meta.Code).WithMessage(raw.Meta.Message)
	}
	return raw, nil
}
