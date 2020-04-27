package mysql_utils

import (
	"github.com/go-sql-driver/mysql"
	"github.com/nhsh1997/bookstore_utils-go/rest_errors"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"strings"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestError  {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows){
			return rest_errors.NewNotFoundError("no record matching given id")
		}
		return rest_errors.NewInternalServerError("error parsing database response", errors.New("database err"))
	}

	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestError("invalid data")
	}
	return rest_errors.NewInternalServerError("error processing request", errors.New("database err"))
}
