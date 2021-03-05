package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alejogs4/synctool/controllers"
	"github.com/alejogs4/synctool/filemanagement"
)

func main() {
	folder := flag.String("folder", "/", "Folder to look up for files and directories")
	port := flag.String("port", ":8080", "Server port to serve the resources from the given folder")
	flag.Parse()

	folderManagement, err := filemanagement.NewFilesManagement(*folder)
	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServeMux()
	server.HandleFunc("/api/v1/files", controllers.HandleGetAllFolderFiles(&folderManagement))
	server.HandleFunc("/api/v1/file", controllers.HandleDownloadFile(&folderManagement))

	log.Println("Listening requests at port " + *port)
	log.Fatal(http.ListenAndServe(*port, server))
}
