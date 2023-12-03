package databaseHelpers

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func IsDuplicateKeyForField(err *mysql.MySQLError, fieldName string) bool {
	// Check if the error message contains the field name
	return strings.Contains(err.Message, fmt.Sprintf("for key '%s'", fieldName))
}
