package repo

import (
	"fmt"
	"reflect"

	"jwt-enabled-api/models"
)

// var tokenGetQuery = fmt.Sprintf(
// 	`SELECT %[1]s FROM %[2]s WHERE %[3]s = crypt($1, %[3]s) AND %[4]s > CURRENT_TIMESTAMP;`,
// 	/*1*/ columnsList(models.TokenColumns),
// 	/*2*/ models.TableNames.Tokens,
// 	/*3*/ models.TokenColumns.TokenHash,
// 	/*4*/ models.TokenColumns.ExpiresAt,
// )

var tokenGetQuery = fmt.Sprintf(
	`WITH user_tokens AS (
		SELECT %[1]s FROM %[2]s WHERE %[3]s = $1 AND %[4]s > CURRENT_TIMESTAMP
	)
	SELECT %[1]s FROM user_tokens WHERE %[5]s = crypt($2, %[5]s);`,
	/*1*/ columnsList(models.TokenColumns),
	/*2*/ models.TableNames.Tokens,
	/*3*/ models.TokenColumns.UserID,
	/*4*/ models.TokenColumns.ExpiresAt,
	/*5*/ models.TokenColumns.TokenHash,
)

func columnsList(tableColumnsStruct any) string {
	v := reflect.ValueOf(tableColumnsStruct)
	columns := ``
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldValue, isString := field.Interface().(string)
		if !isString {
			panic("columnsList: all fields must be of type string")
		}
		if i != 0 {
			columns += `, `
		}
		columns += fmt.Sprintf(`%[1]s AS "%[1]s"`, fieldValue)
	}
	return columns
}
