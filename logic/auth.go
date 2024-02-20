package logic

import (
	"github.com/gin-gonic/gin"
	"sre/jmx_exporter/config"
)

func Auth(c *gin.Context) {
	token := c.Request.URL.Query().Get("token")
	if token != config.C.GetString("auth.token") {
		c.JSON(401, gin.H{
			"msg": "没有接口权限",
		})
		return
	}
	c.Next()
}
