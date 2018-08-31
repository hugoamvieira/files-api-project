package main

import (
	"io/ioutil"
	"os"
	"strings"
)

const (
	// DefaultFolderPermission is the permission flag which is used to create all folders
	DefaultFolderPermission = 0755
	// DefaultFilePermission is the permission flag which is used to create all files
	DefaultFilePermission = 0644
	// RootPath is the root folder for this service
	RootPath = "filesys"
)

// ServiceFile is the internal struct that holds the information about a file
type ServiceFile struct {
	name string
	path string
	data string
}

// NewServiceFile creates a new ServiceFile object, which is the service's representation of a file
func NewServiceFile(vp ValidatedPath, vd ValidatedData) *ServiceFile {
	splitPath := strings.Split(vp.Path, "/")
	filename := splitPath[len(splitPath)-1]
	pathToFile := strings.Join(splitPath[:len(splitPath)-1], "/")

	return &ServiceFile{
		name: filename,
		path: RootPath + pathToFile, // Makes it so that no file will never be read/written to outside rootpath.
		data: vd.Data,
	}
}

// GetFileData takes a path and retrieves the contents of a file if it exists
func (sf *ServiceFile) GetFileData() (*string, error) {
	fileBytes, err := ioutil.ReadFile(sf.path + "/" + sf.name)
	if err != nil {
		return nil, err
	}

	data := string(fileBytes)
	return &data, nil
}

// Create creates a file and its path and writes to it
func (sf *ServiceFile) Create() error {
	// TODO: Check if path already exists
	err := sf.createOSPath()
	if err != nil {
		return err
	}

	osFile, err := sf.createOSFile()
	if err != nil {
		return err
	}

	err = sf.writeToOSFile(osFile)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a file from the filesystem
func (sf *ServiceFile) Delete() error {
	err := os.Remove(sf.constructFullPath())
	if err != nil {
		return err
	}

	return nil
}

func (sf *ServiceFile) createOSPath() error {
	err := os.MkdirAll(sf.path, DefaultFolderPermission)
	if err != nil {
		return err
	}
	return nil
}

func (sf *ServiceFile) createOSFile() (*os.File, error) {
	osFile, err := os.Create(sf.constructFullPath())
	if err != nil {
		return nil, err
	}
	return osFile, nil
}

func (sf *ServiceFile) writeToOSFile(osFile *os.File) error {
	_, err := osFile.Write([]byte(sf.data))
	if err != nil {
		return err
	}
	return nil
}

func (sf *ServiceFile) constructFullPath() string {
	return sf.path + "/" + sf.name
}
