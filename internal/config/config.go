package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	LogLevel       string
	StorageBackend string
	MaxUploadSize  int64
	Local          LocalConfig
	SMB            SMBConfig
	FTP            FTPConfig
	S3             S3Config
}

type LocalConfig struct {
	RootPath string
}

type SMBConfig struct {
	Host     string
	Port     string
	Share    string
	User     string
	Password string
}

type FTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type S3Config struct {
	Bucket string
	Region string
	Prefix string
}

func Load() *Config {
	backend := envOrDefault("STORAGE_BACKEND", "local")

	validBackends := map[string]bool{
		"local": true,
		"smb":   true,
		"ftp":   true,
		"s3":    true,
	}
	if !validBackends[backend] {
		log.Fatalf("invalid STORAGE_BACKEND: %q (must be one of: local, smb, ftp, s3)", backend)
	}

	maxUpload, err := strconv.ParseInt(envOrDefault("MAX_UPLOAD_SIZE", "104857600"), 10, 64)
	if err != nil {
		log.Fatalf("invalid MAX_UPLOAD_SIZE: %v", err)
	}

	cfg := &Config{
		Port:           envOrDefault("PORT", "8080"),
		LogLevel:       envOrDefault("LOG_LEVEL", "info"),
		StorageBackend: backend,
		MaxUploadSize:  maxUpload,
		Local: LocalConfig{
			RootPath: envOrDefault("LOCAL_ROOT_PATH", "./data"),
		},
		SMB: SMBConfig{
			Host:     os.Getenv("SMB_HOST"),
			Port:     envOrDefault("SMB_PORT", "445"),
			Share:    os.Getenv("SMB_SHARE"),
			User:     os.Getenv("SMB_USER"),
			Password: os.Getenv("SMB_PASSWORD"),
		},
		FTP: FTPConfig{
			Host:     os.Getenv("FTP_HOST"),
			Port:     envOrDefault("FTP_PORT", "21"),
			User:     os.Getenv("FTP_USER"),
			Password: os.Getenv("FTP_PASSWORD"),
		},
		S3: S3Config{
			Bucket: os.Getenv("S3_BUCKET"),
			Region: envOrDefault("S3_REGION", "us-east-1"),
			Prefix: os.Getenv("S3_PREFIX"),
		},
	}

	if err := cfg.validateBackend(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	return cfg
}

func (c *Config) validateBackend() error {
	switch c.StorageBackend {
	case "local":
		if c.Local.RootPath == "" {
			return fmt.Errorf("LOCAL_ROOT_PATH is required for local backend")
		}
	case "smb":
		if c.SMB.Host == "" {
			return fmt.Errorf("SMB_HOST is required for smb backend")
		}
		if c.SMB.Share == "" {
			return fmt.Errorf("SMB_SHARE is required for smb backend")
		}
	case "ftp":
		if c.FTP.Host == "" {
			return fmt.Errorf("FTP_HOST is required for ftp backend")
		}
	case "s3":
		if c.S3.Bucket == "" {
			return fmt.Errorf("S3_BUCKET is required for s3 backend")
		}
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
