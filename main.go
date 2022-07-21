package main // import "hello-ftp"

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"github.com/secsy/goftp"
	"golang.org/x/crypto/ssh"
)

type QueSheet struct {
	Name  string
	IsDIR bool
}

var que []QueSheet

var (
	// srcBase = "hihi/haha/samples"
	srcBase = ""
	dstBase = ""

	replacerSlash = strings.NewReplacer("\\", string(os.PathSeparator), "/", string(os.PathSeparator))
)

// Upload file to sftp server
func sftpUploadFile(sc *sftp.Client, localFile, remoteFile string) (err error) {
	// log.Printf("Uploading '%s' to '%s' ..", localFile, remoteFile)

	srcFile, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("Unable to open local file: %v", err)
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		path = strings.ReplaceAll(path, "\\", "/")
		sc.Mkdir(path)
	}

	// Note: SFTP Go doesn't support O_RDWR mode
	dstFile, err := sc.OpenFile(remoteFile, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		return fmt.Errorf("Unable to open remote file: %v", err)
	}
	defer dstFile.Close()

	// bytes, err := io.Copy(dstFile, srcFile)
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("Unable to upload local file: %v", err)
	}
	// log.Printf("%d bytes copied", bytes)

	return nil
}

// Upload file to ftp server
func ftpUploadFile(fc *goftp.Client, localFile, remoteFile string) (err error) {
	// log.Printf("Uploading '%s' to '%s' ..", localFile, remoteFile)

	srcFile, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("Unable to open local file: %v", err)
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		path = strings.ReplaceAll(path, "\\", "/")
		fc.Mkdir(path)
	}

	err = fc.Store(remoteFile, srcFile)
	if err != nil {
		return fmt.Errorf("Unable to upload local file: %v", err)
	}
	// log.Printf("%d bytes copied", bytes)

	return nil
}

// walkDIR - walk and upload files
func walkDIR(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	que = append(que, QueSheet{Name: path, IsDIR: info.IsDir()})

	return nil
}

func sftpMain(addr, id, pwd string) {
	var sshConfig = &ssh.ClientConfig{
		User:            id,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
	}

	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	// log.Println("Connected.")

	// open an SFTP session over an existing ssh connection.
	sc, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	srcBase = replacerSlash.Replace(srcBase)
	srcRoot := filepath.Base(srcBase)
	srcCutPath := replacerSlash.Replace(strings.TrimSuffix(srcBase, srcRoot))

	err = filepath.Walk(srcBase, walkDIR)
	if err != nil {
		log.Println(err)
	}

	for _, q := range que {
		srcPath := filepath.Join("", q.Name)
		dstPath := filepath.Join(dstBase, strings.TrimPrefix(q.Name, srcCutPath))
		dstPath = strings.ReplaceAll(dstPath, "\\", "/")

		switch q.IsDIR {
		case true:
			// log.Println("mkdir", dstPath)
			err = sc.MkdirAll(dstPath)
			if err != nil {
				panic(err)
			}
		case false:
			err = sftpUploadFile(sc, srcPath, dstPath)
			if err != nil {
				log.Fatalf("could not upload file: %v", err)
			}
		}
	}
}

func ftpMain(addr, id, pwd string) {
	config := goftp.Config{
		User:            id,
		Password:        pwd,
		ActiveTransfers: true,
	}

	fc, err := goftp.DialConfig(config, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer fc.Close()
	// log.Println("Connected.")

	srcBase = replacerSlash.Replace(srcBase)
	srcRoot := filepath.Base(srcBase)
	srcCutPath := replacerSlash.Replace(strings.TrimSuffix(srcBase, srcRoot))

	err = filepath.Walk(srcBase, walkDIR)
	if err != nil {
		log.Println(err)
	}

	for _, q := range que {
		srcPath := filepath.Join("", q.Name)
		dstPath := filepath.Join(dstBase, strings.TrimPrefix(q.Name, srcCutPath))
		dstPath = strings.ReplaceAll(dstPath, "\\", "/")

		switch q.IsDIR {
		case true:
			_, err = fc.Mkdir(dstPath)
			if err != nil {
				if err.Error() == "unexpected response: 550-Directory already exists" {
					continue
				}
				if err.Error() == "failed parsing directory name: Directory created successfully" {
					continue
				}

				log.Println("DIR: ", err)
			}
		case false:
			err = ftpUploadFile(fc, srcPath, dstPath)
			if err != nil {
				log.Fatalf("could not upload file: %v", err)
			}
		}
	}

}

func main() {
	que = []QueSheet{}

	addr := "sftp.example.com:22"
	id := "username"
	pwd := "password"

	// addr := "ftp.example.com:21"
	// id := "username"
	// pwd := "password"

	// Default: sftp
	if !strings.Contains(addr, ":") {
		addr = addr + ":22"
	}

	port := strings.Split(addr, ":")[1]

	switch port {
	case "22":
		srcBase = "hihi/haha/samples"
		// srcBase = "samples"
		dstBase = "/home/sites"
		sftpMain(addr, id, pwd)
	case "21":
		srcBase = "hihi/haha/samples"
		// srcBase = "samples"
		dstBase = "/home/sites"
		ftpMain(addr, id, pwd)
	}

	que = []QueSheet{}
}
