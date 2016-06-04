package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

func main() {
	dir := "./"
	renLogFileName := ".renlog"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logs.Errorf("Failed to read dir %s, the error is %v", dir, err)
		return
	}
	renameLog := ""
	for _, f := range files {
		fileName := f.Name()
		if fileName == renLogFileName {
			continue
		}
		if fileName[:1] == "." {
			logs.Debugf("Skip hidden file %s", fileName)
			continue
		}
		newFileName := utils.GetMD5Hash(fileName)
		if !f.IsDir() {
			dotIndex := strings.LastIndex(fileName, ".")
			if dotIndex != -1 && dotIndex != 0 {
				newFileName += fileName[dotIndex:]
			}
		}
		err = os.Rename(fileName, newFileName)
		if err != nil {
			logs.Errorf("Failed to rename file %s to %s, the error is %v", fileName, newFileName, err)
			continue
		}
		renameLog += fmt.Sprintf("%s\t%s\n", fileName, newFileName)
	}
	fmt.Println(renameLog)
	ioutil.WriteFile(renLogFileName, []byte(renameLog), 0666)
}
