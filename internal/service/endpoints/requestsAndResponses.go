package endpoints

import (
	"ProductService/internal/database"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ParseIDFromURL(r *http.Request) (uint, error) {
	params:=mux.Vars(r)
	idStr:=params["id"]
	id,err:=strconv.ParseUint(idStr,10,32)
	if err != nil {
		return 0,err
	}
	return uint(id),nil
}



type GetByIDRequest struct {
	ID uint `json:"id,omitempty"`
}

type GetByIDResponse struct {
	Product database.ProductOut `json:"product"`
}
type SearchRequest struct {
	Search string `json:"search"`
	MinPrice float32 `json:"minPrice"`
	MaxPrice float32 `json:"maxPrice"`
	Discount bool `json:"discount"`
	CategoryID uint `json:"CategoryID"`
	SortName string `json:"sortName"`
	SortPrice string `json:"sortPrice"`
}

type SearchResponse struct {
	Products []database.ProductOut `json:"products"`
}
type CreateRequest struct {
	Data database.ProductIn `json:"data"`
}

type CreateResponse struct {
	Message string `json:"message"`
}
type UpdateRequest struct {
	ID uint `json:"id,omitempty"`
	Data database.ProductIn `json:"data"`
}

type UpdateResponse struct {
	Message string `json:"message"`
}
type DeleteRequest struct {
	ID uint `json:"id,omitempty"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}


func DecodeGetByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request GetByIDRequest
	id,err:=ParseIDFromURL(r)
	if err != nil {
		return request,err
	}
	request.ID=id
	return request,nil
}
func DecodeSearchRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request SearchRequest
	request.Search = r.URL.Query().Get("search")
	request.Discount = r.URL.Query().Get("withDiscount") == "true"
	minPriceParam:=r.URL.Query().Get("minPrice")
	maxPriceParam:=r.URL.Query().Get("maxPrice")
	categoryParam:=r.URL.Query().Get("CategoryId")
	sortNameParam:=r.URL.Query().Get("sortName")
	sortPriceParam:=r.URL.Query().Get("sortPrice")
	if minPriceParam!="" {
		price,err:=strconv.ParseFloat(minPriceParam,32)
		if err == nil {
			request.MinPrice=float32(price)
		} else {
			request.MinPrice=-1
		}
	} else {
		request.MinPrice=-1
	}
	if maxPriceParam!="" {
		price,err:=strconv.ParseFloat(maxPriceParam,32)
		if err == nil {
			request.MaxPrice=float32(price)
		} else {
			request.MaxPrice=-1
		}
	}  else {
		request.MaxPrice=-1
	}
	if categoryParam!=""{
		category,err:=strconv.ParseUint(categoryParam,10,32)
		if err == nil {
			request.CategoryID=uint(category)
		} else {
			request.CategoryID=0
		}
	}else {
		request.CategoryID=0
	}
	request.SortName=sortNameParam
	request.SortPrice=sortPriceParam
	//request.Sort=sortParam
	return request,nil
}
func DecodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request CreateRequest
	data,err:=database.DecodeMultipartRequest(r)
	request.Data=data
	if err != nil {
		return request, err
	}
	return request,nil
}
func DecodeUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request UpdateRequest
	id,err:=ParseIDFromURL(r)
	if err != nil {
		return request,err
	}
	request.ID=id
	data,err:=database.DecodeMultipartRequest(r)
	request.Data=data
	if err != nil {
		return request, err
	}
	return request,nil

}
func DecodeDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request DeleteRequest
	id,err:=ParseIDFromURL(r)
	if err != nil {
		return request,err
	}
	request.ID=id
	return request,nil
}


func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type","application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(response)
}