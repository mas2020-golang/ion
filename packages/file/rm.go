package file

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	out "github.com/mas2020-golang/goutils/output"
)

type Eraser struct {
	dryRun        bool
	Dirs          int
	SpaceFree     int64
	Files         int
	deep, verbose bool
}

func NewEraser(dryRun, deepClean, verbose bool) *Eraser {
	return &Eraser{
		dryRun:  dryRun,
		deep:    deepClean,
		verbose: verbose,
	}
}

func (e *Eraser) Delete(path string) error {
	return e.readPath(path)
}

// readPath reads the path and in case it is a file delete it, otherwise search into the nested
// folders
func (e *Eraser) readPath(path string) error {
	//var symbol string = elemChar
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat error for the path %s, error: %v", path, err)
	}

	if fi.IsDir() {
		e.Dirs++
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("ReadDir error for the path %s, error: %v", path, err)
		}

		for _, fi := range fis {
			if err := e.readPath(filepath.Join(path, fi.Name())); err != nil {
				return err
			}
		}
		e.deleteFolder(path)
	} else {
		e.SpaceFree += fi.Size()
		if err := e.deleteFile(path); err != nil {
			return err
		}
	}
	return nil
}

// deleteFile deletes the given file
func (e *Eraser) deleteFile(f string) error {
	e.Files++
	if e.dryRun {
		e.message(out.RedS(fmt.Sprintf("the file %s would be deleted", f)))
		return nil
	}

	// file deletion
	if e.deep {
		if err := e.deepClean(f); err != nil {
			return fmt.Errorf("error deleting the file %s, error: %v", f, err)
		}
		e.message(out.RedS(fmt.Sprintf("the file %s has been wiped out", f)))
	}

	if err := os.Remove(f); err != nil {
		return fmt.Errorf("error deleting the file %s, error: %v", f, err)
	}
	e.message(out.RedS(fmt.Sprintf("the file %s has been deleted", f)))

	return nil
}

// deleteFolder deletes the given folder
func (e *Eraser) deleteFolder(f string) error {
	if e.dryRun {
		e.message(out.BlueS(fmt.Sprintf("the folder %s would be deleted", f)))
		return nil
	}

	// folder deletion
	if err := os.Remove(f); err != nil {
		return err
	}

	e.message(out.BlueS(fmt.Sprintf("deleted the folder %s", f)))
	return nil
}

func (e *Eraser) message(t string) {
	if e.verbose || e.dryRun {
		fmt.Println(t)
	}
}

// deepClean removes the file in depth
func (e *Eraser) deepClean(path string) error {
	// open the file
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	defer func(*os.File) {
		if err := f.Close(); err != nil {
			out.Error("deepClean", err.Error())
		}
	}(f)

	if err != nil {
		return fmt.Errorf("error opening the file %s: %s", path, err.Error())
	}

	// stat the file
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("error opening the file %s: %s", path, err.Error())
	}

	// slice for the size of the file
	var fSize int64 = fi.Size()
	const chunk = 2 * (1 << 20) // 2 MB

	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fSize) / float64(chunk)))

	lastPosition := 0

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(chunk, float64(fSize-int64(i*chunk))))
		partZeroBytes := make([]byte, partSize)

		// fill out the part with zero value
		copy(partZeroBytes[:], "0")

		// over write every byte in the chunk with 0
		_, err := f.WriteAt([]byte(partZeroBytes), int64(lastPosition))

		if err != nil {
			return fmt.Errorf("error wiping 0 to the file %s: %s", path, err.Error())
		}

		// update last written position
		lastPosition = lastPosition + partSize
	}

	return nil
}
