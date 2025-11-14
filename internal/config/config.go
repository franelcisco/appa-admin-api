package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Port  string
	Debug string

	// DATABASE
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string

	// Firebase
	FirebaseServiceAccountID string
	FirebaseProjectID        string

	// Shopify API credentials
	ShopifyAPIVersion string
	ShopifyAdminToken string
	ShopifyStoreName  string
	ShopifyHMACSecret string

	// CORS
	CORSAllowedOrigins []string

	// R4 API credentials
	R4EntryPoint   string
	R4APIEcommerce string
	R4Secret       string
}

// Load reads configuration from environment variables and returns a Config struct
func Load() (*Config, error) {

	var (
		corsOrigins []string
	)

	if result := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","); len(result) > 0 {
		corsOrigins = result
	}

	cfg := &Config{
		Port:  os.Getenv("PORT"),
		Debug: os.Getenv("DEBUG"),

		// DATABASE
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		SSLMode:    os.Getenv("SSL_MODE"),

		FirebaseServiceAccountID: os.Getenv("FIREBASE_SERVICE_ACCOUNT_ID"),
		FirebaseProjectID:        os.Getenv("FIREBASE_PROJECT_ID"),

		ShopifyAPIVersion: os.Getenv("SHOPIFY_API_VERSION"),
		ShopifyAdminToken: os.Getenv("SHOPIFY_ADMIN_TOKEN"),
		ShopifyStoreName:  os.Getenv("SHOPIFY_STORE_NAME"),
		ShopifyHMACSecret: os.Getenv("SHOPIFY_HMAC_SECRET"),

		CORSAllowedOrigins: corsOrigins,

		R4EntryPoint:   os.Getenv("R4_ENTRY_POINT"),
		R4APIEcommerce: os.Getenv("R4_API_ECOMMERCE"),
		R4Secret:       os.Getenv("R4_SECRET"),
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validate(cfg *Config) error {

	if cfg.DBHost == "" {
		return fmt.Errorf("DBHost is not configured")
	}
	if cfg.DBPort == "" {
		return fmt.Errorf("DBPort is not configured")
	}
	if cfg.DBUser == "" {
		return fmt.Errorf("DBUser is not configured")
	}
	if cfg.DBPassword == "" {
		return fmt.Errorf("DBPassword is not configured")
	}
	if cfg.DBName == "" {
		return fmt.Errorf("DBName is not configured")
	}

	if len(cfg.CORSAllowedOrigins) == 0 {
		return fmt.Errorf("CORSAllowedOrigins is not configured")
	}

	if cfg.FirebaseServiceAccountID == "" {
		return fmt.Errorf("FirebaseServiceAccountID is not configured")
	}
	if cfg.FirebaseProjectID == "" {
		return fmt.Errorf("FirebaseProjectID is not configured")
	}

	if cfg.ShopifyAPIVersion == "" {
		return fmt.Errorf("ShopifyAPIVersion is not configured")
	}
	if cfg.ShopifyAdminToken == "" {
		return fmt.Errorf("ShopifyAdminToken is not configured")
	}
	if cfg.ShopifyStoreName == "" {
		return fmt.Errorf("ShopifyStoreName is not configured")
	}
	if cfg.ShopifyHMACSecret == "" {
		return fmt.Errorf("ShopifyHMACSecret is not configured")
	}

	if cfg.R4EntryPoint == "" {
		return fmt.Errorf("R4EntryPoint is not configured")
	}
	if cfg.R4APIEcommerce == "" {
		return fmt.Errorf("R4APIEcommerce is not configured")
	}
	if cfg.R4Secret == "" {
		return fmt.Errorf("R4 Secret is not configured")
	}

	return nil
}
