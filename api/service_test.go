package api

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/arekbor/file-manager-server/types"
	"github.com/stretchr/testify/assert"
)

func TestGetFilePath(t *testing.T) {
	fr := newFileResponse("/home/main", "/api/test", nil)

	fs := &dummyFileInfo{
		name:  "testfile.txt",
		isDir: false,
	}

	path := filepath.Join(fr.rawPath, fs.Name())
	result := fr.getFilePath(fs)

	test := fmt.Sprintf("/api/manager/%s", path)
	assert.Equal(t, result, test)

	fr = newFileResponse("/home/main", "/api/test/testfile.txt", nil)

	result = fr.getFilePath(fs)

	assert.Equal(t, result, "")

}

func TestGetFileType(t *testing.T) {
	fr := &fileResponse{}

	dirInfo := &dummyFileInfo{
		name:  "testDir",
		isDir: true,
	}
	fileType := fr.getFileType(dirInfo)
	assert.Equal(t, fileType, types.DirectoryFileType)

	audioInfo := &dummyFileInfo{
		name: "testaudio.mp3",
	}
	fileType = fr.getFileType(audioInfo)
	assert.Equal(t, fileType, types.AudioFileType)

	textInfo := &dummyFileInfo{
		name: "index.php",
	}
	fileType = fr.getFileType(textInfo)
	assert.Equal(t, fileType, types.TextFileType)

	unknowInfo := &dummyFileInfo{
		name:  "test",
		isDir: false,
	}
	fileType = fr.getFileType(unknowInfo)
	assert.Equal(t, fileType, types.UnknowFileType)
}

type dummyFileInfo struct {
	name  string
	isDir bool
}

func (df *dummyFileInfo) Name() string {
	return df.name
}

func (df *dummyFileInfo) IsDir() bool {
	return df.isDir
}

func (df *dummyFileInfo) Size() int64 {
	return 0
}

func (fdf *dummyFileInfo) Mode() os.FileMode {
	return 0
}

func (df *dummyFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (df *dummyFileInfo) Sys() interface{} {
	return nil
}
