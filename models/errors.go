package models

import (
	"github.com/RealImage/QLedger/errors"
)

func SearchNamespaceInvalidError(namespace string) errors.ApplicationError {
	return &errors.BaseApplicationError{
		Code:    "search.namespace.invalid",
		Message: "Invalid search namespace: " + namespace,
	}
}

func SearchQueryInvalidError(err error) errors.ApplicationError {
	return &errors.BaseApplicationError{
		Code:    "search.query.invalid",
		Message: "Invalid search query: " + err.Error(),
	}
}

func DBError(err error) errors.ApplicationError {
	return &errors.BaseApplicationError{
		Code:    "db.error",
		Message: "DB Error: " + err.Error(),
	}
}
