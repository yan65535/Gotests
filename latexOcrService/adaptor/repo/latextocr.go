package repo

import "context"

type ILatex interface {
	// 添加文件上传记录
	CreateAdminUploadRecord(ctx context.Context, req *do.AdminUploadRecord) error
}
