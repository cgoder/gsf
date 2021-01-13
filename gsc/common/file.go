package common

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func FilePathCreate(dirPath string, GUID string) error {
	// Get local destination path.
	var tmpDir string
	if GUID == "" {
		tmpDir = dirPath
	} else {
		tmpDir = dirPath + "/" + GUID
	}

	err := os.MkdirAll(tmpDir, 0700)
	if err != nil {
		return err
	}
	return nil
}

func FileDel(filePath string) (bool, error) {
	log.Info("del file: ", filePath)

	err := os.RemoveAll(filePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 判断所给路径文件/文件夹是否存在
func FilePathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断所给路径是否为文件夹
func IsDir(path string) (bool, error) {
	s, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return s.IsDir(), nil
}
