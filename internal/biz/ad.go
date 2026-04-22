package biz

import (
	"context"
	"kratos_single/internal/pkg/auth"
	"kratos_single/internal/pkg/i18n"

	"github.com/go-kratos/kratos/v2/errors"
)

//业务需要的数据
type Ad struct {
	Id         int
	Catid      int
	Title      string
	Img        string
	Url    	   string
	Sort       int16
	Status     int8
	CreateUser string
	UpdateUser string
}

type AdQuery struct {
	Title    string
	Catid    int
	Page     int
	PageSize int
	Limit int
}

type AdRepo interface {
	Create(context.Context, *Ad) (int64, error)
	Update(context.Context, *Ad, string) error
	GetById(context.Context, int64, bool) (*Ad, error)
	List(context.Context, *AdQuery) ([]*Ad, int64, error)
	Delete(context.Context, int64, string) error
}

type AdUsecase struct {
	repo AdRepo
}

func NewAdUsecase(repo AdRepo) *AdUsecase {
	return &AdUsecase{repo: repo}
}

// 添加
func (uc *AdUsecase) Create(ctx context.Context, a *Ad) (int64, error) {
	if a.Title == "" {
		return 0, errors.New(400, "-1", i18n.T(auth.GetLang(ctx), "InvalidTitle"))
	}
	return uc.repo.Create(ctx, a)
}

// 修改（可选：普通更新 / FOR UPDATE / SHARE MODE）
func (uc *AdUsecase) Update(ctx context.Context, a *Ad, lockMode string) error {
	if a.Id == 0 {
		return ErrInvalidID
	}
	return uc.repo.Update(ctx, a, "")
}

// 获取详情(走缓存)
func (uc *AdUsecase) GetById(ctx context.Context, id int64) (*Ad, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	return uc.repo.GetById(ctx, id, true)
}

// 获取详情(后台不走缓存)
func (uc *AdUsecase) AdminGetById(ctx context.Context, id int64) (*Ad, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	return uc.repo.GetById(ctx, id, false)
}

// 获取列表
func (uc *AdUsecase) List(ctx context.Context, q *AdQuery) ([]*Ad, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return uc.repo.List(ctx, q)
}

// 软删除
func (uc *AdUsecase) Delete(ctx context.Context, id int64, user string) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return uc.repo.Delete(ctx, id, user)
}