package photo

import (
	"path/filepath"
	"strings"
)

func isJPEG(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		return true
	default:
		return false
	}
}
