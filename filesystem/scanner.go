// This file handles scanning all of the files in a directory
package filesystem

import (
	"fmt"
	"os"
	"strings"
	fuzzy "github.com/sahilm/fuzzy"
)

// We Create a slice of SystemFiles to get a directory in a consumeable format
func CreateSystemFileList(path string, allFiles bool, fileMode bool, dirMode bool) ([]SystemFile, error) {
	itemList, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var directoryList []SystemFile
	for _, value := range itemList {
		if !allFiles && strings.HasPrefix(value.Name(), ".") {
			continue
		}
		sf, err := CreateSystemFile(value, path)
		if err != nil {
			fmt.Println("Error:", err, "File:", value.Name())
		}
		if fileMode && !sf.IsDir{
			directoryList = append(directoryList, sf)
		} else if dirMode && sf.IsDir{
			directoryList = append(directoryList, sf)
		} else if !dirMode && !fileMode{
			directoryList = append(directoryList, sf)
		}
	}
	sortDirList(directoryList)
	return directoryList, nil
}

// Making an custom type to implement Source interace 
type FileSource []SystemFile

func (f FileSource) String(i int) string{
	return f[i].Name
}

func (f FileSource) Len() int {
	return len(f)
}

func SearchSystemList(searchTerm string, sfl FileSource) ([]SystemFile, error){
	newList := fuzzy.FindFrom(searchTerm, sfl)	

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
			if name1 > name2 {
				tmp := dl[j]
				dl[j] = dl[j+1]
				dl[j+1] = tmp
			}
		}
	}
}
