package main

import (
	"fmt"
)

// CheckPlatform check
func CheckPlatform(platform string) *string {
	allowedPlatforms := []string{"psn", "origin", "xbl"}

	for _, allowedPlatform := range allowedPlatforms {
		if platform == allowedPlatform {
			return nil
		}
	}

	msg := new(string)
	*msg = fmt.Sprintf("platform %s not allowed, should be one of %v", platform, allowedPlatforms)
	return msg
}
