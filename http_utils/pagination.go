package http_utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func Paginate[T any](items []T, page int, pageSize int) []T {
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(items) {
		start = len(items)
	}
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

func GetPageAndPageSize(ctx *gin.Context) (int, int) {
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("pageSize", "10")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 {
		pageSizeInt = 10
	}

	return pageInt, pageSizeInt
}
