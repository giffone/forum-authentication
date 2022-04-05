package paths

import (
	"log"
	"os"
)

func CreatePaths(path, pathBackup string) {
	if NotExist(path) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatalf("create path: can not create pathdb - %v\n", err)
		}
	}
	if NotExist(pathBackup) {
		err := os.Mkdir(pathBackup, os.ModePerm)
		if err != nil {
			log.Fatalf("create pathBackup: can not create pathBackup - %v\n", err)
		}
	}
}

func NotExist(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}
