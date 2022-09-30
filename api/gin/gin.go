package gin

import "github.com/gin-gonic/gin"

const API = "1.0.0"

var engine *gin.Engine

func init() {
	engine = gin.Default()
}
