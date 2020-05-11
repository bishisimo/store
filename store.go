/*
@author '彼时思默'
@time 2020/4/27 下午6:00
@describe:
*/
package store

import "os"

type StoreFace interface {
	//列出所有储存桶
	ListBuckets()
	//列出所有文件
	ListFiles()
	//字符串作为文件上传
	UploadString(msg string, descPath string)
	//从路径获取数据上传
	UploadFileByPath(sourcePath string, descPath string)
	//从文件句柄获取数据上传
	UploadFileByFP(fp *os.File, descPath string)
}
type AccessOption struct {
	Id       string `json:"accessId"`
	Secret   string `json:"accessSecret"`
	Endpoint string `json:"accessEndpoint"`
	Region   string `json:"accessRegion"`
	Bucket   string `json:"accessBucket"`
}
