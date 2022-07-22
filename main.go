package main // import "hello-ftp"

import (
	"hello-ftp/uploader/config"
	"hello-ftp/uploader/ftp"
	"hello-ftp/uploader/sftp"
)

func main() {
	host := config.Host{
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
