package fileutils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb"
)

// CopyFile copies a file from source to dest
func CopyFile(source string, dest string) error {
	fi, err := os.Stat(source)
	if fi == nil {
		return fmt.Errorf("could not access file %s: %v", source, err)
	}

	bar := pb.New(int(fi.Size())).SetUnits(pb.U_BYTES)
	bar.ShowPercent = true
	bar.ShowCounters = true
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	fmt.Printf("### File Copy Task\n#\tSRC: %s\n#\tDST: %s\n", source, dest)
	bar.Start()

	src, err := os.Open(source)
	if err != nil {
		return err
	}

	reader := bar.NewProxyReader(src)

	defer src.Close()

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, reader)
	if err != nil {
		return err
	}

	info, err := os.Stat(source)
	if err != nil {
		err = os.Chmod(dest, info.Mode())
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyDir recursively copies a directory
func CopyDir(source string, dest string) error {
	srcinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, srcinfo.Mode())
	if err != nil {
		return err
	}

	dir, _ := os.Open(source)
	obs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	var errs []error

	for _, obj := range obs {
		fsource := filepath.Join(source, obj.Name())
		fdest := filepath.Join(dest, obj.Name())

		if obj.IsDir() {.
			err = CopyDir(fsource, fdest)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			err = CopyFile(fsource, fdest)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	var errString string
	for _, err := range errs {
		errString += err.Error() + "\n"
	}

	if errString != "" {
		return errors.New(errString)
	}

	return nil
}
