package requests

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func MultipartBody(values map[string]io.Reader) (io.Reader, string, error) {
	var b bytes.Buffer
	var err error
	writer := multipart.NewWriter(&b)

	for field, reader := range values {
		var fw io.Writer

		switch reader.(type) {
		case *os.File:
			file := reader.(*os.File)
			if fw, err = writer.CreateFormFile(field, filepath.Base(file.Name())); err != nil {
				return nil, "", err
			}
		default:
			if fw, err = writer.CreateFormField(field); err != nil {
				return nil, "", err
			}
		}

		if _, err = io.Copy(fw, reader); err != nil {
			return nil, "", err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return bytes.NewReader(b.Bytes()), writer.FormDataContentType(), nil
}
