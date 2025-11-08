package gin

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ExtractPagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10
	pageNumStr := c.Query("pageNum")
	pageSizeStr := c.Query("pageSize")
	pageNumInt, err := strconv.Atoi(pageNumStr)
	if err == nil && pageNumInt > 0 {
		page = pageNumInt
	}
	pageSizeInt, err := strconv.Atoi(pageSizeStr)
	if err == nil && pageSizeInt > 0 {
		pageSize = pageSizeInt
	}
	return page, pageSize
}
