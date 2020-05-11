/*
@author '彼时思默'
@time 2020/4/27 下午6:00
@describe:
*/
package store

import "os"

type StoreFace interface {
	ListBuckets()
	ListFiles()
	UploadString(msg string, descPath string)
	UploadFil2eByPath(sourcePath string, descPath string)
	UploadFileByFP(fp *os.File, descPath string)
}
