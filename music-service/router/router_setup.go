package router

import (
	"music-service/config"
	"music-service/handlers"
	"music-service/repositories"
	"music-service/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) error {

	// --- Auth (пока живёт в router) ---
	r.POST("/auth/signup", CreateUser)
	r.POST("/auth/login", Login)

	// --- User profile (быстрый endpoint) ---
	r.GET("/user/profile", CheckAuth, GetUserProfile)

	// --- Music ---
	musicRepo := repositories.NewMusicRepository(repositories.SQLDB)
	musicSvc := services.NewMusicService(musicRepo)
	musicHandler := handlers.NewMusicHandler(musicSvc)

	r.GET("/artists", musicHandler.ListArtists)
	r.GET("/artists/:id", musicHandler.GetArtistByID)
	r.GET("/tracks/search", musicHandler.SearchTracks)

	// --- Playlists ---
	playlistRepo := repositories.NewPlaylistRepository(repositories.SQLDB)
	playlistSvc := services.NewPlaylistService(playlistRepo)
	playlistHandler := handlers.NewPlaylistHandler(playlistSvc)

	r.POST("/playlists", CheckAuth, playlistHandler.Create)
	r.GET("/users/me/playlists", CheckAuth, playlistHandler.MyPlaylists)
	r.GET("/playlists/:id", playlistHandler.GetByID)
	r.PATCH("/playlists/:id", CheckAuth, playlistHandler.Rename)
	r.POST("/playlists/:id/tracks", CheckAuth, playlistHandler.AddTrack)
	r.DELETE("/playlists/:id/tracks/:trackId", CheckAuth, playlistHandler.RemoveTrack)

	// --- Users ---
	userRepo := repositories.NewUserRepository(repositories.SQLDB)
	userSvc := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	r.GET("/users/me", CheckAuth, userHandler.Me)
	r.PATCH("/users/me", CheckAuth, userHandler.UpdateMe)
	r.PATCH("/users/me/password", CheckAuth, userHandler.ChangePassword)
	r.GET("/users", userHandler.List)

	// --- Artist profile ---
	artistRepo := repositories.NewArtistRepository(repositories.SQLDB)
	artistSvc := services.NewArtistService(artistRepo)
	artistHandler := handlers.NewArtistHandler(artistSvc)

	r.GET("/artist/me", CheckAuth, RequireRole("artist"), artistHandler.Me)
	r.POST("/artist/profile", CheckAuth, RequireRole("artist"), artistHandler.CreateProfile)

	// --- Uploads (artist only) ---
	cfg := config.Load()
	storageSvc := services.NewStorageService(cfg)
	uploadHandler := handlers.NewUploadHandler(storageSvc)

	r.POST("/uploads/audio", CheckAuth, RequireRole("artist"), uploadHandler.UploadAudio)
	r.POST("/uploads/covers", CheckAuth, RequireRole("artist"), uploadHandler.UploadCover)

	// --- Tracks ---
	trackRepo := repositories.NewTrackRepository(repositories.SQLDB)
	trackSvc := services.NewTrackService(trackRepo)
	trackHandler := handlers.NewTrackHandler(trackSvc, artistRepo)

	r.POST("/artist/tracks", CheckAuth, RequireRole("artist"), trackHandler.CreateAsArtist)
	r.GET("/tracks/:id", trackHandler.GetByID)
	r.GET("/tracks/top", trackHandler.Top)
	r.GET("/tracks", trackHandler.List)
	r.GET("/artists/:id/tracks", trackHandler.ListByArtist)
	r.GET("/albums/:id/tracks", trackHandler.ListByAlbum)

	// --- Listens / History ---
	listenRepo := repositories.NewListenRepository(repositories.SQLDB)
	listenSvc := services.NewListenService(listenRepo)
	listenHandler := handlers.NewListenHandler(listenSvc)

	r.POST("/tracks/:id/listen", CheckAuth, listenHandler.Record)
	r.GET("/users/me/history", CheckAuth, listenHandler.MyHistory)

	// --- Likes ---
	likeRepo := repositories.NewLikeRepository(repositories.SQLDB)
	likeSvc := services.NewLikeService(likeRepo)
	likeHandler := handlers.NewLikeHandler(likeSvc)

	r.POST("/tracks/:id/like", CheckAuth, likeHandler.Like)
	r.DELETE("/tracks/:id/like", CheckAuth, likeHandler.Unlike)
	r.GET("/users/me/likes", CheckAuth, likeHandler.MyLikes)

	// --- Search (public) ---
	searchRepo := repositories.NewSearchRepository(repositories.SQLDB)
	searchSvc := services.NewSearchService(searchRepo)
	searchHandler := handlers.NewSearchHandler(searchSvc)

	r.GET("/search", searchHandler.Search)

	// --- Genres ---
	genreRepo := repositories.NewGenreRepository(repositories.SQLDB)
	genreSvc := services.NewGenreService(genreRepo)
	genreHandler := handlers.NewGenreHandler(genreSvc)

	r.GET("/genres", genreHandler.List)
	r.POST("/genres", CheckAuth, RequireRole("artist"), genreHandler.Create)

	// --- Albums ---
	albumRepo := repositories.NewAlbumRepository(repositories.SQLDB)
	albumSvc := services.NewAlbumService(albumRepo)
	albumHandler := handlers.NewAlbumHandler(albumSvc, artistRepo)

	r.POST("/artist/albums", CheckAuth, RequireRole("artist"), albumHandler.CreateAsArtist)
	r.GET("/albums/:id", albumHandler.GetByID)
	r.GET("/artists/:id/albums", albumHandler.ListByArtist)

	return nil
}
