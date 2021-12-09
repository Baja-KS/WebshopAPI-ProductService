package middlewares

import (
	"ProductService/internal/database"
	"context"
	"time"

	//import the service package
	"ProductService/internal/service"
	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   service.Service
}

func (l *LoggingMiddleware) GetByID(ctx context.Context, id uint) (product database.ProductOut,err error) {
	defer func(begin time.Time) {
		err := l.Logger.Log("method", "get by id", "id",id ,"name", product.Name,"err", err, "took", time.Since(begin))
		if err != nil {
			return
		}
	}(time.Now())
	product,err=l.Next.GetByID(ctx,id)
	return
}

func (l *LoggingMiddleware) Search(ctx context.Context, search string, category uint, minPrice float32, maxPrice float32, discount bool, sortName string, sortPrice string) (products []database.ProductOut, err error) {
	defer func(begin time.Time) {
		err := l.Logger.Log("method", "search", "products", len(products),"err", err, "took", time.Since(begin))
		if err != nil {
			return
		}
	}(time.Now())
	products,err=l.Next.Search(ctx, search, category, minPrice, maxPrice, discount, sortName, sortPrice)
	return
}

func (l *LoggingMiddleware) Create(ctx context.Context, data database.ProductIn) (msg string,err error) {
	defer func(begin time.Time) {
		err := l.Logger.Log("method", "create", "message", msg,"err", err, "took", time.Since(begin))
		if err != nil {
			return
		}
	}(time.Now())
	msg,err=l.Next.Create(ctx,data)
	return
}

func (l *LoggingMiddleware) Update(ctx context.Context, id uint, data database.ProductIn) (msg string,err error) {
	defer func(begin time.Time) {
		err := l.Logger.Log("method", "update", "id",id ,"message", msg,"err", err, "took", time.Since(begin))
		if err != nil {
			return
		}
	}(time.Now())
	msg,err=l.Next.Update(ctx,id,data)
	return
}

func (l *LoggingMiddleware) Delete(ctx context.Context, id uint) (msg string,err error) {
	defer func(begin time.Time) {
		err := l.Logger.Log("method", "delete", "id",id ,"message", msg,"err", err, "took", time.Since(begin))
		if err != nil {
			return
		}
	}(time.Now())
	msg,err=l.Next.Delete(ctx,id)
	return
}


