package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arekbor/file-manager-server/types"
	"github.com/arekbor/file-manager-server/utils"
	"github.com/gorilla/mux"
)

func (s *RestApiServer) handleUpload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "error while uploading files", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["file"]

	for _, file := range files {

		path := filepath.Join(os.Getenv("FILES_PATH_TO_DOWNLOAD"), file.Filename)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			http.Error(w, "error while opening dir", http.StatusBadRequest)
			return
		}
		defer f.Close()

		reader, err := file.Open()
		if err != nil {
			http.Error(w, "error while reading file from FORM: "+file.Filename, http.StatusBadRequest)
			return
		}

		defer reader.Close()

		_, err = io.Copy(f, reader)
	}
}

func (s *RestApiServer) handleManager(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)["path"]
		fr   = newFileResponse(os.Getenv("FILES_PATH_TO_DOWNLOAD"), vars, r)
	)
	_, err := ioutil.ReadDir(fr.fullPath)
	if err != nil {
		file, err := fr.getFileResponse()
		if err != nil {
			http.Error(w, "error while reading file or directory", http.StatusBadRequest)
			return
		}
		writeJSON(w, file)
		return
	}

	files, err := fr.getFilesResponse()
	if err != nil {
		http.Error(w, "error while reading files or directories", http.StatusBadRequest)
		return
	}

	writeJSON(w, files)
}

func (s *RestApiServer) handleStreamFile(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)["path"]
		fr   = newFileResponse(os.Getenv("FILES_PATH_TO_DOWNLOAD"), vars, r)
	)

	f, err := os.Open(fr.fullPath)
	if err != nil {
		http.Error(w, "error while reading a file", http.StatusBadRequest)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "error while readinga a file stats", http.StatusBadRequest)
		return
	}

	fr.getFileType(fi)

	switch t := fr.getFileType(fi); t {
	case types.AudioFileType:
		w.Header().Set("Content-Type", fmt.Sprintf("audio/%s", utils.GetFileExt(fi.Name())))

	case types.TextFileType:
		w.Header().Set("Content-Type", fmt.Sprintf("text/%s", utils.GetFileExt(fi.Name())))

	case types.UnknowFileType:
		http.Error(w, "unknow type of file", http.StatusBadRequest)
		return
	}

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "error while copying a file", http.StatusInternalServerError)
		return
	}
}

func (s *RestApiServer) handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)["path"]
		fr   = newFileResponse(os.Getenv("FILES_PATH_TO_DOWNLOAD"), vars, r)
	)

	f, err := os.Open(fr.fullPath)
	if err != nil {
		http.Error(w, "error while reading a file", http.StatusBadRequest)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "error while readinga a file stats", http.StatusBadRequest)
		return
	}

	if fr.getFileType(fi) == types.UnknowFileType {
		http.Error(w, "unknow type of file", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fi.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "error while copying a file", http.StatusInternalServerError)
		return
	}
}
