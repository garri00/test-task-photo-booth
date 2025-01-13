package utils

import (
	"fmt"
	"io"
	"os"

	"test-task-photo-booth/pkg/logger"
	"test-task-photo-booth/src/entities/customErrors"
)

// IsFileExists check is file exists
// NOTE: do not return full file destination in error handling
func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

func GetFileBytes(filePath string) ([]byte, error) {
	//Check if file exists
	if !IsFileExists(filePath) {
		return nil, ErrNoCertFileFound
	}

	//Load file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open() failed: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Log.Error().Err(customErrors.ErrorOsCloseFailed).Err(err).Send()
		}
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll() failed: %w", err)
	}

	return fileBytes, nil
}
