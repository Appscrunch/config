package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	//"github.com/KristinaEtc/config"

	uuid "github.com/satori/go.uuid"
)

const uuidFile = "local.conf"

const chmodFile = 0755

func GetUUID(DirWithUUID string) string {

	uuidPath, err := GetFileNameWithUUID(DirWithUUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: File with uuid %s does not exist; creating new\n", err.Error())
	} else {
		defaultLogF("[config] uuidPath=%s\n", uuidPath)
		UUID, err := GetUUIDFromFile(uuidPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] [config]: Could not get uuid from file [%s]; creating new\n", err.Error())
		} else {
			return UUID
		}
	}

	u, err := CreateNewUUID(uuidPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: Could not create uuid in file %s; creating <%s>\n", err.Error(), uuidFile)
	}
	u, err = CreateNewUUID(uuidPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: Could not create uuid in file %s; return just %s instead\n", err.Error(), u)
	}
	return u

}

func GetFileNameWithUUID(DirWithUUID string) (string, error) {

	var fileWithUUIDExists = uuidFile

	//formatting if we had related path
	fpath, err := GetPathForDir(DirWithUUID)
	if err != nil {
		defaultLogF("[config] UUID dir: %s\n", err.Error())
		return "", err
	}

	// Does such path exists?
	fileWithUUIDExists = filepath.Join(fpath, uuidFile)
	exist, err := Exists(fileWithUUIDExists)
	if err != nil && !exist {
		defaultLogF("[config]: Wrong path for UUOD (%s): %s\n", fileWithUUIDExists, err.Error())
		return "", err
	}

	return fileWithUUIDExists, nil
}

func GetUUIDFromFile(uuidPath string) (string, error) {

	//Reading and validate UUID
	u, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		defaultLogF("[config]: Could not read uuid from %s\n", uuidPath)
		return "", err
	}

	uValidated, err := uuid.FromString(string(u))
	if err != nil {
		defaultLogF("[config]: Could not validate uuid %s from file %s\n", string(u), uuidPath)
		return "", err
	}
	return uValidated.String(), nil
}

func CreateNewUUID(uuidPath string) (string, error) {

	defaultLogF("[config]: CreateNewUUID... fileWithUUIDExists=%s\n", uuidPath)

	u := uuid.NewV4()
	uStr := u.String()

	exist, err := Exists(uuidPath)
	if err == nil && exist {
		err := ioutil.WriteFile("uuidPath", u.Bytes(), chmodFile)
		if err == nil {
			return u.String(), nil
		}
	}

	defaultLogF("[config]: directories %s doesn't exist; creating...\n", uuidPath)
	err = os.MkdirAll(uuidPath, chmodFile)
	if err == nil {
		err := ioutil.WriteFile("uuidPath", u.Bytes(), chmodFile)
		if err == nil {
			return u.String(), nil
		}
	}

	defaultLogF("[config] could not create file %s; searching near binary file\n", uuidPath)
	uStr, err = GetUUIDFromFile(uuidFile)
	if err == nil {
		return u.String(), nil
	}

	defaultLogF("[config]: File with uuid %s does not exist [%s]; creating new near binary file", err.Error())
	uStr, err = CreateNewUUID(uuidFile)
	if err == nil {
		return uStr, nil
	}

	defaultLogF("[config]: Could not create uuid in file %s; return just u4=%s instead\n", err.Error(), u.String())
	return u.String(), err
}
