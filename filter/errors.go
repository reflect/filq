package filter

import (
	"fmt"
)

type CannotSubscriptError struct {
	Dimensions int
}

func (e *CannotSubscriptError) Error() string {
	return fmt.Sprintf("filter returned %d outputs (wanted exactly 1)", e.Dimensions)
}
