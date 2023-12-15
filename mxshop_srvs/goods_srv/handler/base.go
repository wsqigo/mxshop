package handler

import (
	"golang.org/x/exp/constraints"
	"gorm.io/gorm"
)

func Paginate[T constraints.Integer](pageNum, pageSize T) func(db *gorm.DB) *gorm.DB {
	num := int(pageNum)
	size := int(pageSize)

	return func(db *gorm.DB) *gorm.DB {
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (num - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
