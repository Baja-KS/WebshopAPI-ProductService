package endpoints

import (
	"ProductService/internal/service"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type EndpointSet struct {
	GetByIDEndpoint endpoint.Endpoint
	SearchEndpoint endpoint.Endpoint
	CreateEndpoint endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc service.Service) EndpointSet {
	return EndpointSet{
		GetByIDEndpoint:    MakeGetByIDEndpoint(svc),
		SearchEndpoint:    MakeSearchEndpoint(svc),
		CreateEndpoint:    MakeCreateEndpoint(svc),
		UpdateEndpoint:    MakeUpdateEndpoint(svc),
		DeleteEndpoint:    MakeDeleteEndpoint(svc),
	}
}

func MakeGetByIDEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        req:=request.(GetByIDRequest)
		product,err:=svc.GetByID(ctx,req.ID)
		if err != nil {
			return nil, err
		}
		return GetByIDResponse{Product: product},nil
	}
}
func MakeSearchEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        req:=request.(SearchRequest)
		products,err:=svc.Search(ctx, req.Search, req.CategoryID, req.MinPrice, req.MaxPrice, req.Discount, req.SortName, req.SortPrice)
		if err != nil {
			return nil, err
		}
		return SearchResponse{Products: products},nil
	}
}
func MakeCreateEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        req:=request.(CreateRequest)
		msg,err:=svc.Create(ctx,req.Data)
		if err != nil {
			return nil, err
		}
		return CreateResponse{Message: msg},nil
	}
}
func MakeUpdateEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        req:=request.(UpdateRequest)
		msg,err:=svc.Update(ctx,req.ID,req.Data)
		if err != nil {
			return nil, err
		}
		return UpdateResponse{Message: msg},nil
	}
}
func MakeDeleteEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        req:=request.(DeleteRequest)
		msg,err:=svc.Delete(ctx,req.ID)
		if err != nil {
			return nil, err
		}
		return DeleteResponse{Message: msg},nil
	}
}
