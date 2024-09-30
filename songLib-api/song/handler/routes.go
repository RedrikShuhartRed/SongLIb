package handler

import (
	_ "github.com/RedrikShuhartRed/EfMobSongLib/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var r *gin.Engine

var RegisterRoutes = func(handler *handler) *gin.Engine {
	r := gin.Default()
	r.POST("/song", handler.addSong)
	r.GET("/songs", handler.getAllSongs)
	r.GET("/song/:id", handler.getSongByID)
	r.DELETE("/song/:id", handler.deleteSongByID)
	r.PATCH("/song/:id", handler.updateSongByID)
	r.GET("/song/:id/verse", handler.getVersesBySongID)
	r.DELETE("/song/:id/verse", handler.deleteVerseByID)
	r.PATCH("/song/:id/verse", handler.updateVerseBySongID)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
