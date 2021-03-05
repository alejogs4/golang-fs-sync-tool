package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/alejogs4/synctool/filemanagement"
)

// HandleGetAllFolderFiles controller to get all files from a given folder holded by filesManagement
func HandleGetAllFolderFiles(filesManagement *filemanagement.FilesManagement) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")

		files, err := filesManagement.GetFolderFiles()
		if err != nil {
			response, _ := json.Marshal(map[string]interface{}{"message": err.Error()})
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(response))
			return
		}

		response, _ := json.Marshal(map[string]interface{}{"data": files})
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(response))
	}
}

// HandleDownloadFile download specific file from the folder, file to download comes from query parameter filepath
func HandleDownloadFile(fileManagement *filemanagement.FilesManagement) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")

		folderFilepath := r.URL.Query().Get("filepath")
		fileMetadata, err := fileManagement.GetFile(folderFilepath)

		if err != nil {
			response, _ := json.Marshal(map[string]interface{}{"message": err.Error()})
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(response))
			return
		}

		defer fileMetadata.File.Close()
		fileType := http.DetectContentType(fileMetadata.FileHeader)

		// Set response headers
		rw.Header().Set("Expires", "0")
		rw.Header().Set("Content-Transfer-Encoding", "binary")
		rw.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")
		rw.Header().Set("Content-Disposition", "attachment; filename="+folderFilepath)
		rw.Header().Set("Content-Type", fileType)
		rw.Header().Set("Content-Length", strconv.FormatInt(fileMetadata.Size, 10))

		fileMetadata.File.Seek(0, 0)
		io.Copy(rw, fileMetadata.File)
		return
	}
}
