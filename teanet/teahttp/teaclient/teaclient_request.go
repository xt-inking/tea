package teaclient

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tea-frame-go/tea/teaencoding/teajson"
)

func NewRequest(ctx context.Context, method string, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, nil)
}

func NewRequestMultipartForm(ctx context.Context, method string, url string, value map[string][]string, file map[string][]string) (*http.Request, error) {
	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)
	for fieldName, values := range value {
		for _, value := range values {
			err := writer.WriteField(fieldName, value)
			if err != nil {
				return nil, err
			}
		}
	}
	for fieldName, filePaths := range file {
		for _, filePath := range filePaths {
			err := writeFile(writer, fieldName, filePath)
			if err != nil {
				return nil, err
			}
		}
	}
	writer.Close()
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func writeFile(writer *multipart.Writer, fieldName string, filePath string) error {
	p, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(p, file)
	return err
}

func NewRequestForm(ctx context.Context, method string, url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func NewRequestJson(ctx context.Context, method string, url string, data any) (*http.Request, error) {
	body := bytes.NewBuffer(nil)
	err := teajson.Encode(body, data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
