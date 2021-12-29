package cloud

import (
	"bytes"
	a1 "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/anhdht/go-common/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	pb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"os"
)

type SecretManager interface {
	Load() ([]byte, error)
}

type secretManager struct {
}

func NewSecretManager() SecretManager {
	return &secretManager{}
}

func (ptr *secretManager) Load() ([]byte, error) {
	ctx := context.Background()
	logger := log.Logger(ctx)

	projectId := os.Getenv("PROJECT_ID")
	serviceId := os.Getenv("SERVICE_ID")
	if projectId == "" {
		return nil, nil
	}

	viper.SetConfigType("yaml")
	data, err := os.ReadFile(fmt.Sprintf("/%s/env/%s.yaml", serviceId, projectId))
	if err != nil {
		logger.Fatal("encounter error while reading config file", zap.Error(err))
	}
	err = viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		logger.Fatal("encounter error while parsing config file", zap.Error(err))
	}

	sid := viper.Get("SID")
	sv := viper.Get("SV")

	client, err := a1.NewClient(ctx)
	if err != nil {
		logger.Fatal("failed to setup client", zap.Error(err))
	}
	defer client.Close()

	accessRequest := &pb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%d/secrets/config/versions/%d", sid, sv),
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		logger.Fatal("failed to access secret version", zap.Error(err))
	}

	return result.Payload.Data, nil
}
