package logger

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

const maxLogSize = 10 * 1024 * 1024 // 10MB

// RotatingWriter wraps a file with automatic rotation
type RotatingWriter struct {
	filename string
	file     *os.File
	mu       sync.Mutex
}

// NewRotatingWriter creates a new rotating file writer
func NewRotatingWriter(filename string) (*RotatingWriter, error) {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Open or create file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &RotatingWriter{
		filename: filename,
		file:     file,
	}, nil
}

// Write implements io.Writer
func (rw *RotatingWriter) Write(p []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// Check file size
	info, err := rw.file.Stat()
	if err != nil {
		return 0, err
	}

	// If file exceeds max size, rotate it
	if info.Size() >= maxLogSize {
		if err := rw.rotate(); err != nil {
			return 0, err
		}
	}

	return rw.file.Write(p)
}

// Sync implements zapcore.WriteSyncer
func (rw *RotatingWriter) Sync() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.file.Sync()
}

// rotate closes current file and creates a new one
func (rw *RotatingWriter) rotate() error {
	// Close current file
	if err := rw.file.Close(); err != nil {
		return err
	}

	// Rename old file with timestamp
	oldFilename := rw.filename + "." + time.Now().Format("20060102-150405")
	if err := os.Rename(rw.filename, oldFilename); err != nil {
		// If rename fails, try to remove and create new
		os.Remove(rw.filename)
	}

	// Create new file
	file, err := os.OpenFile(rw.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	rw.file = file
	return nil
}

// Close closes the file
func (rw *RotatingWriter) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	if rw.file != nil {
		return rw.file.Close()
	}
	return nil
}

