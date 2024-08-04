package helper

import "os"


func CreateDirectoryIfNotExists(dirctoryName string) error {
    var err = os.MkdirAll(dirctoryName, os.ModeDir)

    if err == nil || os.IsExist(err) {
        return nil
    } else {
        return err
    }
}