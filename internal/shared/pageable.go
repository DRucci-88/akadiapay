package shared

import (
	"github.com/gin-gonic/gin"
)

type Pageable struct {
	Page int `form:"page,default=1"`
	Size int `form:"size,default=15"`
}

func (p *Pageable) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Size <= 0 {
		p.Size = 10
	}

	if p.Size > 100 {
		p.Size = 100
	}
}

func (p *Pageable) Limit() int {
	return p.Size
}

func (p *Pageable) Offset() int {
	return (p.Page - 1) * p.Size
}

func GetPageable(c *gin.Context) (*Pageable, error) {

	var req Pageable
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, ErrInvalidPaginationParam
	}

	pageable := Pageable{Page: req.Page, Size: req.Size}
	pageable.Normalize()
	// log.Printf("%+v\n", pageable)
	return &pageable, nil
}
