package client

import (
	"context"
	"ecm-sdk-go/utils"
)

// customCredential
type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {

	maxRetryTimes := 10
	serviceName, backendName, token, err := utils.ParseBackendInfo(maxRetryTimes)
	if err != nil {
		return map[string]string{
			"serviceName": serviceName,
			"backendName": backendName,
			"token":       token,
		}, err
	}

	return map[string]string{
		"serviceName": serviceName,
		"backendName": backendName,
		"token":       token,
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return false
}
