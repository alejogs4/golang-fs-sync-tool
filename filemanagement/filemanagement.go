package filemanagement

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FilesManagement struct to manipulate folder files information
type FilesManagement struct {
	folder string
}

// NewFilesManagement returns filesManagement internally new absolute folder
func NewFilesManagement(folderPath string) (FilesManagement, error) {
	absolutePath, err := filepath.Abs(folderPath)
	if err != nil {
		return FilesManagement{}, err
	}

	return FilesManagement{
		folder: absolutePath,
	}, nil
}

// FileInformation simple struct to keep general information of a file
type FileInformation struct {
	ID        string `json:"id"`
	Path      string `json:"path"`
	ShortPath string `json:"short_path"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
}

// FileBinaryMetadata hold OS file metadata
type FileBinaryMetadata struct {
	File       *os.File
	Size       int64
	FileHeader []byte
}

// GetFolderFiles get all files from teh given dir, including subdirectories
func (fm *FilesManagement) GetFolderFiles() ([]FileInformation, error) {
	var files []FileInformation
	err := filepath.Walk(fm.folder, fm.traverseDirectory(&files))

	if err != nil {
		return []FileInformation{}, err
	}

	return files, nil
}

// GetFile look for requested file and returns it alongside header information to get type
// besides its size and verifying if exists
// as likely this file will be used outside this fn, user of this fn must ensure that will be closed as soon as the operation finish
func (fm *FilesManagement) GetFile(filePath string) (FileBinaryMetadata, error) {
	searchedFilePath := filepath.Join(fm.folder, filePath)
	if _, err := os.Stat(searchedFilePath); os.IsNotExist(err) {
		return FileBinaryMetadata{}, err
	}

	searchedFile, err := os.Open(searchedFilePath)
	if err != nil {
		return FileBinaryMetadata{}, err
	}

	fileHeader := make([]byte, 512)
	_, err = searchedFile.Read(fileHeader)
	if err != nil {
		return FileBinaryMetadata{}, err
	}

	fileInformation, err := searchedFile.Stat()
	if err != nil {
		return FileBinaryMetadata{}, err
	}
	fileSize := fileInformation.Size()

	return FileBinaryMetadata{File: searchedFile, FileHeader: fileHeader, Size: fileSize}, nil
}

func (fm *FilesManagement) traverseDirectory(files *[]FileInformation) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			*files = append(*files, FileInformation{
				ID:        info.Name() + "-" + strconv.FormatInt(info.Size(), 10),
				Name:      info.Name(),
				Path:      path,
				Size:      info.Size(),
				ShortPath: fm.buildFileShortPath(path),
			})
		}

		// Avoid hidden files, those which begins with a dot(.) are included in the final response
		var notHiddenFiles []FileInformation = []FileInformation{}
		for _, gotFile := range *files {
			fileTree := strings.Split(gotFile.Path, "/")
			hasPoint := false

			for _, treePart := range fileTree {
				if strings.HasPrefix(treePart, ".") {
					hasPoint = true
					break
				}
			}

			if !hasPoint {
				notHiddenFiles = append(notHiddenFiles, gotFile)
			}
		}

		*files = notHiddenFiles

		return nil
	}
}

func (fm *FilesManagement) buildFileShortPath(filePath string) string {
	return strings.Replace(filePath, fm.folder, "", 1)[1:]
}
