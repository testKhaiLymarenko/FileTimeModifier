package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func main() {

	defer func() {
		fmt.Println("\n\nPress any key to continue ...")
		fmt.Scanln()
	}()

	var userPath string
	//userPath := "I:\\П  Р  И  Р  О   Д  А"

	for {
		fmt.Print("Path: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if scanner.Err() != nil {
			fmt.Println("err")

			return
		}

		userPath = scanner.Text()

		if !fileExists(userPath) {
			fmt.Println("\tSuch directory doesn't exist. Input path again")
			continue
		}

		break
	}

	start := time.Now()
	fileCount := 0

	err := filepath.Walk(userPath, func(path string, file fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println()
			fmt.Println(err)
		} else {
			if userPath != path {
				if !file.IsDir() {
					if perr := setCreationTime(path, file.ModTime()); perr != nil {
						fmt.Println(err)
					}

					fileCount++
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n\nFileTimeModified => FileTimeCreated = success!\tFiles: %d | Elapsed time: %v\n",
		fileCount, time.Since(start))
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func setCreationTime(path string, ctime time.Time) error {
	ctimespec := syscall.NsecToTimespec(ctime.UnixNano())
	pathp, e := syscall.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	h, e := syscall.CreateFile(pathp,
		syscall.FILE_WRITE_ATTRIBUTES, syscall.FILE_SHARE_WRITE, nil,
		syscall.OPEN_EXISTING, syscall.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil {
		return e
	}
	defer syscall.Close(h)
	c := syscall.NsecToFiletime(syscall.TimespecToNsec(ctimespec))
	return syscall.SetFileTime(h, &c, nil, nil)
}

/*
	type Win32FileAttributeData struct {
	FileAttributes uint32
	CreationTime   syscall.Filetime
	LastAccessTime syscall.Filetime
	LastWriteTime  syscall.Filetime
	FileSizeHigh   uint32
	FileSizeLow    uint32
	}

	d := file.Sys().(*syscall.Win32FileAttributeData)
	cTime := time.Unix(0, d.CreationTime.Nanoseconds())

	fd, err := syscall.Open(path, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println(path, err)
	}

	err = syscall.SetFileTime(fd, &d.LastWriteTime, &d.LastAccessTime, &d.LastWriteTime)
	fmt.Println(err)

	fmt.Println("ctime: %v", cTime)
*/
