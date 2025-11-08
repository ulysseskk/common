package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/ulysseskk/common/components/opensearch/model"
	"github.com/ulysseskk/common/model/errors"
	"time"
)

func (c *Clients) Search(ctx context.Context, req interface{}, index ...string) (datas *model.SearchResultResp, meta model.OpenSearchMeta, err error) {
	datas, _, meta, err = c.searchScroll(ctx, false, req, index...)
	return datas, meta, err
}

func (c *Clients) SearchScroll(ctx context.Context, req interface{}, index ...string) (datas *model.SearchResultResp, scrollId string, meta model.OpenSearchMeta, err error) {
	return c.searchScroll(ctx, true, req, index...)
}

func (c *Clients) searchScroll(ctx context.Context, scroll bool, req interface{}, index ...string) (datas *model.SearchResultResp, scrollId string, meta model.OpenSearchMeta, err error) {
	searchApi := c.GetOpenSearchClient().Search
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, "", model.OpenSearchMeta{
			HasError: true,
			Error:    err.Error(),
		}, errors.NewError().WithError(err).WithCode(errors.RequestParameterInvalid).WithMessage("Fail to marshal request")
	}
	bodyReq := bytes.NewBuffer(bodyBytes)
	funcs := []func(*opensearchapi.SearchRequest){
		searchApi.WithContext(ctx), searchApi.WithBody(bodyReq), searchApi.WithIndex(index...),
	}
	if scroll {
		funcs = append(funcs, searchApi.WithScroll(5*time.Minute))
	}
	resp, err := searchApi(funcs...)
	if err != nil {
		return nil, "", model.OpenSearchMeta{
			HasError: true,
			Error:    err.Error(),
		}, errors.NewError().WithError(err).WithCode(errors.OpensearchError).WithMessage("Fail to search")
	}
	result := &model.SearchResultResp{}
	meta, err = c.solveResponse(resp, result)
	if err != nil {
		return nil, "", meta, err
	}
	return result, result.ScrollId, meta, nil
}

func (c *Clients) SearchScrollNext(ctx context.Context, scrollId string) (datas *model.SearchResultResp, meta model.OpenSearchMeta, err error) {
	scrollApi := c.GetOpenSearchClient().Scroll
	resp, err := scrollApi(scrollApi.WithContext(ctx), scrollApi.WithScrollID(scrollId), scrollApi.WithScroll(5*time.Minute))
	if err != nil {
		return nil, model.OpenSearchMeta{
			HasError: true,
			Error:    err.Error(),
		}, errors.NewError().WithError(err).WithCode(errors.OpensearchError).WithMessage("Fail to scroll")
	}
	result := &model.SearchResultResp{}
	meta, err = c.solveResponse(resp, result)
	if err != nil {
		return nil, meta, err
	}
	return result, meta, nil
}

func (c *Clients) ClearScroll(ctx context.Context, scrollId string) (meta model.OpenSearchMeta, err error) {
	deleteApi := c.GetOpenSearchClient().ClearScroll
	resp, err := deleteApi(deleteApi.WithContext(ctx), deleteApi.WithScrollID(scrollId))
	if err != nil {
		return model.OpenSearchMeta{
			HasError: true,
			Error:    err.Error(),
		}, errors.NewError().WithError(err).WithCode(errors.OpensearchError).WithMessage("Fail to delete scroll")
	}
	result := map[string]interface{}{}
	meta, err = c.solveResponse(resp, &result)
	if err != nil {
		return meta, err
	}
	return meta, nil
}
