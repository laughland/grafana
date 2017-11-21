package imguploader

import (
	"context"
	"fmt"
	"log"
    "net/url"
    "time"

	"github.com/minio/minio-go"
	//"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/util"
)

type MinioUploader struct {
	endpoint        string
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	expiry          string
	useSSL          bool
	log             log.Logger
}

func NewMinioUploader(endpoint string, bucketName string, accessKeyID string, secretAccessKey string, useSSL bool) *MinioUploader {
	return &MinioUploader{
		endpoint:        endpoint,
		bucketName:      bucketName,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		expiry:          expiry,
		useSSL:          useSSL,
	}
}

func (u *MinioUploader) Upload(ctx context.Context, imageDiskPath string) (string, error) {

    // Initialize minio client object.
    minioClient, err := minio.New(u.endpoint, u.accessKeyID, u.secretAccessKey, u.useSSL)

    if err != nil {
       log.Fatalln(err)
    }

    //create random name for image
    objectName := util.GetRandomString(20) + ".png"
    contentType := "image/png"

    // Upload the image file with FPutObject
    n, err := minioClient.FPutObject(u.bucketName, objectName, imageDiskPath, minio.PutObjectOptions{ContentType:contentType})
    log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

    if err != nil {
      log.Fatalln(err)
    }

    // Set request parameters for content-disposition.
    reqParams := make(url.Values)
    reqParams.Set("response-content-disposition", "attachment; filename=\"objectName\"")

    // Generates a presigned url which expires per the expiry config setting
    secondsToExpire := (time.Second * 24 * 60 * 60) * (u.expiry)
    presignedURL, err := minioClient.PresignedGetObject(u.bucketName, objectName, secondsToExpire, reqParams)
    if err != nil {
        fmt.Println(err)
    }

    return presignedURL.String(), nil
}


