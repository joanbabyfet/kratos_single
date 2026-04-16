package biz

import (
	"context"
	"kratos_single/internal/pkg/auth"
	"kratos_single/internal/pkg/i18n"

	"github.com/go-kratos/kratos/v2/errors"
)

//业务需要的数据
type Article struct {
	Id      	int
	Catid   	int
	Title   	string
	Info    	string
	Content 	string
	Img     	string
	Author  	string
	Extra   	string
	Sort    	int16
	Status		int8
	CreateUser 	string
	UpdateUser 	string
}

type ArticleQuery struct {
	Title    string
	Catid    int
	Page     int
	PageSize int
	Limit int
}

type ArticleRepo interface {
	Create(context.Context, *Article) (int64, error)
	Update(context.Context, *Article) error
	GetById(context.Context, int64, bool) (*Article, error)
	List(context.Context, *ArticleQuery) ([]*Article, int64, error)
	Delete(context.Context, int64, string) error
}

type ArticleUsecase struct {
	repo ArticleRepo
}

func NewArticleUsecase(repo ArticleRepo) *ArticleUsecase {
	return &ArticleUsecase{repo: repo}
}

// 添加
func (uc *ArticleUsecase) Create(ctx context.Context, a *Article) (int64, error) {
	if a.Title == "" {
		return 0, errors.New(400, "-1", i18n.T(auth.GetLang(ctx), "InvalidTitle"))
	}
	return uc.repo.Create(ctx, a)
}

// 修改
func (uc *ArticleUsecase) Update(ctx context.Context, a *Article) error {
	if a.Id == 0 {
		return ErrInvalidID
	}
	return uc.repo.Update(ctx, a)
}

// 获取详情(走缓存)
func (uc *ArticleUsecase) GetById(ctx context.Context, id int64) (*Article, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	return uc.repo.GetById(ctx, id, true)
}

// 获取详情(后台不走缓存)
func (uc *ArticleUsecase) AdminGetById(ctx context.Context, id int64) (*Article, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	return uc.repo.GetById(ctx, id, false)
}

// 获取列表
func (uc *ArticleUsecase) List(ctx context.Context, q *ArticleQuery) ([]*Article, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return uc.repo.List(ctx, q)
}

// 软删除
func (uc *ArticleUsecase) Delete(ctx context.Context, id int64, user string) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return uc.repo.Delete(ctx, id, user)
}