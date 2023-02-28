package mio

import (
	"context"
	"errors"
	"file-server/conf"
	"file-server/server/presenter"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var lock = &sync.Mutex{}
var client *minio.Client

func GetClient() (*minio.Client, error) {
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
				return nil, err
			}
		}
	}
	return client, nil
}

func CreateBucket(name string) error {
	var err error
	aClient, err := GetClient()
	if err != nil {
		return err
	}
	if aClient != nil {
		isExist, err := aClient.BucketExists(context.Background(), name)
		if err != nil {
			return err
		}
		if !isExist {
			err = aClient.MakeBucket(context.Background(), name, minio.MakeBucketOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UploadFile(project string, folder string, file *multipart.FileHeader, fileType string) (*presenter.FileObject, error) {

	aClient, err := GetClient()
	if aClient == nil || err != nil {
		return nil, errors.New("Không tạo được MinIO Client")
	}

	id := (uuid.New()).String()
	filename := fmt.Sprintf("%s/%s/%s_%s", project, folder, id, file.Filename)
	openFile, err := file.Open()
	defer openFile.Close()
	if err != nil {
		return nil, err
	}

	bucket := conf.GetMinioConfig()["bucket"]
	_, err = aClient.PutObject(context.Background(), bucket, filename, openFile, file.Size, minio.PutObjectOptions{ContentType: fileType})
	if err != nil {
		return nil, err
	}
	return &presenter.FileObject{
		ID:   id,
		Name: filename,
	}, nil
}

func GenerateDownloadLink(project string, folder string, name string) (*url.URL, error) {
	filename := fmt.Sprintf("%s/%s/%s", project, folder, name)
	aClient, err := GetClient()
	if aClient == nil || err != nil {
		return nil, errors.New("Không tạo được MinIO Client")
	}

	bucket := conf.GetMinioConfig()["bucket"]
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
	link, err := aClient.PresignedGetObject(context.Background(), bucket, filename, 6*time.Hour, reqParams)
	if link == nil || err != nil {
		return nil, err
	}
	return link, nil
}
