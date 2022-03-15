package utils

import (
	"fmt"
	"io"
)

// CloseOrLog closes closable body and logs if close fails
func CloseOrLog(closable io.Closer) {
	if closable != nil {
		err := closable.Close()
		if err != nil {
			fmt.Println("closing object failed") // TODO: make this log better
		}
	}
}
