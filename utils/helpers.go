package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	AudioExtTypes   = []string{"mp3", "wav", "m4a", "flac", "wma"}
	TextExtTypes    = []string{"php", "js", "go", "rust", "c", "cpp", "cp", "json", "jar", "html", "css", "txt", "log"}
	ImageExtTypes   = []string{"png", "jpeg", "jpg", "gif", "svg", "webp", "avif"}
	VideoExtTypes   = []string{"mp4", "wmv", "avi", "mkv", "webm", "flv", "mov"}
	ArchiveExtTypes = []string{"rar", "zip", "iso", "tar", "mar", "a", "ar", "lbr", "br", "gz", "lz", "7z", "s7s", "ace"}
)

// Converts bytes to human readable stats
func ByteToHuman(f os.FileInfo) string {
	b := f.Size()
	if !f.IsDir() {
		const unit = 1000
		if b < unit {
			return fmt.Sprintf("%d B", b)
		}
		div, exp := int64(unit), 0
		for n := b / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}
		return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
	}
	return ""
}

// Joins filename with absolute path and return full string path
//
// fileName string - relative path of file
func JoinPaths(fileName string) string {
	path := filepath.Join(os.Getenv("FILES_PATH_TO_DOWNLOAD"), fileName)
	return path
}

// Gets file extension
//
// filePath string - full path of file
func GetFileExt(filePath string) string {
	return strings.TrimPrefix(filepath.Ext(filePath), ".")
}

// Checks if any element is included in slice array
//
// s []T - slice of array
//
// v T element to compare
func IsSliceHas[T comparable](s []T, v T) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}

// Gets from request full host path
func GetFullHostRequest(r *http.Request) string {
	proto := GetProtoFromRequest(r)

	return fmt.Sprintf("%s://%s", proto, r.Host)
}

// Chekcs if file extension has a value
func IsExtHasAnyValue(filePath string) bool {
	ext := filepath.Ext(filePath)

	if len(ext) > 0 {
		ext = strings.Split(ext, ".")[1]
		if len(ext) > 0 {
			return true
		}
	}

	return false
}

func RandomNumber() int {
	rand.Seed(time.Now().UnixNano())

	min := 1
	max := 1000000

	return rand.Intn(max-min+1) + min
}

func GetProtoFromRequest(r *http.Request) string {
	return os.Getenv("DEFAULT_PROTO")
}
