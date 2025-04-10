package repo

import (
	"context"
	"latexOcrService/service/do"
)

type ILatex interface {
	// 添加文件上传记录
	CreateAdminUploadRecord(ctx context.Context, req *do.AdminUploadRecord) error
}
