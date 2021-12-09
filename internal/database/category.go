package database

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type CategoryOut struct {
	ID uint `json:"id,omitempty"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	GroupID uint `json:"GroupId"`
	Deletable bool `json:"deletable"`
}

type CategoryServiceResponse struct {
	Categories []CategoryOut `json:"categories"`
}

func GetCategories(GroupID uint,categoryServiceURL string) ([]CategoryOut,error) {

	var categories []CategoryOut
	var response CategoryServiceResponse
	res,err:=http.Get(categoryServiceURL+"/GetByGroupID/"+ strconv.Itoa(int(GroupID)))
	if err != nil {
		return categories,err
	}
	err=json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return categories,err
	}
	return response.Categories,nil
}




