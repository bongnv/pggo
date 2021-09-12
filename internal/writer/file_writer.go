package writer

import (
	"os"
	"path"
)

// FileWriter writes content to file.
type FileWriter struct {
	Dir string
}

// Write writes content into a file.
func (w FileWriter) Write(fileName string, content []byte) error {
	return os.WriteFile(path.Join(w.Dir, fileName), content, 0644)
}
