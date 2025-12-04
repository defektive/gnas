package files

import (
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"sort"
)

// FileInfo represents file information for listing
type FileInfo struct {
	Name        string
	Size        int64
	IsDirectory bool
	Permissions string // Unix-like permissions string
}

// GetFileInfo retrieves file information for a given path.
// It handles errors and provides consistent output across platforms.
func GetFileInfo(path string) ([]FileInfo, error) {
	var files []FileInfo
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	for _, entry := range entries {
		entryInfo, err := entry.Info()
		var size int64
		var mode = fs.ModePerm
		if err != nil {
			size = 0
		} else {
			size = entryInfo.Size()
			mode = entryInfo.Mode()

		}

		fileInfo := FileInfo{
			Name:        entry.Name(),
			Size:        size,
			IsDirectory: entry.IsDir(),
		}

		// Handle permissions differently based on the operating system.
		if runtimeOS() == "windows" {
			fileInfo.Permissions = "" // Windows doesn't have traditional permissions
		} else {
			// Unix-like permissions
			fileInfo.Permissions = getPermissionsString(mode)
		}

		files = append(files, fileInfo)
	}
	return files, nil
}

// runtimeOS detects the operating system (windows, darwin, linux).
func runtimeOS() string {
	osType := runtime.GOOS
	switch osType {
	case "windows":
		return "windows"
	case "darwin":
		return "darwin" // macOS
	case "linux":
		return "linux"
	default:
		return "unknown" // Handle unexpected OS
	}
}

// PrintFileList prints the file list in a format similar to 'ls'.
func PrintFileList(files []FileInfo) {
	// Sort files alphabetically
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		fileType := ""
		if file.IsDirectory {
			fileType = "d" // directory
		} else {
			fileType = "-"
		}

		fmt.Printf("%s%s %10d %s\n", fileType, file.Permissions, file.Size, file.Name)
	}
}
