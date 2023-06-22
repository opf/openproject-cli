package work_packages

import (
	"fmt"
	"github.com/opf/openproject-cli/models"
	"path/filepath"
	"strconv"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

const path = "api/v3/work_packages"

func Lookup(id int64) *models.WorkPackage {
	status, response := requests.Get(filepath.Join(path, strconv.FormatInt(id, 10)), nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	workPackage := parser.Parse[WorkPackageDto](response)
	return workPackage.convert()
}

func Update(id int64) {
	fmt.Println("I UPDATE WORK PACKAGES")
}
