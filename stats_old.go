package main

/*
I'll leave this here since I did spend time on it ¯\_(ツ)_/¯


import (
	"io/ioutil"
	"os"
	"strings"
)

type WalkFunction func(path string, fi []os.FileInfo) error

// GetStatsOld populates the `Stats` struct with data and returns it to the caller
func GetStatsOld(path string) (*Stats, error) {
	s := new(Stats)
	if strings.HasSuffix(path, "/") {
		path = strings.TrimRight(path, "/")
	}

	err := walkReverse(path, s.walkFunction)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Stats) walkFunction(path string, fi []os.FileInfo) error {
	// TODO
	return nil
}

// This is my version of filepath.Walk(), but in reverse
func walkReverse(path string, f WalkFunction) error {
	path = RootPath + path

	sp := strings.Split(path, "/")
	for i := len(sp) - 1; i >= 0; i-- {
		// Visit folder
		currentPath := strings.Join(sp, "/")
		fi, err := ioutil.ReadDir(currentPath)
		if err != nil {
			return err
		}

		// Call function
		err = f(currentPath, fi)
		if err != nil {
			return err
		}

		// Pop last element
		sp = sp[:len(sp)-1]
	}
	return nil
}
*/
