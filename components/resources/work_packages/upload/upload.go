package upload

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/opf/openproject-cli/components/requests"
)

func BodyReader(path string) (io.Reader, string, error) {
	var multipartFields = make(map[string]io.Reader)

	metaReader, err := metadata(filepath.Base(path))
	if err != nil {
		return nil, "", err
	}
	multipartFields["metadata"] = metaReader

	fileReader, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	multipartFields["file"] = fileReader

	return requests.MultipartBody(multipartFields)
}

func metadata(fileName string) (io.Reader, error) {
	type meta struct {
		FileName string `json:"fileName"`
	}

	marshal, err := json.Marshal(meta{FileName: fileName})
	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(marshal)), nil
}
