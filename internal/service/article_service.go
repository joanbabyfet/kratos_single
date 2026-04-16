package service

import (
	"context"

	v1 "kratos_single/api/client/article/v1"
	"kratos_single/internal/biz"

	"github.com/jinzhu/copier"
)

//一套 proto = 一个 service 实现
type ArticleService struct {
	v1.UnimplementedArticleServer
	uc *biz.ArticleUsecase
}

func NewArticleService(uc *biz.ArticleUsecase) *ArticleService {
	return &ArticleService{uc: uc}
}

// 获取首页文章列表(前几条)
func (s *ArticleService) GetHomeArticleList(ctx context.Context, req *v1.GetHomeArticleListReq) (*v1.GetHomeArticleListReply, error) {

	//把 proto 请求参数转成 biz 结构
	list, _, err := s.uc.List(ctx, &biz.ArticleQuery{
		Limit:		int(req.Limit),
	})
	
	if err != nil {
		//不要自己包装
		return nil, err
	}

	items := make([]*v1.ArticleItem, 0, len(list))

	for i := range list {
		var item v1.ArticleItem

		if err := copier.Copy(&item, list[i]); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return &v1.GetHomeArticleListReply{
		List:  items,
	}, nil
}

// 获取列表
func (s *ArticleService) GetArticleList(ctx context.Context, req *v1.GetArticleListReq) (*v1.GetArticleListReply, error) {

	//把 proto 请求参数转成 biz 结构
	list, total, err := s.uc.List(ctx, &biz.ArticleQuery{
		Title:		req.Title,
		Catid:		int(req.Catid),
		Page:     	int(req.Page),
		PageSize: 	int(req.PageSize),
	})
	
	if err != nil {
		//不要自己包装
		return nil, err
	}

	items := make([]*v1.ArticleItem, 0, len(list))

	for i := range list {
		var item v1.ArticleItem

		if err := copier.Copy(&item, list[i]); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return &v1.GetArticleListReply{
		List:  items,
		Total: total,
	}, nil
}

// 获取详情
func (s *ArticleService) GetArticle(ctx context.Context, req *v1.GetArticleReq) (*v1.GetArticleReply, error) {
	
	a, err := s.uc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var res v1.GetArticleReply

	// 自动拷贝同名字段
	if err := copier.Copy(&res, a); err != nil {
		return nil, err
	}

	//手动处理类型
	res.Id = int64(a.Id)

	return &res, nil
}