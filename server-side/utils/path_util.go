package utils

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type PathUtil struct {
	rootDir string
	fileDir string
	tempDir string
}

func NewPathUtil(rootDir string) *PathUtil {
	fileDir := filepath.Join(rootDir, ".cloud_drive_files")
	tempDir := filepath.Join(fileDir, "temp")
	if err := CreateDir(fileDir); err != nil {
		logrus.Errorf("Failed to create file directory: %v", err)
	}
	if err := CreateDir(tempDir); err != nil {
		logrus.Errorf("Failed to create temp directory: %v", err)
	}

	return &PathUtil{
		rootDir: rootDir,
		fileDir: fileDir,
		tempDir: tempDir,
	}
}

func (pu *PathUtil) GetRootDir() string {
	return pu.rootDir
}

func (pu *PathUtil) GetFileDir() string {
	return pu.fileDir
}

func (pu *PathUtil) GetTempDir() string {
	return pu.tempDir
}
