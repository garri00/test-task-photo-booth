package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func GetB64MimeType(data []byte) string {
	mimeType := http.DetectContentType(data)

	return mimeType
}

// GetB64WithMimeType used for embedding b64 data in json response for (PNG,JPEG)
func GetB64WithMimeType(data []byte) string {
	mimeType := http.DetectContentType(data)

	// Encode file bytes to Base64
	encoded := base64.StdEncoding.EncodeToString(data)

	b64DataWithMimeType := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	return b64DataWithMimeType
}
