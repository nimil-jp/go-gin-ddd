package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Paging struct {
	page   int
	length int

	offset int
	limit  int
}

func NewPaging(c *gin.Context) *Paging {
	paging := new(Paging)
	paging.page, _ = strconv.Atoi(c.Query("page"))
	if paging.page == 0 {
		paging.page = 1
	}
	paging.length, _ = strconv.Atoi(c.Query("length"))

	paging.offset, _ = strconv.Atoi(c.Query("offset"))
	paging.limit, _ = strconv.Atoi(c.Query("limit"))
	return paging
}

func (p *Paging) GetCount(query *gorm.DB, model interface{}) (*gorm.DB, uint, error) {
	var count int64

	copiedQuery := &gorm.DB{
		Config:       query.Config,
		Error:        query.Error,
		RowsAffected: query.RowsAffected,
		Statement:    query.Statement,
	}

	if err := copiedQuery.Model(model).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return query, uint(count), nil
}

func (p *Paging) Query() func(*gorm.DB) *gorm.DB {
	offset := (p.page - 1) * p.length
	limit := p.length

	if p.offset > 0 {
		offset = p.offset
	}

	if p.limit > 0 {
		limit = p.limit
	}

	return func(db *gorm.DB) *gorm.DB {
		db.Offset(offset).Limit(limit)
		return db
	}
}
