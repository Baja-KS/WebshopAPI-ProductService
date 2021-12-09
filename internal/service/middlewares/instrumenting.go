package middlewares

import (
	"ProductService/internal/database"
	"ProductService/internal/service"
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"strconv"
	"time"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Next           service.Service
}

func (i *InstrumentingMiddleware) GetByID(ctx context.Context, id uint) (product database.ProductOut,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","GetByID","product_id", strconv.Itoa(int(id)),"error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	product,err=i.Next.GetByID(ctx,id)
	return
}

func (i *InstrumentingMiddleware) Search(ctx context.Context, search string, category uint, minPrice float32, maxPrice float32, discount bool, sortName string, sortPrice string) (products []database.ProductOut, err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Search","product_id", "none","error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	products,err=i.Next.Search(ctx, search, category, minPrice, maxPrice, discount, sortName, sortPrice)
	return
}

func (i *InstrumentingMiddleware) Create(ctx context.Context, data database.ProductIn) (msg string,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Create","product_id", "none","error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	msg,err=i.Next.Create(ctx,data)
	return
}

func (i *InstrumentingMiddleware) Update(ctx context.Context, id uint, data database.ProductIn) (msg string,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Update","product_id", strconv.Itoa(int(id)),"error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	msg,err=i.Next.Update(ctx,id,data)
	return
}

func (i *InstrumentingMiddleware) Delete(ctx context.Context, id uint) (msg string,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Delete","product_id", strconv.Itoa(int(id)),"error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	msg,err=i.Next.Delete(ctx,id)
	return
}
