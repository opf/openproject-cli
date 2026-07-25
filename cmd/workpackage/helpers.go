package workpackage

import opErrors "github.com/opf/openproject-cli/components/errors"

func isNotFound(err error) bool {
	if respErr, ok := err.(*opErrors.ResponseError); ok {
		return respErr.Status() == 404
	}
	return false
}
