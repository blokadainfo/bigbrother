package autodelete

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type FileInfo struct {
	Name    string
	Size    int64
	IsDir   bool
	ModTime time.Time
}

func DeleteOldestFilesCron(dir string, limitInGB int, dry bool, dur time.Duration) {
	for {
		if err := deleteOldestFiles(dir, limitInGB, dry); err != nil {
			slog.Error("Failed to delete oldest files", "error", err)
		}
		time.Sleep(dur)
	}
}

// deleteOldestFiles deletes enough oldest files to make the directory size smaller than the limit
func deleteOldestFiles(dir string, limitInGB int, dry bool) error {
	limit := int64(limitInGB * 1024 * 1024 * 1024) // Convert GB to bytes
	var currentSize int64
	var files []FileInfo

	// If the limit is 0, no need to delete anything
	if limit == 0 {
		slog.Warn("Limit is set to 0, which means no limit")
		return nil
	}

	// Read all files in the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			currentSize += info.Size()
			files = append(files, FileInfo{
				Name:    info.Name(),
				Size:    info.Size(),
				IsDir:   info.IsDir(),
				ModTime: info.ModTime(),
			})
		}
		return nil
	})
	if err != nil {
		return err
	}

	// If current size is already under the limit, no need to delete anything
	if currentSize <= limit {
		return nil
	}

	// Sort files by modification time (oldest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.Before(files[j].ModTime)
	})

	// Delete files until the size is under the limit
	for _, file := range files {
		if currentSize <= limit {
			break
		}

		// Get full path of the file
		filePath := filepath.Join(dir, file.Name)

		// Delete the file
		if !dry {
			err := os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("failed to delete file %s: %w", filePath, err)
			}
		}

		// Update the current size after deletion
		currentSize -= file.Size

		if dry {
			slog.Warn("Dry run enabled, otherwise this file would get removed", "file_path", filePath, "limit", limit)
		} else {
			slog.Info("Deleted file to get the current size below the limit", "file_path", filePath, "current_size", currentSize, "limit", limit)
		}
	}

	return nil
}
