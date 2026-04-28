package biz

import "context"

// MailRepo = 业务层接口（依赖倒置）
type MailRepo interface {
	SendSMTP(ctx context.Context, to string, subject string, body string) error
	Send(ctx context.Context, to string, subject string, body string) error
}

// MailUsecase
type MailUsecase struct {
	repo MailRepo
}

// Wire 注入
func NewMailUsecase(repo MailRepo) *MailUsecase {
	return &MailUsecase{repo: repo}
}

// 发欢迎信(暂未使用)
func (uc *MailUsecase) SendWelcomeMail(
	ctx context.Context,
	email string,
) error {

	subject := "Welcome"
	body := "Thanks for register."

	return uc.repo.SendSMTP(
		ctx,
		email,
		subject,
		body,
	)
}