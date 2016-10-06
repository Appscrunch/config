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

const chmodDir = 0755
const chmodFile = 0744

func GetUUID(DirWithUUID string) string {

	UUID, err := GetUUIDFromFile(DirWithUUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: Could not get uuid from file [%s]; creating new\n", err.Error())
	} else {
		return UUID
	}

	u, err := CreateNewUUID(DirWithUUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: Could not create uuid in file %s; creating <%s>\n", err.Error(), uuidFile)
	}
	u, err = CreateNewUUID(DirWithUUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: Could not create uuid in file %s; return just %s instead\n", err.Error(), u)
	}
	return u

}

func GetUUIDFromFile(DirWithUUID string) (string, error) {

	uuidPath, err := GetFileNameWithUUID(DirWithUUID)
	if err != nil {
		defaultLogF("[config]:  File with uuid [%s] does not exist; creating new\n", err.Error())
	}

	defaultLogF("[config] uuidPath=[%s]\n", uuidPath)

	//Reading and validate UUID
	u, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		defaultLogF("[config]: Could not read uuid from [%s]: [%s]\n", uuidPath, err.Error())
		return "", err
	}

	uValidated, err := uuid.FromString(string(u))
	if err != nil {
		defaultLogF("[config]: Could not validate uuid [%s] from file [%s]\n", string(u), uuidPath)
		return "", err
	}
	return uValidated.String(), nil
}

func GetFileNameWithUUID(DirWithUUID string) (string, error) {

	var fileWithUUIDExists = uuidFile

	//formatting if we had related path
	fpath, err := GetPathForDir(DirWithUUID)
	if err != nil {
		defaultLogF("[config] UUID dir: [%s]\n", err.Error())
		return "", err
	}

	// Does such path exists?
	fileWithUUIDExists = filepath.Join(fpath, uuidFile)
	exist, err := Exists(fileWithUUIDExists)
	if err != nil && !exist {
		defaultLogF("[config]: Wrong path for UUOD [%s]: [%s]\n", fileWithUUIDExists, err.Error())
		return "", err
	}

	return fileWithUUIDExists, nil
}

func CreateNewUUID(DirWithUUID string) (string, error) {

	uuidPath, err := GetFileNameWithUUID(DirWithUUID)
	if err != nil {
		defaultLogF("[config]:  File with uuid [%s] does not exist; creating new\n", err.Error())
	}

	defaultLogF("[config]: CreateNewUUID... uuidPath=[%s]\n", uuidPath)

	u := uuid.NewV4()
	uStr := u.String()

	//file exists
	exist, err := Exists(uuidPath)
	if err == nil && exist {
		err := ioutil.WriteFile(uuidPath, []byte(u.String()), chmodFile)
		if err == nil {
			return u.String(), nil
		}
	}

	// directory exists
	exist, err = Exists(DirWithUUID)
	if err == nil && exist {
		err := ioutil.WriteFile(uuidPath, []byte(u.String()), chmodFile)
		if err == nil {
			return u.String(), nil
		}
	}

	defaultLogF("[config]: Directory [%s] doesn't exist; creating...\n", DirWithUUID)
	err = os.MkdirAll(DirWithUUID, chmodDir)
	if err == nil {
		err := ioutil.WriteFile(uuidPath, []byte(u.String()), chmodFile)
		if err == nil {
			return u.String(), nil
		}
	}

	defaultLogF("[config]: Could not create file [%s]; searching near binary file\n", uuidPath)
	uStr, err = GetUUIDFromFile(uuidFile)
	if err == nil {
		return u.String(), nil
	}

	defaultLogF("[config]: File with uuid %s does not exist [%s]; creating new near binary file", err.Error())
	uStr, err = CreateNewUUID(uuidFile)
	if err == nil {
		return uStr, nil
	}

	defaultLogF("[config]: Could not create uuid in file [%s]; return just u4=[%s] instead\n", err.Error(), u.String())
	return u.String(), err
}
