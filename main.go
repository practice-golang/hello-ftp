package main // import "hello-ftp"

import (
	"hello-ftp/ftp"
	"hello-ftp/model"
	"hello-ftp/sftp"
)

func main() {
	host := model.Host{
		// Hostname: "ftp.example.com",
		// Port:     "21",
		Hostname: "sftp.example.com",
		Port:     "22",
		Username: "username",
		Password: "password",
		SrcBase:  "hihi/haha/samples",
		DstBase:  "/home/sites",
	}

	switch host.Port {
	case "21":
		ftp.ProcMain(host)
	case "22":
		sftp.ProcMain(host)
	}
}
