package path

import (
	"log"
	"os"
)

const (
	TEST_MAPBASE_PATH = "TEST_MAPBASE_PATH"
)
// Root current dir
var Root string

func init() {
	log.Println("tg-core path init start...")

	if rootPath := os.Getenv(TEST_MAPBASE_PATH);rootPath != ""{
		Root = rootPath
	}else {
		var err error
		Root, err = os.Getwd()
		if err != nil {
			log.Fatal("Initialize Root error: ", err)
		}
	}
	log.Println("tg-core path init successful, path="+Root)
}