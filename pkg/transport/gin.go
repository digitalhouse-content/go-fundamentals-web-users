package transport

import (
	"github.com/gin-gonic/gin"
)

func GinServer(
	endpoint Endpoint,
	decode func(c *gin.Context) (interface{}, error),
	encode func(c *gin.Context, resp interface{}),
	encondeError func(c *gin.Context, err error)) func(c *gin.Context) {

	return func(c *gin.Context) {
		data, err := decode(c)
		if err != nil {
			encondeError(c, err)
			return
		}

		res, err := endpoint(c.Request.Context(), data)
		if err != nil {
			encondeError(c, err)
			return
		}

		encode(c, res)
	}

}
