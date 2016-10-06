package config

import (
	"io/ioutil"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
	"github.com/ventu-io/slf"
)

var pwdCurr = "KristinaEtc/config"
var log = slf.WithContext(pwdCurr)

func GetUUID(DirWithUUID string) string {

	uuidPath, err := GetFileNameWithUUID(DirWithUUID)
	if err != nil {
		log.Errorf("File with uuid %s does not exist; creating new", err.Error())
	} else {
		UUID, err := GetUUIDFromFile(uuidPath)
		if err != nil {
			log.Errorf("Could not get uuid from file %s; getting from <local.conf>", err.Error())
			UUID, err = GetUUIDFromFile("local.conf")
			if err != nil {
				log.Errorf("File with uuid %s does not exist; creating new", err.Error())
			}
		} else {
			return UUID
		}
	}

	u, err := CreateNewUUID(uuidPath)
	if err != nil {
		log.Errorf("Could not create uuid in file %s; creating <local.conf>", err.Error())
	}
	u, err = CreateNewUUID(uuidPath)
	if err != nil {

	}
	return u

}

func GetFileNameWithUUID(DirWithUUID string) (string, error) {

	var fileWithUUIDExists = "local.conf"

	//formatting if we had related path
	fpath, err := GetPathForDir(DirWithUUID)
	if err != nil {
		log.Errorf("UUID dir: %s", err.Error())
		return "", err
	}

	// Does such path exists?
	fileWithUUIDExists = filepath.Join(fpath, "local.conf")
	exist, err := Exists(fileWithUUIDExists)
	if err != nil && !exist {
		log.Errorf("Wrong path for UUOD (%s): %s", fileWithUUIDExists, err.Error())
		return "", err
	}

	return fileWithUUIDExists, nil
}

func GetUUIDFromFile(uuidPath string) (string, error) {

	//Reading and validate UUID
	u, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		log.Errorf("Could not read uuid from %s file", uuidPath)
		return "", err
	}

	uValidated, err := uuid.FromString(string(u))
	if err != nil {
		log.Errorf("Could not validate uuid %s from file %s", string(u), uuidPath)
		return "", err
	}
	return uValidated.String(), nil
}

func CreateNewUUID(uuidPath string) (string, error) {

	/*	fileWithUUIDExists := filepath.Join(glo, UUIDfilename)

		log.Debugf("fileWithUUIDExists=%s", fileWithUUIDExists)
		exist, err := config.Exists(fileWithUUIDExists)
		if err != nil && !exist {
			log.Errorf("Wrong path for UUOD: %s;  searching near binary file", fileWithUUIDExists)
			fileWithUUIDExists = UUIDfilename
		}

		d1 := []byte("hello\ngo\n")
		err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
		check(err)
	*/
	return "", nil
}
