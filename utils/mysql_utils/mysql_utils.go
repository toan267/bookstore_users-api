package mysql_utils

import (
	"github.com/go-sql-driver/mysql"
	"strings"
	"toan267/bookstore_users-api/utils/errors"
)

const (
	errorNoRows = "sql: no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		//mysql error
		switch sqlErr.Number {
		case 1062:
			return errors.NewBadRequestError("invalid data")
		}
		return errors.NewInternalServerError("error processing request")
	}
	if strings.Contains(err.Error(), errorNoRows) {
		return errors.NewNotFoundError("no record matching given id")
	}
	return errors.NewInternalServerError("error parsing database response")
}