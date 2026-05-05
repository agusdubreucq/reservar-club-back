package main

import (
	"fmt"
	"log"

	"github.com/adubr/reservar-club-back/internal/config"
	"github.com/adubr/reservar-club-back/internal/db"
	"github.com/adubr/reservar-club-back/internal/handlers"
	"github.com/adubr/reservar-club-back/internal/repository/postgres"
	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/adubr/reservar-club-back/pkg/jwt"
	"github.com/adubr/reservar-club-back/pkg/oauth2"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	dbConn, err := db.InitPostgres(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	userRepo := postgres.NewUserRepository(dbConn)
	sportRepo := postgres.NewSportRepository(dbConn)
	courtRepo := postgres.NewCourtRepository(dbConn)
	reservationRepo := postgres.NewReservationRepository(dbConn)

	authService := services.NewAuthService(userRepo)
	sportService := services.NewSportService(sportRepo)
	courtService := services.NewCourtService(courtRepo, sportRepo)
	reservationService := services.NewReservationService(reservationRepo, courtRepo)

	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret)
	googleProvider := oauth2.NewGoogleOAuth2Provider(
		cfg.OAuth2.ClientID,
		cfg.OAuth2.ClientSecret,
		cfg.OAuth2.RedirectURL,
	)

	authHandler := handlers.NewAuthHandler(authService, tokenManager, googleProvider)
	sportHandler := handlers.NewSportHandler(sportService)
	courtHandler := handlers.NewCourtHandler(courtService)
	reservationHandler := handlers.NewReservationHandler(reservationService)

	router := gin.Default()

	router.Use(handlers.CORSMiddleware())

	publicRoutes := router.Group("")
	{
		publicRoutes.GET("/auth/login", authHandler.GetAuthURL)
		publicRoutes.POST("/auth/login", authHandler.Login)
		publicRoutes.GET("/sports", sportHandler.ListSports)
		publicRoutes.GET("/sports/:id", sportHandler.GetSport)
		publicRoutes.GET("/courts", courtHandler.ListCourts)
		publicRoutes.GET("/courts/:id", courtHandler.GetCourt)
		publicRoutes.GET("/sports/:sport_id/courts", courtHandler.ListCourtsBySport)
	}

	protectedRoutes := router.Group("")
	protectedRoutes.Use(handlers.AuthMiddleware(tokenManager))
	{
		protectedRoutes.GET("/auth/me", authHandler.GetMe)

		protectedRoutes.POST("/sports", sportHandler.CreateSport)
		protectedRoutes.DELETE("/sports/:id", sportHandler.DeleteSport)

		protectedRoutes.POST("/courts", courtHandler.CreateCourt)
		protectedRoutes.PUT("/courts/:id", courtHandler.UpdateCourt)
		protectedRoutes.DELETE("/courts/:id", courtHandler.DeleteCourt)

		protectedRoutes.GET("/reservations", reservationHandler.ListUserReservations)
		protectedRoutes.GET("/reservations/:id", reservationHandler.GetReservation)
		protectedRoutes.POST("/reservations", reservationHandler.CreateReservation)
		protectedRoutes.POST("/reservations/:id/cancel", reservationHandler.CancelReservation)
	}

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server started at %s", address)
	if err := router.Run(address); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
