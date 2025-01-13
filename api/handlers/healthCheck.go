package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"

	"test-task-photo-booth/pkg/logger"
	"test-task-photo-booth/src/config"
	"test-task-photo-booth/src/entities"
)

type healthCheckResponse struct {
	Service        string `json:"service"`
	ServiceVersion string `json:"serviceVersion"`
	Status         string `json:"status"`
	Mode           string `json:"mode"`
}

const statusOK = "OK"

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	configs, err := config.GetConfig()
	if err != nil {
		RespondErr(w, &logger.Log, fmt.Errorf("GetConfig() failed: %w", err), http.StatusInternalServerError)

		return
	}

	serviceName := viper.GetString(entities.ServiceName)
	serviceVersion := viper.GetString(entities.ServicesVersion)

	response := healthCheckResponse{
		Service:        serviceName,
		ServiceVersion: serviceVersion,
		Status:         statusOK,
		Mode:           configs.Mode,
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		RespondErr(w, &logger.Log, fmt.Errorf("json.Marshal() failed: %w", err), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		logger.Log.Error().Err(fmt.Errorf("error writing response: %w", err)).Send()

		return
	}
}
