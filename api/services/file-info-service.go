package services

import "os"

type FileInfoService interface {
	GetFileInfo(recordId, pageId string, format string) (os.FileInfo, error)
}
