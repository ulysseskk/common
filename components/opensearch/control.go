package opensearch

import (
	"gitlab.ulyssesk.top/common/common/components/opensearch/model"
	"gitlab.ulyssesk.top/common/common/model/errors"
	"context"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"strings"
)

func (c *Clients) CatIndices(ctx context.Context) (model.CatIndicesResult, model.OpenSearchMeta, error) {
	catIndicesApi := c.GetOpenSearchClient().Cat.Indices
	result := model.CatIndicesResult{}
	resp, err := catIndicesApi(catIndicesApi.WithContext(ctx), catIndicesApi.WithFormat("json"))
	if err != nil {
		return nil, model.OpenSearchMeta{
			HasError: true,
			Error:    err.Error(),
		}, err
	}
	meta, err := c.solveResponse(resp, &result)
	if err != nil {
		return nil, meta, err
	}
	return result, meta, nil
}

func (c *Clients) solveResponse(resp *opensearchapi.Response, result interface{}) (model.OpenSearchMeta, error) {
	if resp.IsError() {
		body := resp.String()
		return model.OpenSearchMeta{
			Code:     resp.StatusCode,
			HasError: true,
			Error:    resp.String(),
		}, errors.NewError().WithCode(errors.OpensearchError).WithMessagef("Fail to get indices from opensearch,Error %s", body)
	}
	err := json.Unmarshal([]byte(strings.TrimPrefix(resp.String(), "[200 OK] ")), result)
	if err != nil {
		return model.OpenSearchMeta{
			Code:     resp.StatusCode,
			HasError: true,
			Error:    err.Error(),
		}, err
	}
	return model.OpenSearchMeta{
		Code:       resp.StatusCode,
		HasError:   false,
		Error:      "",
		HasWarning: false,
		Warnings:   nil,
	}, nil
}
