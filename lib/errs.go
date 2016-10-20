package lib

import (
	"fmt"
)

type errUnexpectedResponseCode struct {
	URL      string
	Actual   int
	Expected int
}

// Error is the error interface implementation
func (e errUnexpectedResponseCode) Error() string {
	return fmt.Sprintf("%s - expected response code %d, actual %d", e.URL, e.Expected, e.Actual)
}
