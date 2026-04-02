package restc

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileUpload struct {
	FieldName string
	FileName  string
	Reader    io.Reader
}

func (r *Request) SetFormData(data map[string]string) *Request {
	r.ensureFormData()
	maps.Copy(r.formData, data)
	return r
}

func (r *Request) SetFileReader(fieldName, fileName string, reader io.Reader) *Request {
	r.files = append(r.files, &FileUpload{
		FieldName: fieldName,
		FileName:  fileName,
		Reader:    reader,
	})
	return r
}

func (r *Request) SetFile(fieldName, filePath string) *Request {
	file, err := os.Open(filePath)
	if err != nil {
		r.multipartErr = err
		return r
	}

	r.files = append(r.files, &FileUpload{
		FieldName: fieldName,
		FileName:  filepath.Base(filePath),
		Reader:    file,
	})
	return r
}

func (r *Request) ensureFormData() {
	if r.formData == nil {
		r.formData = make(map[string]string)
	}
}

func (r *Request) buildMultipartBody() (io.Reader, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for k, v := range r.formData {
		if err := writer.WriteField(k, v); err != nil {
			return nil, "", fmt.Errorf("%w: %w",
				ErrMultipart,
				fmt.Errorf("failed to write form field '%q': %w", k, err),
			)
		}
	}

	for _, file := range r.files {
		part, err := writer.CreateFormFile(file.FieldName, file.FileName)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %w",
				ErrMultipart,
				fmt.Errorf("failed to create form file '%q': %w", file.FieldName, err),
			)
		}
		if _, err := io.Copy(part, file.Reader); err != nil {
			return nil, "", fmt.Errorf("%w: %w",
				ErrMultipart,
				fmt.Errorf("failed to write form file '%q': %w", file.FieldName, err),
			)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrMultipart, err)
	}

	return &buf, writer.FormDataContentType(), nil
}
