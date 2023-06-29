package work_packages

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
)

func upload(dto *dtos.WorkPackageDto, path string) error {
	if dto.Links.PrepareAttachment != nil {
		printer.ErrorText(fmt.Sprintf("Uploads to fog storages are currently not supported. :("))
	}

	printer.Info(fmt.Sprintf("Uploading %s to work package ...", printer.Yellow(filepath.Base(path))))
	link := dto.Links.AddAttachment
	reader, contentType, err := bodyReader(path)
	if err != nil {
		return err
	}

	body := &requests.RequestData{ContentType: contentType, Body: reader}
	_, err = requests.Do(link.Method, link.Href, nil, body)
	if err != nil {
		return err
	}

	printer.Done()
	return nil
}

func bodyReader(path string) (io.Reader, string, error) {
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
