// biz/mq.go
package biz

import (
	"context"
	"fmt"
)

// MQRepo 定义 RabbitMQ 行为接口
type MQRepo interface {
	Publish(ctx context.Context, queue string, body string) error
}

// MQUsecase 业务层
type MQUsecase struct {
	repo MQRepo
}

func NewMQUsecase(repo MQRepo) *MQUsecase {
	return &MQUsecase{repo: repo}
}

// SendWelcomeMessage 示例业务逻辑
func (uc *MQUsecase) SendWelcomeMessage(ctx context.Context, userID int64) error {
	body := "welcome user id = " + fmt.Sprint(userID)
	return uc.repo.Publish(ctx, "test_queue", body)
}