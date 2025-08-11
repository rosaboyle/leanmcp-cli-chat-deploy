package filesystem

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileInfo represents information about a file
type FileInfo struct {
	Path     string
	RelPath  string
	Size     int64
	IsDir    bool
	Mode     os.FileMode
}

// FileStats represents statistics about scanned files
type FileStats struct {
	TotalFiles int
	TotalSize  int64
	TotalDirs  int
}

// DirectoryScanner scans directories and respects ignore patterns
type DirectoryScanner struct {
	rootPath        string
	gitignoreRules  []string
	excludePatterns []string
}

// NewDirectoryScanner creates a new directory scanner
func NewDirectoryScanner(rootPath string) *DirectoryScanner {
	scanner := &DirectoryScanner{
		rootPath: rootPath,
		excludePatterns: []string{
			".git",
			".DS_Store",
			"*.log",
			"node_modules",
			".env",
			".env.local",
			".env.production",
			".env.development",
			"*.tmp",
			"*.temp",
			"dist",
			"build",
			"coverage",
			".nyc_output",
			".idea",
			".vscode",
			"*.swp",
			"*.swo",
			"*~",
		},
	}
	
	// Load .gitignore rules
	scanner.loadGitignoreRules()
	
	return scanner
}

// loadGitignoreRules loads patterns from .gitignore file
func (ds *DirectoryScanner) loadGitignoreRules() {
	gitignorePath := filepath.Join(ds.rootPath, ".gitignore")
	
	file, err := os.Open(gitignorePath)
	if err != nil {
		// No .gitignore file, that's okay
		return
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		ds.gitignoreRules = append(ds.gitignoreRules, line)
	}
}

// ScanDirectory recursively scans a directory and returns file information
func (ds *DirectoryScanner) ScanDirectory() ([]FileInfo, FileStats, error) {
	var files []FileInfo
	stats := FileStats{}
	
	err := filepath.Walk(ds.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Get relative path from root
		relPath, err := filepath.Rel(ds.rootPath, path)
		if err != nil {
			return err
		}
		
		// Skip root directory itself
		if relPath == "." {
			return nil
		}
		
		// Check if this path should be ignored
		if ds.shouldIgnore(relPath, info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		
		fileInfo := FileInfo{
			Path:    path,
			RelPath: relPath,
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			Mode:    info.Mode(),
		}
		
		files = append(files, fileInfo)
		
		if info.IsDir() {
			stats.TotalDirs++
		} else {
			stats.TotalFiles++
			stats.TotalSize += info.Size()
		}
		
		return nil
	})
	
	if err != nil {
		return nil, stats, fmt.Errorf("failed to scan directory: %w", err)
	}
	
	return files, stats, nil
}

// shouldIgnore checks if a path should be ignored based on patterns
func (ds *DirectoryScanner) shouldIgnore(relPath string, isDir bool) bool {
	// Check exclude patterns
	for _, pattern := range ds.excludePatterns {
		if ds.matchPattern(pattern, relPath, isDir) {
			return true
		}
	}
	
	// Check .gitignore rules
	for _, rule := range ds.gitignoreRules {
		if ds.matchGitignoreRule(rule, relPath, isDir) {
			return true
		}
	}
	
	return false
}

// matchPattern checks if a path matches a simple pattern
func (ds *DirectoryScanner) matchPattern(pattern, path string, isDir bool) bool {
	// Handle wildcards
	if strings.Contains(pattern, "*") {
		matched, _ := filepath.Match(pattern, filepath.Base(path))
		return matched
	}
	
	// Exact match for directories
	if isDir && pattern == filepath.Base(path) {
		return true
	}
	
	// Check if path starts with pattern (for directory patterns)
	if strings.HasPrefix(path, pattern+string(filepath.Separator)) {
		return true
	}
	
	// Exact match
	return pattern == path || pattern == filepath.Base(path)
}

// matchGitignoreRule checks if a path matches a .gitignore rule
func (ds *DirectoryScanner) matchGitignoreRule(rule, path string, isDir bool) bool {
	// Handle negation rules (starting with !)
	if strings.HasPrefix(rule, "!") {
		return false // Negation rules are complex, skip for now
	}
	
	// Handle directory-only rules (ending with /)
	if strings.HasSuffix(rule, "/") {
		if !isDir {
			return false
		}
		rule = strings.TrimSuffix(rule, "/")
	}
	
	// Handle absolute paths (starting with /)
	if strings.HasPrefix(rule, "/") {
		rule = strings.TrimPrefix(rule, "/")
		return ds.matchPattern(rule, path, isDir)
	}
	
	// Handle wildcards and regular patterns
	return ds.matchPattern(rule, path, isDir) || ds.matchPattern(rule, filepath.Base(path), isDir)
}

// GetFileList returns only non-directory files
func (ds *DirectoryScanner) GetFileList() ([]FileInfo, FileStats, error) {
	allFiles, stats, err := ds.ScanDirectory()
	if err != nil {
		return nil, stats, err
	}
	
	var files []FileInfo
	fileStats := FileStats{TotalDirs: stats.TotalDirs}
	
	for _, file := range allFiles {
		if !file.IsDir {
			files = append(files, file)
			fileStats.TotalFiles++
			fileStats.TotalSize += file.Size
		}
	}
	
	return files, fileStats, nil
}

// ValidateDirectory checks if a directory is suitable for project creation
func ValidateDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return fmt.Errorf("cannot access directory: %w", err)
	}
	
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}
	
	// Check if directory is readable
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot read directory: %w", err)
	}
	file.Close()
	
	return nil
}
