package utils

import (
	"fmt"
	"os"
)

func GenHref(id string) string {
	if os.Getenv("HOST") == "" {
		return fmt.Sprintf("%s/todo/%s", "{{HOST}}", id)
	}
	return fmt.Sprintf("%s/todo/%s", os.Getenv("HOST"), id)
}
