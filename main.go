package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"strings"

	"github.com/client9/xson/hjson"
	goftp "github.com/n-marshall/ftp"
	rndm "github.com/n-marshall/rndm-go"
)

var (
	chars = "1234567890QWERTYUIOPMLKJHGFDSAZXCVBNqwertyuiopmlkjhgfdsazxcvbn"
)

type Job struct {
	SrcPath              string
	DestPath             string
	AuthorizedExtensions string
	DestFolderSize       int
	NbFilesToLeave       int
}

type Config struct {
	FtpServer string
	Jobs      []Job
}

func main() {
	config := loadConfig()

	ftp, err := goftp.Connect(config.FtpServer)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		errq := ftp.Quit()
		if errq != nil {
			panic(errq)
		}
	}()

	if err = ftp.Login("username", "password"); err != nil {
		log.Fatal(err)
	}

	for _, job := range config.Jobs {
		baseSrcPath := path.Clean(job.SrcPath)
		baseDestPath := path.Clean(job.DestPath)

		fileNames, err := ftp.NameList(baseSrcPath)
		if err != nil {
			log.Panic(fmt.Errorf("%s: %v", baseSrcPath, err))
		}

		var nbFiles = len(fileNames)
		var folders []string
		if job.DestFolderSize > 0 {
			nbDestFolders := nbFiles/job.DestFolderSize + 1
			for i := 0; i < nbDestFolders; i++ {
				folders = append(folders, rndm.String(8))
			}
		}

		for idx, fileName := range fileNames {
			suffixes := strings.Split(job.AuthorizedExtensions, ",")
			if hasOneSuffix(fileName, suffixes) {
				if idx < nbFiles-job.NbFilesToLeave || job.NbFilesToLeave == 0 {
					srcPath := path.Join(baseSrcPath, fileName)
					srcFile, err := ftp.Retr(srcPath)
					if err != nil {
						log.Panic(err)
					}
					var srcFile2 bytes.Buffer
					tee := io.TeeReader(srcFile, &srcFile2)

					var srcHasher = sha256.New()
					if _, err = io.Copy(srcHasher, tee); err != nil {
						log.Fatal(err)
					}
					srcFile.Close()
					srcHash := fmt.Sprintf("%s %x", fileName, string(srcHasher.Sum(nil)))

					var destFolderPath, destFilePath string
					if job.DestFolderSize > 0 {
						destFolderPath = path.Join(baseDestPath, folders[idx/job.DestFolderSize])
						destFilePath = path.Join(destFolderPath, fileName)
					} else {
						destFolderPath = baseDestPath
						destFilePath = path.Join(baseDestPath, fileName)
					}

					// Open the destination file for writing
					err = os.MkdirAll(destFolderPath, 0777)
					if err != nil {
						panic(err)
					}
					dstFile, err := os.Create(destFilePath)
					if err != nil {
						panic(err)
					}

					if _, err = io.Copy(dstFile, &srcFile2); err != nil {
						panic(err)
					}
					err = dstFile.Sync()
					cerr := dstFile.Close()
					if err == nil {
						err = cerr
					}

					// check hash of dest file
					dstFile2, err := os.Open(destFilePath)
					if err != nil {
						log.Fatal(err)
					}

					var dstHasher = sha256.New()
					if _, err = io.Copy(dstHasher, dstFile2); err != nil {
						log.Fatal(err)
					}
					dstFile2.Close()
					dstHash := fmt.Sprintf("%s %x", fileName, string(dstHasher.Sum(nil)))
					if srcHash != dstHash {
						log.Fatal(fmt.Errorf("src: %s \ndst: %s", srcHash, dstHash))
					}

					err = ftp.Delete(srcPath)
					if err != nil {
						panic(err)
					}
					fmt.Println("done: " + srcPath)
				}
			}
		}
	}
}

func loadConfig() Config {
	configFileBytes, err := ioutil.ReadFile("config.hjson")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := hjson.Unmarshal(configFileBytes, &config); err != nil {
		panic(err)
	}
	return config
}

func hasOneSuffix(s string, suffixList []string) bool {
	for _, suffix := range suffixList {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}
