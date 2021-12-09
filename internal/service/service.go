package service

import (
	"ProductService/internal/database"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//ProductService should implement the Service interface


type ProductService struct {
	DB *gorm.DB
}

func ValidateCategory(categoryServiceURL string, id uint) bool {
	_,err:=http.Get(categoryServiceURL+"/GetByID/"+ strconv.Itoa(int(id)))
	if err != nil {
		return false
	}
	return true
}

type Service interface {
	GetByID(ctx context.Context,id uint) (database.ProductOut,error)
	Search(ctx context.Context, search string, category uint, minPrice float32, maxPrice float32, discount bool, sortName string, sortPrice string) ([]database.ProductOut, error)
	Create(ctx context.Context,data database.ProductIn) (string,error)
	Update(ctx context.Context,id uint,data database.ProductIn) (string,error)
	Delete(ctx context.Context,id uint) (string,error)
}

func (p *ProductService) GetByID(ctx context.Context, id uint) (database.ProductOut, error) {
	token:=ctx.Value("auth").(string)
	authHeader:=fmt.Sprintf("Bearer %s",token)
	var product database.Product
	p.DB.Where("id = ?",id).First(&product)
	return product.Out(authHeader),nil
}

func (p *ProductService) Search(ctx context.Context, search string, category uint, minPrice float32, maxPrice float32, discount bool, sortName string, sortPrice string) ([]database.ProductOut, error) {
	var products []database.Product


	token:=ctx.Value("auth").(string)
	authHeader:=fmt.Sprintf("Bearer %s",token)


	result:=p.DB.Where(
		p.DB.Where("cast(id as varchar) ilike ?","%"+search+"%").Or("name ilike ?","%"+search+"%"),
	)
	if minPrice != -1 {
		result=result.Where("price >= ?",minPrice)
	}
	if maxPrice != -1 {
		result=result.Where("price <= ?",maxPrice)
	}
	if category != 0 {
		result=result.Where("category_id = ?",category)
	}
	if discount  {
		result=result.Where("discount > ?","0")
	}
	if len(sortName)>0 {
		result=result.Order("name "+strings.ToLower(sortName))
	}
	if len(sortPrice) > 0 {
		result=result.Order("price "+strings.ToLower(sortPrice))
	}
	if result.Debug().Find(&products).Error != nil {
		return database.ProductArrayOut(products,authHeader),result.Error
	}
	out:=database.ProductArrayOut(products,authHeader)
	return out,nil
}

func (p *ProductService) Create(ctx context.Context, data database.ProductIn) (string, error) {
	product:=data.In()
	if !ValidateCategory(os.Getenv("CATEGORY_SERVICE"),product.CategoryID) {
		return "Non existent category",errors.New("category with that ID doesnt exist")
	}
	result:=p.DB.Create(&product)
	if result.Error != nil {
		return "Error", result.Error
	}
	return "Successfully created", nil
}

func (p *ProductService) Update(ctx context.Context, id uint, data database.ProductIn) (string, error) {
	var product database.Product
	notFound:=p.DB.Where("id = ?",id).First(&product).Error
	if notFound != nil {
		return "That product doesn't exist", notFound
	}
	if !ValidateCategory(os.Getenv("CATEGORY_SERVICE"),product.CategoryID) {
		return "Non existent category",errors.New("category with that ID doesnt exist")
	}
	product=product.Update(data)
	err:=p.DB.Save(&product).Error
	if err != nil {
		return "Error updating product", err
	}

	return "Product updated successfully", nil
}

func (p *ProductService) Delete(ctx context.Context, id uint) (string, error) {
	var product database.Product
	notFound:=p.DB.Where("id = ?",id).First(&product).Error
	if notFound != nil {
		return "That product doesn't exist", notFound
	}
	err:=p.DB.Delete(&database.Product{},id).Error
	if err != nil {
		return "Error deleting product", err
	}

	return "Product deleted successfully", nil
}
