// internal/biz/upload.go
package biz

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	v1 "kratos_single/api/common/v1"
	"kratos_single/internal/pkg/utils"

	"github.com/disintegration/imaging"
	"github.com/go-kratos/kratos/v2/log"
)

type CommonUsecase struct {
	log *log.Helper
}

func NewCommonUsecase(logger log.Logger) *CommonUsecase {
	return &CommonUsecase{
		log: log.NewHelper(logger),
	}
}

func (uc *CommonUsecase) Upload(ctx context.Context, req *v1.UploadReq, uploadDir, fileURL string, maxSize int64) (*v1.UploadReply, error) {

	if int64(len(req.File)) > maxSize*1024*1024 {
		return nil, fmt.Errorf("文件过大")
	}

	ext := strings.ToLower(path.Ext(req.Filename))
	allowExt := map[string]bool{
		".jpg": true,
		".jpeg": true,
		".png": true,
		".gif": true,
	}

	if !allowExt[ext] {
		return nil, fmt.Errorf("文件格式不支持")
	}

	dirDate := time.Now().Format("20060102")
	fullDir := filepath.Join(uploadDir, req.Dir, dirDate)

	if err := os.MkdirAll(fullDir, 0775); err != nil {
		return nil, err
	}

	fileName := utils.GenFileName(ext)
	savePath := filepath.Join(fullDir, fileName)

	if err := os.WriteFile(savePath, req.File, 0644); err != nil {
		return nil, err
	}

	fileLink := fmt.Sprintf("%s/%s/%s/%s", fileURL, req.Dir, dirDate, fileName)

	if req.ThumbW > 0 || req.ThumbH > 0 {
		img, err := imaging.Open(savePath)
		if err != nil {
			return nil, err
		}

		thumb := imaging.Resize(img, int(req.ThumbW), int(req.ThumbH), imaging.Lanczos)

		fileName = utils.GenFileName(ext)
		savePath = filepath.Join(fullDir, fileName)
		fileLink = fmt.Sprintf("%s/%s/%s/%s", fileURL, req.Dir, dirDate, fileName)

		if err := imaging.Save(thumb, savePath); err != nil {
			return nil, err
		}
	}

	return &v1.UploadReply{
		RealName: req.Filename,
		FileName: dirDate + "/" + fileName,
		FileLink: fileLink,
	}, nil
}

