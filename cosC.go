/*
@author '彼时思默'
@time 2020/4/27 下午5:59
@describe:
*/
package store

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type CosConnect struct {
	Client     *cos.Client
	BucketName string
	Bucket     *cos.BucketService
	Service    *cos.ServiceService
	Object     *cos.ObjectService
}

func NewCosConnector(option ...AccessOption) *CosConnect {
	var (
		accessId,
		accessSecret,
		endpoint,
		region,
		bucketName string
	)
	if option != nil {
		accessId = option[0].Id
		accessSecret = option[0].Secret
		endpoint = option[0].Endpoint
		region = option[0].Region
		bucketName = option[0].Bucket
	} else {
		accessId = os.Getenv("CosId")
		accessSecret = os.Getenv("CosSecret")
		endpoint = os.Getenv("CosEndpoint")
		bucketName = os.Getenv("CosBucket")
	}
	_ = region
	u, _ := url.Parse(endpoint)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  accessId,
			SecretKey: accessSecret,
		},
	})

	return &CosConnect{
		Client:     c,
		BucketName: bucketName,
		Bucket:     c.Bucket,
		Service:    c.Service,
		Object:     c.Object,
	}
}

func (c CosConnect) ListBuckets() {
	s, _, err := c.Service.Get(context.Background())
	if err != nil {
		panic(err)
	}

	for _, b := range s.Buckets {
		fmt.Printf("%#v\n", b)
	}
}

func (c CosConnect) ListFiles() {
	opt := &cos.BucketGetOptions{
		Prefix:  c.BucketName,
		MaxKeys: 3,
	}
	s, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		panic(err)
	}

	for _, b := range s.Contents {
		fmt.Printf("%#v\n", b)
	}
}

func (c CosConnect) UploadString(msg string, descPath string) {
	fs := strings.NewReader(msg)
	_, err := c.Object.Put(context.Background(), path.Join(c.BucketName, descPath), fs, nil)
	if err != nil {
		panic(err)
	}
}

func (c CosConnect) UploadFileByPath(sourcePath string, descPath string) {
	_, err := c.Object.PutFromFile(context.Background(), path.Join(c.BucketName, descPath), sourcePath, nil)
	if err != nil {
		panic(err)
	}
}

func (c CosConnect) UploadFileByFP(fp *os.File, descPath string) {
	_, err := c.Object.Put(context.Background(), path.Join(c.BucketName, descPath), fp, nil)
	if err != nil {
		panic(err)
	}
}
