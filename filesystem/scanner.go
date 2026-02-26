// This file handles scanning all of the files in a directory
package filesystem

import (
	"fmt"
	"os"
	"strings"
)

// We Create a slice of SystemFiles to get a directory in a consumeable format
func CreateSystemFileList(path string) ([]SystemFile, error) {
	itemList, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var directoryList []SystemFile
	for _, value := range itemList {
		sf, err := CreateSystemFile(value, path)
		if err != nil {
			fmt.Println("Error:", err, "File:", value.Name())
		}
		directoryList = append(directoryList, sf)
	}
	sortDirList(directoryList)
	return directoryList, nil
}

// Private helper to help sort the dir list
// Because we are using slices we don't actually need to pass a reference
func sortDirList(dl []SystemFile) {
	n := len(dl)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			name1 := strings.ToLower(dl[j].Name)
			name2 := strings.ToLower(dl[j+1].Name)
			// Go checks every letter in this comparison
			if name1 < name2 {
				tmp := dl[j]
				dl[j] = dl[j+1]
				dl[j+1] = tmp
			}
		}
	}
}
