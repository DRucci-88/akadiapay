package shared

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pageable struct {
	Page int `form:"page"`
	Size int `form:"size"`
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

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return nil, ErrInvalidPaginationParam
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		return nil, ErrInvalidPaginationParam
	}

	pageable := Pageable{Page: page, Size: size}
	pageable.Normalize()
	// log.Printf("%+v\n", pageable)
	return &pageable, nil
}
