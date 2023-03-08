package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arekbor/file-manager-server/types"
	"github.com/arekbor/file-manager-server/utils"
)

// Writes any value to JSON
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

type fileResponse struct {
	//fullPath is joined path from http request and base path
	fullPath string

	//rawPath is only requested cleaned path from http variable
	rawPath string

	//basePath is a path to base directory
	//
	//for example you migh pass some variable from environment
	basePath string

	request *http.Request
}

// Creates reference to fileResponse
func newFileResponse(basePath string, rawPath string, request *http.Request) *fileResponse {

	var (
		cleanPath = filepath.Clean(rawPath)
		fullPath  = filepath.Join(basePath, cleanPath)
	)

	return &fileResponse{
		fullPath: fullPath,
		rawPath:  cleanPath,
		request:  request,
		basePath: basePath,
	}
}

// Gets all files from base directory
func (fr *fileResponse) getFilesResponse() ([]*types.File, error) {
	filesInfo, err := ioutil.ReadDir(fr.fullPath)
	if err != nil {
		return nil, err
	}

	files := []*types.File{}

	for _, file := range filesInfo {
		f, err := fr.getFile(file)
		if err == nil {
			files = append(files, f)
		}
	}

	return files, nil
}

// Gets signle file from directory
func (fr *fileResponse) getFileResponse() (*types.File, error) {
	file, err := os.Open(fr.fullPath)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	f, err := fr.getFile(fileInfo)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Transfors a file info into the file struct
func (fr *fileResponse) getFile(fs fs.FileInfo) (*types.File, error) {
	if utils.IsExtHasAnyValue(fs.Name()) || fs.IsDir() {
		return &types.File{
			Id:           utils.RandomNumber(),
			FileType:     fr.getFileType(fs),
			FileName:     fs.Name(),
			Size:         utils.ByteToHuman(fs),
			PathFile:     fr.getFilePath(fs),
			StreamPath:   fr.getStreamPath(fs),
			PathDownload: fr.getDownloadPath(fs),
		}, nil
	}
	return nil, errors.New("error while transforming a file")
}

// Gets url path to file
func (fr *fileResponse) getFilePath(fs fs.FileInfo) string {
	if !utils.IsExtHasAnyValue(fr.rawPath) {
		path := filepath.Join(fr.rawPath, fs.Name())

		return fmt.Sprintf("/api/manager/%s", path)
	}
	return ""
}

// Gets url path to file stream
func (fr *fileResponse) getStreamPath(fs fs.FileInfo) string {
	if utils.IsExtHasAnyValue(fr.rawPath) {
		host := utils.GetFullHostRequest(fr.request)
		fmt.Printf("%s/api/stream/%s\n", host, fr.rawPath)
		return fmt.Sprintf("%s/api/stream/%s", host, fr.rawPath)
	}
	return ""
}

// Gets url path to download a file
func (fr *fileResponse) getDownloadPath(fs fs.FileInfo) string {
	if utils.IsExtHasAnyValue(fr.rawPath) {
		host := utils.GetFullHostRequest(fr.request)

		return fmt.Sprintf("%s/api/download/%s", host, fr.rawPath)
	}
	return ""
}

// Transforms FileInfo into specific enum of FileType
func (fr *fileResponse) getFileType(fi os.FileInfo) types.FileType {
	var (
		fullPath = utils.JoinPaths(fi.Name())
		ext      = utils.GetFileExt(fullPath)
	)

	if fi.IsDir() {
		return types.DirectoryFileType
	}
	if utils.IsSliceHas(utils.AudioExtTypes, ext) {
		return types.AudioFileType
	}
	if utils.IsSliceHas(utils.TextExtTypes, ext) {
		return types.TextFileType
	}
	if utils.IsSliceHas(utils.ImageExtTypes, ext) {
		return types.ImageFileType
	}

	return types.UnknowFileType
}
