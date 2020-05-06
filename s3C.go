/*
@author '彼时思默'
@time 2020/4/8 15:28
@describe:
*/
package store

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

type S3Connector struct {
	Sess       *session.Session
	Svc        *s3.S3
	Endpoint   string
	Region     string
	BucketName string
}

func NewS3Connector() *S3Connector {
	accessId := os.Getenv("S3Id")
	accessSecret := os.Getenv("S3Secret")
	endpoint := os.Getenv("S3Endpoint")
	region := os.Getenv("S3Region")
	bucketName := os.Getenv("S3Bucket")
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessId, accessSecret, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
	})
	if err != nil {
		fmt.Println("aws session error!")
	}
	svc := s3.New(sess)
	return &S3Connector{
		Sess:       sess,
		Svc:        svc,
		BucketName: bucketName,
	}
}
func (s S3Connector) ListBuckets() {
	result, err := s.Svc.ListBuckets(nil)
	if err != nil {
		fmt.Printf("Unable to list buckets, %v\n", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func (s S3Connector) ListFiles() {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(s.BucketName),
	}
	resp, err := s.Svc.ListObjects(params)
	if err != nil {
		fmt.Printf("Unable to list buckets, %v\n", err)
	}
	for _, item := range resp.Contents {
		_ = item
		fmt.Println(*item.Key, *item.LastModified, *item.Size, *item.StorageClass)
	}
}

/*
上传一个文件
*/
func (s S3Connector) UploadFil2eByPath(filePath string, descPath string) {
	fp, err := os.Open(filePath)
	if err != nil {
		fmt.Println(filePath, "打开错误! ", err.Error())
	} else if fp != nil {
		defer fp.Close()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()
	res, _ := s.Svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(descPath),
		Body:   fp,
	})
	fmt.Println(res)
}

func (s S3Connector) UploadFileByFP(fp *os.File, descPath string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()
	res, err := s.Svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(descPath),
		Body:   fp,
	})
	if err != nil {
		logrus.Errorf("上传文件失败:%s", err)
	} else {
		logrus.Infof("上传s3文件成功:%s", *res.ETag)
	}
}

/*
上传字符串保存到文件
*/
func (s S3Connector) UploadString(msg string, desPath string) {
	fp := strings.NewReader(msg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()
	res, err := s.Svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(desPath),
		Body:   fp,
	})
	if err != nil {
		logrus.Errorf("上传文件失败:%s", err)
	} else {
		logrus.Infof("上传s3文件成功:%s", *res.ETag)
	}
}

func (s S3Connector) Download(filePath string) {
	filePath = path.Join("s3Download", filePath)
	dirSp := strings.Split(filePath, "/")
	dir := path.Join(dirSp[:len(dirSp)-1]...)
	if _, err := os.Stat(dir); !os.IsExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Unable to open file %q, %v", err)
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(s.Sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(filePath),
		})
	if err != nil {
		fmt.Println("Unable to download item %q, %v", filePath, err)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func (s S3Connector) DownloadAll() []string {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(s.BucketName),
	}
	resp, err := s.Svc.ListObjects(params)
	if err != nil {
		logrus.Panic("Unable to list bucket files:", err)
	}
	filePath := make([]string, 0)
	for _, item := range resp.Contents {
		s3Path := *item.Key
		if s3Path != "" && s3Path[len(s3Path)-1] != '/' {
			fmt.Println(s3Path, *item.LastModified, *item.Size, *item.StorageClass)
			s.Download(s3Path)
			filePath = append(filePath, s3Path)
		}
	}
	return filePath
}
