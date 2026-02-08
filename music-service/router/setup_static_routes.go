package router

import "github.com/gin-gonic/gin"

func SetupStatic(r *gin.Engine) {
	// Раздаёт файлы из папки storage по URL /static/...
	// Пример: storage/audio/x.mp3 -> GET /static/storage/audio/x.mp3
	r.Static("/static", ".")
}
