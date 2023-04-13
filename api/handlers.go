package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/arekbor/file-manager-server/types"
	"github.com/arekbor/file-manager-server/utils"
	"github.com/gorilla/mux"
)

func (s *RestApiServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle Upload")

	limit, err := strconv.ParseInt(os.Getenv("HEADERS_LIMIT_MB"), 0, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = r.ParseMultipartForm(limit << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	files := r.MultipartForm.File["file"]

	folderName := r.Form.Get("folderName")
	fmt.Println(folderName)

	fullPath := os.Getenv("FILES_PATH_TO_DOWNLOAD")

	if folderName != "root" {
		fullPath = filepath.Join(os.Getenv("FILES_PATH_TO_DOWNLOAD"), folderName)
		_, err = os.Stat(fullPath)
		if err != nil {
			http.Error(w, "error while opening dir from body", http.StatusBadRequest)
			log.Println(err)
			return
		}
	}

	fmt.Println(fullPath)

	for _, file := range files {

		path := filepath.Join(fullPath, file.Filename)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			http.Error(w, "error while opening dir", http.StatusBadRequest)
			log.Println(err)
			return
		}
		defer f.Close()

		reader, err := file.Open()
		if err != nil {
			http.Error(w, "error while reading file from form: "+file.Filename, http.StatusBadRequest)
			log.Println(err)
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
			log.Println(err)
			return
		}
		writeJSON(w, file)
		return
	}

	files, err := fr.getFilesResponse()
	if err != nil {
		http.Error(w, "error while reading files or directories", http.StatusBadRequest)
		log.Println(err)
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
		log.Println(err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "error while readinga a file stats", http.StatusBadRequest)
		log.Println(err)
		return
	}
	fmt.Println(fi.ModTime())

	fr.getFileType(fi)

	switch t := fr.getFileType(fi); t {
	case types.AudioFileType:
		w.Header().Set("Content-Type", "audio/audio/mpeg")

	case types.TextFileType:
		w.Header().Set("Content-Type", fmt.Sprintf("text/%s", utils.GetFileExt(fi.Name())))

	case types.UnknowFileType:
		http.Error(w, "unknow type of file", http.StatusBadRequest)
		return
	}

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "error while copying a file", http.StatusInternalServerError)
		log.Println(err)
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
		log.Println(err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "error while readinga a file stats", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if fr.getFileType(fi) == types.UnknowFileType {
		http.Error(w, "unknow type of file", http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fi.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "error while copying a file", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (s *RestApiServer) handleGetAllFolderNames(w http.ResponseWriter, r *http.Request) {
	var (
		fr = newFileResponse(os.Getenv("FILES_PATH_TO_DOWNLOAD"), "", r)
	)
	result, err := fr.getAllFolderNames()
	if err != nil {
		http.Error(w, "error while reading folders", http.StatusBadRequest)
		log.Println(err)
		return
	}
	writeJSON(w, result)
}
