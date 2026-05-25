package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTP         HTTPConfig
	Log          LogConfig
	OTel         OTelConfig
	ServiceURLs  ServiceURLsConfig
	RateLimit    RateLimitConfig
	CORS         CORSConfig
	JWKS         JWKSConfig
	UpstreamTimeout time.Duration
}

type HTTPConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

type OTelConfig struct {
	Enabled     bool
	Exporter    string
	Endpoint    string
	ServiceName string
}

type ServiceURLsConfig struct {
	AuthURL    string
	CatalogURL string
	CartURL    string
	OrderURL   string
}

type RateLimitConfig struct {
	RPS int
}

type CORSConfig struct {
	AllowedOrigins []string
}

type JWKSConfig struct {
	RefreshInterval time.Duration
}

func Load() (*Config, error) {
	authURL, err := requireEnv("AUTH_SERVICE_URL")
	if err != nil {
		return nil, err
	}
	catalogURL, err := requireEnv("CATALOG_SERVICE_URL")
	if err != nil {
		return nil, err
	}
	cartURL, err := requireEnv("CART_SERVICE_URL")
	if err != nil {
		return nil, err
	}
	orderURL, err := requireEnv("ORDER_SERVICE_URL")
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTP: HTTPConfig{
			Addr:         getEnv("HTTP_ADDR", ":80"),
			ReadTimeout:  getEnvDuration("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("HTTP_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		OTel: OTelConfig{
			Enabled:     getEnvBool("OTEL_ENABLED", false),
			Exporter:    getEnv("OTEL_EXPORTER", "stdout"),
			Endpoint:    getEnv("OTEL_ENDPOINT", ""),
			ServiceName: getEnv("OTEL_SERVICE_NAME", "api-gateway"),
		},
		ServiceURLs: ServiceURLsConfig{
			AuthURL:    authURL,
			CatalogURL: catalogURL,
			CartURL:    cartURL,
			OrderURL:   orderURL,
		},
		RateLimit: RateLimitConfig{
			RPS: getEnvInt("RATE_LIMIT_RPS", 100),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvStringSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		},
		JWKS: JWKSConfig{
			RefreshInterval: getEnvDuration("JWKS_REFRESH_INTERVAL", 5*time.Minute),
		},
		UpstreamTimeout: getEnvDuration("UPSTREAM_TIMEOUT", 10*time.Second),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func getEnvStringSlice(key string, fallback []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parts := strings.Split(v, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			result = append(result, s)
		}
	}
	return result
}

func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("required environment variable %q is not set", key)
	}
	return v, nil
}
