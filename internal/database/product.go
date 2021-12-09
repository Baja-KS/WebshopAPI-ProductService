package database

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

type Product struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name string `gorm:"not null;unique" json:"name"`
	Description string `gorm:"" json:"description,omitempty"`
	Img string `gorm:"" json:"img,omitempty"`
	Price float32 `gorm:"not null" json:"price"`
	Discount int `gorm:"" json:"discount"`
	CategoryID uint `gorm:"not null" json:"categoryID"`
}

type ProductIn struct {
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Img string `json:"img,omitempty"`
	Price float32 `json:"price"`
	Discount int `json:"discount"`
	CategoryID uint `json:"CategoryId"`
}

type ProductOut struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Img string `json:"img,omitempty"`
	Price float32 `json:"price"`
	Discount int `json:"discount"`
	CategoryID uint `json:"CategoryId"`
	Deletable bool `json:"deletable"`
}

func (p *Product) Out(authHeader string) ProductOut {
	return ProductOut{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Img:         p.Img,
		Price:       p.Price,
		Discount:    p.Discount,
		CategoryID:  p.CategoryID,
		Deletable: p.IsDeletable(authHeader),
	}
}

type QuantityOrderedResponse struct {
	Quantity uint `json:"quantity"`
}

func (p *Product) IsDeletable(authHeader string) bool {
	client:=&http.Client{}
	req,err:=http.NewRequest("GET",os.Getenv("ORDER_SERVICE")+"/QuantityOrdered/"+strconv.Itoa(int(p.ID)),nil)
	if err != nil {
		return false
	}
	req.Header.Add("Authorization",authHeader)
	res, err := client.Do(req)
	if err != nil || res.StatusCode!=200 {
		return false
	}
	var response QuantityOrderedResponse
	err=json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return false
	}

	return response.Quantity==0
}

func (i *ProductIn) In() Product {
	return Product{
		Name:        i.Name,
		Description: i.Description,
		Img:         i.Img,
		Price:       i.Price,
		Discount:    i.Discount,
		CategoryID:  i.CategoryID,
	}
}

func ProductArrayOut(models []Product,authHeader string) []ProductOut {
	outArr:=make([]ProductOut,len(models))
	for i,item := range models {
		outArr[i]=item.Out(authHeader)
	}
	return outArr
}

func (p *Product) Update(data ProductIn) Product {
	updated:=*p
	oldImg:=p.Img
	forUpdate:=reflect.ValueOf(data)
	for i:=0;i<forUpdate.NumField();i++ {
		field:=forUpdate.Type().Field(i).Name
		value:=forUpdate.Field(i)
		v := reflect.ValueOf(&updated).Elem().FieldByName(field)
		if v.IsValid() {
			v.Set(value)
		}

	}
	if data.Img=="" || data.Img=="undefined" || data.Img=="null" {
		updated.Img=oldImg
	}
	return updated
}

func UploadImage(file *multipart.File,authHeader string,imgName string) error {
	var reqBody bytes.Buffer
	multiPartWriter:=multipart.NewWriter(&reqBody)
	fileWriter,err:=multiPartWriter.CreateFormFile("img",imgName)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, *file)
	if err != nil {
		return err
	}

	_=multiPartWriter.Close()

	req,err:=http.NewRequest("POST",os.Getenv("IMG_SERVICE")+"/Upload",&reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type",multiPartWriter.FormDataContentType())
	req.Header.Set("Authorization",authHeader)

	client:=&http.Client{}
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil

}

func DecodeMultipartRequest(r *http.Request) (ProductIn,error) {
	err:=r.ParseMultipartForm(32<<20)
	if err != nil {
		//logger.Log("status","successn't","place","Parse-Multipart-Form","err",err)
		return ProductIn{}, err
	}
	//text fields
	price,err:=strconv.ParseFloat(r.FormValue("price"),32)
	if err != nil {
		return ProductIn{}, err
	}
	discount,err:=strconv.Atoi(r.FormValue("discount"))
	if err != nil {
		discount=0
	}
	category,err:=strconv.ParseUint(r.FormValue("CategoryId"),10,32)
	if err != nil {
		return ProductIn{}, err
	}
	data:=ProductIn{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       float32(price),
		Discount:    discount,
		CategoryID:  uint(category),
	}
	file, hdr,err:=r.FormFile("img")
	if file == nil || err != nil {
		return data, nil
	}

	imgName:=RandomString(32)+filepath.Ext(hdr.Filename)

	err = UploadImage(&file,r.Header["Authorization"][0],imgName)
	if err != nil {
		_ = file.Close()
		return data, err
	}


	err = file.Close()
	if err != nil {
		return data, err
	}

	data.Img=imgName
	return data,nil
}
