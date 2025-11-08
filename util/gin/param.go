package gin

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetParamInt32(c *gin.Context, name string) (int32, error) {
	param := c.Param(name)
	if param == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return int32(value), nil
}
