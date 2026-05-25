package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sharederrors "api-on/internal/shared/errors"
)

type Config struct {
	AppEnv              string
	Port                string
	SecretKey           string
	JWTIssuer           string
	FrontendURL         string
	PasswordResetTTLMin int
	SMTPHost            string
	SMTPPort            string
	SMTPUsername        string
	SMTPPassword        string
	SMTPFrom            string
	StorageDriver       string
	DataFile            string
	DatabaseURL         string
	DatabaseAutoMigrate bool
	DatabaseMaxConns    int32
	CloudProvider       string
	AWSRegion           string
	AWSS3Bucket         string
	AWSSecretsPrefix    string
	AWSUseIAMRole       bool
	AWSCloudWatchNS     string
}

func Load() (*Config, error) {
	dotEnvPath, err := findDotEnvPath()
	if err != nil {
		return nil, err
	}

	if dotEnvPath != "" {
		if err := loadDotEnv(dotEnvPath); err != nil {
			return nil, err
		}
	}

	if err := loadDotEnv(filepath.Join(".", ".env")); err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		port = "8080"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if strings.TrimSpace(issuer) == "" {
		issuer = "psicosistem-backend"
	}

	dataFile := os.Getenv("DATA_FILE")
	if strings.TrimSpace(dataFile) == "" {
		dataFile = filepath.Join("data", "app_state.json")
	}

	storageDriver := strings.ToLower(strings.TrimSpace(envOrDefault("STORAGE_DRIVER", "json")))
	if storageDriver != "json" && storageDriver != "postgres" {
		return nil, sharederrors.Invalid("INVALID_STORAGE_DRIVER", "STORAGE_DRIVER must be json or postgres")
	}

	secret := os.Getenv("JWT_SECRET")
	if strings.TrimSpace(secret) == "" {
		secret = os.Getenv("SECRET_KEY")
	}
	if strings.TrimSpace(secret) == "" {
		return nil, sharederrors.Internal("JWT_SECRET or SECRET_KEY is required")
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if storageDriver == "postgres" && databaseURL == "" {
		return nil, sharederrors.Internal("DATABASE_URL is required when STORAGE_DRIVER=postgres")
	}

	return &Config{
		AppEnv:              envOrDefault("APP_ENV", "development"),
		Port:                port,
		SecretKey:           secret,
		JWTIssuer:           issuer,
		FrontendURL:         strings.TrimRight(strings.TrimSpace(envOrDefault("FRONTEND_URL", "http://localhost:3000")), "/"),
		PasswordResetTTLMin: envInt("PASSWORD_RESET_TTL_MINUTES", 30),
		SMTPHost:            strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:            strings.TrimSpace(envOrDefault("SMTP_PORT", "587")),
		SMTPUsername:        strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
		SMTPPassword:        strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		SMTPFrom:            strings.TrimSpace(os.Getenv("SMTP_FROM")),
		StorageDriver:       storageDriver,
		DataFile:            dataFile,
		DatabaseURL:         databaseURL,
		DatabaseAutoMigrate: envBool("DATABASE_AUTO_MIGRATE", true),
		DatabaseMaxConns:    int32(envInt("DATABASE_MAX_CONNS", 10)),
		CloudProvider:       strings.ToLower(strings.TrimSpace(envOrDefault("CLOUD_PROVIDER", "local"))),
		AWSRegion:           strings.TrimSpace(os.Getenv("AWS_REGION")),
		AWSS3Bucket:         strings.TrimSpace(os.Getenv("AWS_S3_BUCKET")),
		AWSSecretsPrefix:    strings.TrimSpace(os.Getenv("AWS_SECRETS_PREFIX")),
		AWSUseIAMRole:       envBool("AWS_USE_IAM_ROLE", false),
		AWSCloudWatchNS:     strings.TrimSpace(os.Getenv("AWS_CLOUDWATCH_NAMESPACE")),
	}, nil
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback
	}

	return value == "1" || value == "true" || value == "yes" || value == "on"
}

func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open .env: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("set env %s: %w", key, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan .env: %w", err)
	}

	return nil
}

func findDotEnvPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	current := wd
	for {
		candidate := filepath.Join(current, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("stat %s: %w", candidate, err)
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", nil
}
