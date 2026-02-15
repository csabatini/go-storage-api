package config

import (
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("STORAGE_BACKEND", "local")
	t.Setenv("LOCAL_ROOT_PATH", "./data")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default Port 8080, got %s", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default LogLevel info, got %s", cfg.LogLevel)
	}
	if cfg.StorageBackend != "local" {
		t.Errorf("expected StorageBackend local, got %s", cfg.StorageBackend)
	}
	if cfg.MaxUploadSize != 104857600 {
		t.Errorf("expected default MaxUploadSize 104857600, got %d", cfg.MaxUploadSize)
	}
}

func TestLoadCustomValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("STORAGE_BACKEND", "local")
	t.Setenv("LOCAL_ROOT_PATH", "/tmp/files")
	t.Setenv("MAX_UPLOAD_SIZE", "52428800")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected Port 9090, got %s", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel debug, got %s", cfg.LogLevel)
	}
	if cfg.MaxUploadSize != 52428800 {
		t.Errorf("expected MaxUploadSize 52428800, got %d", cfg.MaxUploadSize)
	}
	if cfg.Local.RootPath != "/tmp/files" {
		t.Errorf("expected Local.RootPath /tmp/files, got %s", cfg.Local.RootPath)
	}
}

func TestLoadSMBBackendConfig(t *testing.T) {
	t.Setenv("STORAGE_BACKEND", "smb")
	t.Setenv("SMB_HOST", "fileserver.local")
	t.Setenv("SMB_PORT", "4455")
	t.Setenv("SMB_SHARE", "shared")
	t.Setenv("SMB_USER", "admin")
	t.Setenv("SMB_PASSWORD", "secret")

	cfg := Load()

	if cfg.SMB.Host != "fileserver.local" {
		t.Errorf("expected SMB.Host fileserver.local, got %s", cfg.SMB.Host)
	}
	if cfg.SMB.Port != "4455" {
		t.Errorf("expected SMB.Port 4455, got %s", cfg.SMB.Port)
	}
	if cfg.SMB.Share != "shared" {
		t.Errorf("expected SMB.Share shared, got %s", cfg.SMB.Share)
	}
}

func TestLoadFTPBackendConfig(t *testing.T) {
	t.Setenv("STORAGE_BACKEND", "ftp")
	t.Setenv("FTP_HOST", "ftp.example.com")
	t.Setenv("FTP_PORT", "2121")
	t.Setenv("FTP_USER", "ftpuser")
	t.Setenv("FTP_PASSWORD", "ftppass")

	cfg := Load()

	if cfg.FTP.Host != "ftp.example.com" {
		t.Errorf("expected FTP.Host ftp.example.com, got %s", cfg.FTP.Host)
	}
	if cfg.FTP.Port != "2121" {
		t.Errorf("expected FTP.Port 2121, got %s", cfg.FTP.Port)
	}
}

func TestLoadS3BackendConfig(t *testing.T) {
	t.Setenv("STORAGE_BACKEND", "s3")
	t.Setenv("S3_BUCKET", "my-bucket")
	t.Setenv("S3_REGION", "eu-west-1")
	t.Setenv("S3_PREFIX", "uploads/")

	cfg := Load()

	if cfg.S3.Bucket != "my-bucket" {
		t.Errorf("expected S3.Bucket my-bucket, got %s", cfg.S3.Bucket)
	}
	if cfg.S3.Region != "eu-west-1" {
		t.Errorf("expected S3.Region eu-west-1, got %s", cfg.S3.Region)
	}
	if cfg.S3.Prefix != "uploads/" {
		t.Errorf("expected S3.Prefix uploads/, got %s", cfg.S3.Prefix)
	}
}

func TestS3DefaultRegion(t *testing.T) {
	t.Setenv("STORAGE_BACKEND", "s3")
	t.Setenv("S3_BUCKET", "my-bucket")

	cfg := Load()

	if cfg.S3.Region != "us-east-1" {
		t.Errorf("expected default S3.Region us-east-1, got %s", cfg.S3.Region)
	}
}

func TestValidateBackendSMBMissingHost(t *testing.T) {
	cfg := &Config{
		StorageBackend: "smb",
		SMB:            SMBConfig{Share: "shared"},
	}
	err := cfg.validateBackend()
	if err == nil {
		t.Error("expected error for missing SMB_HOST")
	}
}

func TestValidateBackendSMBMissingShare(t *testing.T) {
	cfg := &Config{
		StorageBackend: "smb",
		SMB:            SMBConfig{Host: "server"},
	}
	err := cfg.validateBackend()
	if err == nil {
		t.Error("expected error for missing SMB_SHARE")
	}
}

func TestValidateBackendFTPMissingHost(t *testing.T) {
	cfg := &Config{
		StorageBackend: "ftp",
	}
	err := cfg.validateBackend()
	if err == nil {
		t.Error("expected error for missing FTP_HOST")
	}
}

func TestValidateBackendS3MissingBucket(t *testing.T) {
	cfg := &Config{
		StorageBackend: "s3",
	}
	err := cfg.validateBackend()
	if err == nil {
		t.Error("expected error for missing S3_BUCKET")
	}
}

func TestValidateBackendLocalMissingRoot(t *testing.T) {
	cfg := &Config{
		StorageBackend: "local",
		Local:          LocalConfig{RootPath: ""},
	}
	err := cfg.validateBackend()
	if err == nil {
		t.Error("expected error for missing LOCAL_ROOT_PATH")
	}
}
