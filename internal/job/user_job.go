package job

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"kratos_single/internal/biz"
)

//共享同一套 biz/data
type UserJob struct {
	uc *biz.UserUsecase
}

func NewUserJob(uc *biz.UserUsecase) *UserJob {
	return &UserJob{uc: uc}
}

// 定时同步用户
func (j *UserJob) SyncUser() {

	log.Info("开始同步用户")

	err := j.uc.SyncUser(context.Background())
	if err != nil {
		log.Info("同步失败:", err)
		return
	}

	log.Info("同步成功")
}