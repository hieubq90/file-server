package mio

import (
	"errors"
	"file-server/conf"
	"log"
	"mime/multipart"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var lock = &sync.Mutex{}
var client *minio.Client

func GetClient() *minio.Client {
	if client == nil {
		lock.Lock()
		defer lock.Unlock()
		if client == nil {
			var err error
			minioConfig := conf.GetMinioConfig()
			client, err = minio.New(minioConfig["endpoint"], &minio.Options{
				Creds:  credentials.NewStaticV4(minioConfig["access"], minioConfig["secret"], ""),
				Secure: false,
			})
			if err != nil {
				log.Fatalln(err)
				return nil
			}
		}
	}
	return client
}

func UploadFile(file *multipart.FileHeader, fileType string) (err error) {
	client = GetClient()
	if client == nil {
		err = errors.New("Không tạo được MinIO Client")
	}
	return nil
}
