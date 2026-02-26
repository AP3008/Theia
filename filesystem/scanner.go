// This file handles scanning all of the files in a directory
package filesystem

import (
	"fmt"
	"os"
)

// We Create a slice of SystemFiles to get a directory in a consumeable format
func CreateSystemFileList (path string) ([]SystemFile, error){
	itemList, err := os.ReadDir(path)
	if err != nil{
		return nil, err
	}
	var directoryList []SystemFile
	for _, value := range itemList{
		sf, err := CreateSystemFile(value, path)	
		if err != nil{
			fmt.Println("Error:", err, "File:", value.Name())
		}
		directoryList = append(directoryList, sf)
	}
	return directoryList, nil
}
