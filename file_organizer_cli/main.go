package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var category = map[string]string{
	".jpg":  "Images",
	".jpeg": "Images",
	".png":  "Images",
	".gif":  "Images",
	".mp4":  "Videos",
	".avi":  "Videos",
	".mov":  "Videos",
	".pdf":  "Docs",
	".docx": "Docs",
	".txt":  "Docs",
	".mp3":  "Music",
	".wav":  "Music",
}

func copyFile(ext, name, base string) error {
	fileType := category[ext]
	if fileType == "" {
		return nil
	}
	Oname := filepath.Join(base, name)
	fmt.Println(fileType, " ", Oname)
	sourceFile, err := os.Open(Oname)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	os.MkdirAll(filepath.Join(base, fileType), 0755)
	dest := filepath.Join(base, fileType, name)
	destFile, err := os.Create((dest))

	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	err = destFile.Sync()
	if err != nil {
		return err
	}
	return nil
}
func deleteFile(name, path string) error {
	Oname := filepath.Join(path, name)
	err := os.Remove(Oname)
	if err != nil {
		fmt.Println("Failed to delete:", err)
	} else {
		fmt.Println("Original deleted")
	}

	return nil
}

func findFile(path string) {
	entries, _ := os.ReadDir(path)
	for _, e := range entries {
		if !e.IsDir() {
			ext := filepath.Ext(e.Name())
			copyFile(ext, e.Name(), path)
			deleteFile(e.Name(), path)
		}
	}
}

func main() {
	src := flag.String("src", "", "Source file path")
	flag.Parse()
	findFile(*src)
}
