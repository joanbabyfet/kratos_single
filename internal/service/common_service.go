// internal/service/upload.go
package service

import (
	"context"

	v1 "kratos_single/api/common/v1"
	"kratos_single/internal/biz"

	"github.com/go-kratos/kratos/v2/config"
)

type CommonService struct {
	v1.UnimplementedUploadServer
	uc  *biz.CommonUsecase
	cfg config.Config
}

func NewUploadService(uc *biz.CommonUsecase) *CommonService {
	return &CommonService{
		uc:  uc,
	}
}

//上传
func (s *CommonService) UploadFile(ctx context.Context, req *v1.UploadReq) (*v1.UploadReply, error) {

	var uploadDir string
	var fileURL string
	var maxSize int64

	// s.cfg.Value("file.upload_dir").Scan(&uploadDir)
	// s.cfg.Value("file.file_url").Scan(&fileURL)
	// s.cfg.Value("file.upload_max_size").Scan(&maxSize)
	uploadDir = "image"
	fileURL = ""
	maxSize = 5

	return s.uc.Upload(ctx, req, uploadDir, fileURL, maxSize)
}