package utils

import "path/filepath"

// CheckExt checks for valid file extension
func CheckExt(filename string, ext string) bool {
	fileExt := filepath.Ext(filename)
	if fileExt != ext {
		return false
	}

	return true
}

// IsInSlice checks for existing string in slice, returns true if exists, false if not
func IsInSlice(str string, slice []string) bool {
	for _, v := range slice {
		if str == v {
			return true
		}
	}

	return false
}
