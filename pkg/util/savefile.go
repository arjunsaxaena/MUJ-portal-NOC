package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

const uploadPath = "../uploads"

func SaveFile(fileHeader *multipart.FileHeader, folder, regNumber string) string {
	dstDir := filepath.Join(uploadPath, folder)
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return ""
	}

	dstPath := filepath.Join(dstDir, fmt.Sprintf("%s_%s", regNumber, fileHeader.Filename))
	dstFile, err := os.Create(dstPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return ""
	}
	defer dstFile.Close()

	srcFile, err := fileHeader.Open()
	if err != nil {
		fmt.Printf("Error opening uploaded file: %v\n", err)
		return ""
	}
	defer srcFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return ""
	}

	return folder + "/" + filepath.Base(dstPath)
}
