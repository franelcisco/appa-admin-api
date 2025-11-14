package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"appa_admin_api/internal/config"
	"appa_admin_api/internal/handlers"
	routes "appa_admin_api/internal/routers"
	"appa_admin_api/internal/services"
	"appa_admin_api/pkg/bcv"
	"appa_admin_api/pkg/db"
	"appa_admin_api/pkg/firebase"
	"appa_admin_api/pkg/logs"
	"appa_admin_api/pkg/middleware"
	"appa_admin_api/pkg/shopify"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	if p := os.Getenv("PORT"); p != "" {
		cfg.Port = p
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.Debug == "" {
		cfg.Debug = "0"
	}

	logger := logs.NewZapLogger(cfg.Debug == "1")
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("error syncing logger: %v\n", err)
		}
	}()

	sslmode := cfg.SSLMode
	fmt.Printf("sslmode -> %s\n", sslmode)
	if len(sslmode) > 0 {
		sslmode = "sslmode=" + sslmode
	}

	//connect the database
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s %s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslmode)
	// gorm connect
	gormDB, err := db.NewDBSQLHandler(connStr)
	if err != nil {
		logger.Error(err.Error(), zap.Any("host", cfg.DBHost), zap.Any("port", cfg.DBPort), zap.Any("user", cfg.DBUser), zap.Any("dbname", cfg.DBName))
	}

	db, err := gormDB.DB()
	if err != nil {
		logger.Error(err.Error(), zap.Any("host", cfg.DBHost), zap.Any("port", cfg.DBPort), zap.Any("user", cfg.DBUser), zap.Any("dbname", cfg.DBName))
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("error db body: %v\n", err)
		}
	}()

	loc, err := time.LoadLocation("America/Caracas")
	if err != nil {
		logger.Fatal("could not load Venezuela time zone", zap.Error(err))
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) { c.String(200, "ok") })
	router.GET("/_ah/health", func(c *gin.Context) { c.String(200, "ok") })

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins, // es necesario con las cookies
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true, // <- clave para cookies
		MaxAge:           12 * time.Hour,
	}))
	logger.Info("Allow Origins", zap.Any("origins", cfg.CORSAllowedOrigins))

	// init resources
	shopifyRepo := shopify.NewRepository(
		cfg.ShopifyStoreName, cfg.ShopifyAPIVersion, cfg.ShopifyAdminToken, logger,
	)

	firebaseApp, err := firebase.NewRepository(
		logger, cfg.FirebaseServiceAccountID, cfg.FirebaseProjectID,
	)
	if err != nil {
		logger.Error(err.Error(), zap.Any("service_account_id", cfg.FirebaseServiceAccountID), zap.Any("project_id", cfg.FirebaseProjectID))
	}

	bcvClient, err := bcv.NewClient(loc, logger)
	if err != nil {
		logger.Error("error initializing BCV client", zap.Error(err))
	}

	// init middlewares
	AuthMiddleware, err := middleware.NewAuthMiddleware(firebaseApp)
	if err != nil {
		logger.Error(err.Error())
	}

	// init Services
	orderService := services.NewOrdersService(gormDB, shopifyRepo, logger)

	// init handlers
	orderHandler := handlers.NewOrdersHandler(orderService, bcvClient)

	// init routes
	orderRoutes := routes.NewOrdersRoutes(orderHandler, AuthMiddleware)

	// set routes
	orderRoutes.SetRouter(router)

	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatal(err.Error(), zap.Any("port", cfg.Port))
	}
}
