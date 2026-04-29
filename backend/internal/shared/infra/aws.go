package infra

import (
	"strings"

	"api-on/internal/shared/config"
	sharederrors "api-on/internal/shared/errors"
)

type RuntimeProfile struct {
	CloudProvider    string
	StorageDriver    string
	AWSRegion        string
	AWSS3Bucket      string
	AWSSecretsPrefix string
	AWSUseIAMRole    bool
	AWSCloudWatchNS  string
}

// NewRuntimeProfile valida a configuração de infraestrutura usada pelo app.
func NewRuntimeProfile(cfg *config.Config) (*RuntimeProfile, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.CloudProvider))
	if provider == "" {
		provider = "local"
	}

	if provider != "local" && provider != "aws" {
		return nil, sharederrors.Invalid("INVALID_CLOUD_PROVIDER", "CLOUD_PROVIDER must be local or aws")
	}

	if provider == "aws" && strings.TrimSpace(cfg.AWSRegion) == "" {
		return nil, sharederrors.Invalid("INVALID_AWS_REGION", "AWS_REGION is required when CLOUD_PROVIDER=aws")
	}

	return &RuntimeProfile{
		CloudProvider:    provider,
		StorageDriver:    cfg.StorageDriver,
		AWSRegion:        cfg.AWSRegion,
		AWSS3Bucket:      cfg.AWSS3Bucket,
		AWSSecretsPrefix: cfg.AWSSecretsPrefix,
		AWSUseIAMRole:    cfg.AWSUseIAMRole,
		AWSCloudWatchNS:  cfg.AWSCloudWatchNS,
	}, nil
}
