package cache

import (
	"ecm-sdk-go/constants"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func GetFileName(cacheDir, cacheFilePrefix string) string {
	return cacheDir + string(os.PathSeparator) + cacheFilePrefix + "_" + constants.CachFileName
}

func WriteConfigToFile(cacheDir, cacheFilePrefix, content string) {
	mkdirIfNecessary(cacheDir)
	fileName := GetFileName(cacheDir, cacheFilePrefix)
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		log.Printf("[ERROR]:faild to write config  cache:%s ,value:%s ,err:%s \n", fileName, string(content), err.Error())
	}
}

func ReadConfigFromFile(cacheDir, cacheFilePrefix string) (string, error) {
	fileName := GetFileName(cacheDir, cacheFilePrefix)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to read config cache file:%s,err:%s! ", fileName, err.Error()))
	}
	return string(b), nil
}
