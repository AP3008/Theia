// This file creates the struct of what each file will look like
package filesystem

import (
	"fmt"
	"os"
	"time"
)

type SystemFile struct{
	Name string
	Path string
	IsDir bool
	Size int64
	Permission uint32
	ModifiedTime time.Time
}
