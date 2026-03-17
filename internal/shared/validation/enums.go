package validation

import (
	"fmt"
	"strings"
)

func Enum(field, value string, allowed []string) error {
	for _, item := range allowed {
		if value == item {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of [%s]", field, strings.Join(allowed, ", "))
}
