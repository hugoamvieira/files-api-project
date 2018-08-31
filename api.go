package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

var (
	errInternalInvalidJSONBody       = errors.New("Invalid JSON Body")
	errInternalInvalidPath           = errors.New("Invalid Path")
	errInternalNoData                = errors.New("No data to write to file")
	errInternalInvalidFilename       = errors.New("Invalid filename")
	errInternalCreatingOrWritingFile = errors.New("Error creating/writing file")
	errInternalFileNotFound          = errors.New("File not found")
	errInternalPathNotFound          = errors.New("Path not found")
	errInternalGeneric               = errors.New("Internal server error")

	/* This is a map of internal `error` to external API error.
	It's used to return the external API error based on the internal that was thrown.
	These two are seggregated - For each error there should be a conscious decision of what should be
	the internal and external user error.
	*/
	errorsMap = map[error]string{
		errInternalInvalidJSONBody:       `{"error": "Verify your JSON payload."}`,
		errInternalInvalidPath:           `{"error": "Verify your path (must start with "/" and end with "<filename>.txt")."}`,
		errInternalNoData:                `{"error": "There is no data to write to file."}`,
		errInternalInvalidFilename:       `{"error": "Invalid filename."}`,
		errInternalCreatingOrWritingFile: `{"error": "Error creating/writing file."}`,
		errInternalFileNotFound:          `{"error": "File not found."}`,
		errInternalPathNotFound:          `{"error": "Path not found."}`,
		errInternalGeneric:               `{"error": "Internal server error."}`,
	}

	successGetFileDataRsp = `{"data": "%v"}`
)

type newAndUpdateFileRequest struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

// ValidatedPath and ValidatedData are structs that hold a validated path/data. This is used in code to make sure that no `ServiceFile` gets created without the path/data being validated.
type ValidatedPath struct {
	Path string
}
type ValidatedData struct {
	Data string
}

func toServiceFile(vp ValidatedPath, vd ValidatedData) *ServiceFile {
	return NewServiceFile(vp, vd)
}

// StartAPI creates the API Endpoints and starts the API
func StartAPI(addr string) {
	r := mux.NewRouter()

	/* Returns the contents of a file.
	Query param is `path`: denotes the path to file and must end in <FILE_NAME>.txt
	Example:
	GET /file?path=/this/is/a/path/to/a/file.txt
	*/
	r.HandleFunc("/file", getFileAPIHandler).Methods("GET").Queries("path", "{path}")

	/* Removes a file and its contents from the system.
	Query param is `path`: denotes the path to file and must end in <FILE_NAME>.txt
	Example:
	DELETE /file?path=/please/delete/this/file.txt
	*/
	r.HandleFunc("/file", deleteFileAPIHandler).Methods("DELETE").Queries("path", "{path}")

	/* Creates/Updates a new file within the given path, with the given data.
	POST to /file. Body is in JSON and follows the following format:
	{
		"path": <PATH_TO_FILE> -> string (must end in `<FILE_NAME>.txt`)
		"data": <FILE_DATA> -> string
	}
	*/
	r.HandleFunc("/file", newFileAPIHandler).Methods("POST")

	/* Returns the statistics for a path.
	Query param is `path`: denotes the path to file and must end in <FILE_NAME>.txt
	Example:
	GET /file?path=/this/is/a/path/to/a/file.txt
	*/
	r.HandleFunc("/stats", getStatsAPIHandler).Methods("GET").Queries("path", "{path}")

	log.Fatal(http.ListenAndServe(addr, r))
}

func getFileAPIHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	vp, err := validatePath(path, true)
	if err != nil {
		log.Printf("Error validating path. Error: %v", err)
		http.Error(w, errorsMap[err], http.StatusBadRequest)
		return
	}

	// We don't know the data yet, so we pretend we do
	vd, err := validateData("Data")
	if err != nil {
		log.Printf("Error validating path. Error: %v", err)
		http.Error(w, errorsMap[err], http.StatusBadRequest)
		return
	}

	sf := toServiceFile(*vp, *vd)
	data, err := sf.GetFileData()
	if err != nil {
		log.Printf("Error getting file's data. Error: %v", err)
		// We assume, for external user error purposes that
		// if there is an error, most likely is that the file was not found
		http.Error(w, errorsMap[errInternalFileNotFound], http.StatusNotFound)
		return
	}

	s := fmt.Sprintf(successGetFileDataRsp, *data)
	w.Write([]byte(s))
	return
}

func deleteFileAPIHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	vp, err := validatePath(path, true)
	if err != nil {
		log.Printf("Error validating path. Error: %v", err)
		http.Error(w, errorsMap[err], http.StatusBadRequest)
		return
	}

	// We don't care about data, so we pretend we do
	vd, err := validateData("Data")
	if err != nil {
		log.Printf("Error validating path. Error: %v", err)
		http.Error(w, errorsMap[err], http.StatusBadRequest)
		return
	}

	sf := toServiceFile(*vp, *vd)
	err = sf.Delete()
	if err != nil {
		log.Printf("Error deleting file. Error: %v", err)
		// We assume, for external user error purposes that
		// if there is an error deleting a file, most likely is that the file was not found
		http.Error(w, errorsMap[errInternalFileNotFound], http.StatusNotFound)
		return
	}

	return
}

func newFileAPIHandler(w http.ResponseWriter, r *http.Request) {
	rqBody := new(newAndUpdateFileRequest)
	d := json.NewDecoder(r.Body)
	err := d.Decode(rqBody)
	if err != nil {
		log.Printf("Error decoding JSON Body. Error %v", err)
		http.Error(w, errorsMap[errInternalInvalidJSONBody], http.StatusBadRequest)
		return
	}

	validPath, err := validatePath(rqBody.Path, true)
	if err != nil {
		log.Printf("Error validating path. Error: %v", err)
		http.Error(w, errorsMap[err], http.StatusBadRequest)
		return
	}

	validData, err := validateData(rqBody.Data)
	if err != nil {
		log.Printf("Error validating path. Error: %v", err)
		http.Error(w, errorsMap[err], http.StatusBadRequest)
		return
	}

	file := toServiceFile(*validPath, *validData)
	err = file.Create()
	if err != nil {
		log.Printf("Error creating/writing to file. Error: %v", err)
		http.Error(w, errorsMap[errInternalCreatingOrWritingFile], http.StatusInternalServerError)
		return
	}

	return
}

func getStatsAPIHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	vp, err := validatePath(path, false)
	if err != nil {
		log.Printf("Error validating path. Error: %v", err)
		http.Error(w, errorsMap[err], http.StatusBadRequest)
		return
	}

	s, err := GetStats(vp.Path)
	if err != nil {
		log.Printf("Error getting stats. Error: %v", err)
		// We assume, for external user error purposes that
		// if there is an error here, most likely is that the path doesn't exist
		http.Error(w, errorsMap[errInternalPathNotFound], http.StatusNotFound)
		return
	}

	rsp, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling Stats struct. Error: %v", err)
		http.Error(w, errorsMap[errInternalGeneric], http.StatusInternalServerError)
		return
	}

	w.Write(rsp)
	return
}

func validatePath(path string, filenameCheck bool) (*ValidatedPath, error) {
	if path == "" {
		return nil, errInternalInvalidPath
	}

	// Every path must start with `/`
	if !strings.HasPrefix(path, "/") {
		return nil, errInternalInvalidPath
	}

	// Every filename must end in `.txt``
	if filenameCheck {
		splitPath := strings.Split(path, "/")
		filename := splitPath[len(splitPath)-1]
		if !strings.HasSuffix(filename, ".txt") {
			return nil, errInternalInvalidFilename
		}
		// But it cannot be only `.txt`
		if len(strings.TrimRight(filename, ".txt")) == 0 {
			return nil, errInternalInvalidFilename
		}
	}

	// Stops funny people from doing funny things :)
	path = filepath.Clean(path)

	return &ValidatedPath{Path: path}, nil
}

func validateData(data string) (*ValidatedData, error) {
	if data == "" {
		return nil, errInternalNoData
	}

	return &ValidatedData{Data: data}, nil
}
