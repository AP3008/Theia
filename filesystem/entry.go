// This file creates the struct of what each file will look like
package filesystem

import (
	"io/fs"
	"log"
	"os"
	"time"
	"path/filepath"
)

// This is how we will interact with each file / dir throughout the projet
type SystemFile struct{
	Name string
	Path string
	IsDir bool
	IsSymLink bool
	Size int64
	Permission os.FileMode
	ModifiedTime time.Time
}

//Simple way to get a formatted permission string

func (s SystemFile) FormatPermission() string{
	return s.Permission.String()
}

// From os.ReadDir() we get a []fs.DirEntry so we want to grab all necessary info from a singular DirEntry 

func CreateSystemFile(de os.DirEntry, parentDir string) (SystemFile, error){
	info, err := de.Info()
	if err != nil{
		log.Fatal(err)
	}
	fileName := de.Name()
	perms := info.Mode()
	fullPath := filepath.Join(parentDir, fileName)
	return SystemFile{
		Name : fileName,
		Path : fullPath,
		IsDir: de.IsDir(),
		IsSymLink: (perms & os.ModeSymlink) != 0,
		Size: info.Size(),
		Permission: perms,
		ModifiedTime: info.ModTime(),
	}, nil
}
