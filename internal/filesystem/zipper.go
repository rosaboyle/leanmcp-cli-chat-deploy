package filesystem

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// ZipResult represents the result of a zip operation
type ZipResult struct {
	Data      []byte
	FileCount int
	TotalSize int64
}

// ProjectZipper handles zipping project files
type ProjectZipper struct {
	scanner *DirectoryScanner
}

// NewProjectZipper creates a new project zipper
func NewProjectZipper(projectPath string) *ProjectZipper {
	return &ProjectZipper{
		scanner: NewDirectoryScanner(projectPath),
	}
}

// CreateZip creates a zip file from the project directory
func (pz *ProjectZipper) CreateZip() (*ZipResult, error) {
	// Scan directory for files
	files, stats, err := pz.scanner.GetFileList()
	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}
	
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found to zip (directory might be empty or all files are ignored)")
	}
	
	// Create zip in memory
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	
	for _, file := range files {
		err := pz.addFileToZip(zipWriter, file)
		if err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("failed to add file %s to zip: %w", file.RelPath, err)
		}
	}
	
	err = zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to finalize zip: %w", err)
	}
	
	result := &ZipResult{
		Data:      buf.Bytes(),
		FileCount: stats.TotalFiles,
		TotalSize: stats.TotalSize,
	}
	
	return result, nil
}

// addFileToZip adds a single file to the zip archive
func (pz *ProjectZipper) addFileToZip(zipWriter *zip.Writer, file FileInfo) error {
	// Open source file
	sourceFile, err := os.Open(file.Path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer sourceFile.Close()
	
	// Get file info for permissions
	info, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	
	// Create zip file header
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("failed to create zip header: %w", err)
	}
	
	// Use forward slashes for zip paths (cross-platform compatibility)
	header.Name = strings.ReplaceAll(file.RelPath, "\\", "/")
	
	// Set compression method
	header.Method = zip.Deflate
	
	// Create file in zip
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create zip entry: %w", err)
	}
	
	// Copy file content
	_, err = io.Copy(writer, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}
	
	return nil
}

// ValidateZipSize checks if the zip size is within reasonable limits
func ValidateZipSize(data []byte) error {
	const maxSizeMB = 500 // 500MB limit
	const maxSizeBytes = maxSizeMB * 1024 * 1024
	
	if len(data) > maxSizeBytes {
		return fmt.Errorf("zip file too large: %d bytes (max %d MB)", len(data), maxSizeMB)
	}
	
	return nil
}

// GetHumanReadableSize converts bytes to human readable format
func GetHumanReadableSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// PreviewFiles returns a preview of files that would be included in the zip
func (pz *ProjectZipper) PreviewFiles(limit int) ([]FileInfo, FileStats, error) {
	files, stats, err := pz.scanner.GetFileList()
	if err != nil {
		return nil, stats, err
	}
	
	if len(files) > limit {
		return files[:limit], stats, nil
	}
	
	return files, stats, nil
}
