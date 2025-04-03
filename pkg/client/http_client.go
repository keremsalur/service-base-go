package client

import (
	"context"
	"net"
	"net/http"
	"service-base-go/pkg/logger"
	"time"

	"github.com/rs/zerolog"
)

func HttpClient() {
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com", nil)
	if err != nil {
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Request oluşturulamadı", map[string]interface{}{"error": err})
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Request atılamadı", map[string]interface{}{"error": err})
	}

	logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, "Response alındı", map[string]interface{}{"response": resp})
}
