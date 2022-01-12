package datastore

import "strings"

func IsErrNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no rows in result set")
}
