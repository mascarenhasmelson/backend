package error

import (
	"fmt"
	//  "log"
)

// Error test package
func IfError(err error) bool {
	if err == nil {
		return false
	}

	fmt.Sprintf("%v", err)
	return true
}
