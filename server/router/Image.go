package router

import (
	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	// 暂时注释掉图像识别路由，因为它依赖于 ONNX 库，在 Windows 上存在问题
	// r.POST("/recognize", image.RecognizeImage)
}
