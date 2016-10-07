package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

const uuidFile = "local.conf"

const chmodDir = 0755
const chmodFile = 0644

// getUUIDFromFile read and validate uuid from file
func getUUIDFromFile(DirWithUUID string) (string, error) {

	uuidPath, err := getFileNameWithUUID(DirWithUUID)
	if err != nil {
		defaultLogF("[config]: guid: File with uuid [%s] does not exist; creating new\n", err.Error())
	}

	defaultLogF("[config] guid: uuidPath=[%s]\n", uuidPath)

	//Reading and validate UUID
	u, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		defaultLogF("[config]: guid: Could not read uuid from [%s]: [%s]\n", uuidPath, err.Error())
		return "", err
	}

	uValidated, err := uuid.FromString(string(u))
	if err != nil {
		defaultLogF("[config]: guid: Could not validate uuid [%s] from file [%s]\n", string(u), uuidPath)
		return "", err
	}
	return uValidated.String(), nil
}

func getFileNameWithUUID(DirWithUUID string) string {

	//var fileWithUUIDExists = uuidFile

	//formatting if we had related path
	fpath, err := GetPathForDir(DirWithUUID)
	if err != nil {
		defaultLogF("[config] guid: UUID dir: [%s]\n", err.Error())
		//	return "", err
	}

	// Does such path exists?
	filename := filepath.Join(fpath, uuidFile)
	/*exist, err := Exists(filename)
	if err != nil && !exist {
		defaultLogF("[config]: guid: Wrong path for UUOD [%s]: [%s]\n", filename, err.Error())
		return "", err
	}*/

	return filename
}

func createNewUUID(DirWithUUID string) (string, error) {

	uuidPath := getFileNameWithUUID(DirWithUUID)
	if err != nil {
		defaultLogF("[config]: guid: File with uuid [%s] does not exist; creating new\n", err.Error())
	}

	defaultLogF("[config]: guid: CreateNewUUID... uuidPath=[%s]\n", uuidPath)

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

	defaultLogF("[config]: guid: Directory [%s] doesn't exist; creating...\n", DirWithUUID)
	err = os.MkdirAll(DirWithUUID, chmodDir)
	if err == nil {
		err := ioutil.WriteFile(uuidPath, []byte(u.String()), chmodFile)
		if err == nil {
			return u.String(), nil
		}
	}

	defaultLogF("[config]: guid: Could not create file [%s]; searching near binary file\n", uuidPath)
	uStr, err = getUUIDFromFile(uuidFile)
	if err == nil {
		return u.String(), nil
	}

	defaultLogF("[config]: guid: File with uuid %s does not exist [%s]; creating new near binary file", err.Error())
	uStr, err = createNewUUID(uuidFile)
	if err == nil {
		return uStr, nil
	}

	defaultLogF("[config]: guid: Could not create uuid in file [%s]; return just u4=[%s] instead\n", err.Error(), u.String())
	return u.String(), err
}

// GetUUID read or create UUID from local.conf file
func GetUUID(DirWithUUID string) string {

	UUID, err := getUUIDFromFile(DirWithUUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: guid: Could not get uuid from file [%s]; creating new\n", err.Error())
	} else {
		defaultLogF("[config] guid: uuid=[%s]\n", UUID)
		return UUID
	}

	UUID, err = createNewUUID(DirWithUUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: guid: Could not create uuid in file %s; creating <%s>\n", err.Error(), uuidFile)
	}
	UUID, err = createNewUUID(DirWithUUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] [config]: guid: Could not create uuid in file %s; return just %s instead\n", err.Error(), UUID)
	}
	defaultLogF("[config] guid: uuid=[%s]\n", UUID)
	return UUID

}
