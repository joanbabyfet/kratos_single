package service

import (
	"context"

	adminv1 "kratos_single/api/admin/v1"
	"kratos_single/internal/biz"
	"kratos_single/internal/pkg/auth"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
)

//一套 proto = 一个 service 实现
type AdminArticleService struct {
	adminv1.UnimplementedAdminArticleServer
	uc *biz.ArticleUsecase
	log *log.Helper //不要全局用 log.xxx，应该依赖注入 logger
}

func NewAdminArticleService(uc *biz.ArticleUsecase, logger log.Logger) *AdminArticleService {
	return &AdminArticleService{
		uc: uc,
		log: log.NewHelper(log.With(logger, "module", "admin-article-service")),
	}
}

// 获取列表
func (s *AdminArticleService) GetArticleList(ctx context.Context, req *adminv1.GetArticleListReq) (*adminv1.GetArticleListReply, error) {

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

	items := make([]*adminv1.ArticleItem, 0, len(list))

	for i := range list {
		var item adminv1.ArticleItem

		if err := copier.Copy(&item, list[i]); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return &adminv1.GetArticleListReply{
		List:  items,
		Total: total,
	}, nil
}

// 获取详情(不走缓存)
func (s *AdminArticleService) GetArticle(ctx context.Context, req *adminv1.GetArticleReq) (*adminv1.GetArticleReply, error) {

	a, err := s.uc.AdminGetById(ctx, req.Id)
	if err != nil {
		//写入日志
		s.log.Errorw(
            "msg", "文章不存在",
            "id", req.Id,
            "err", err,
        )
		return nil, err
	}

	var res adminv1.GetArticleReply

	// 自动拷贝同名字段
	if err := copier.Copy(&res, a); err != nil {
		return nil, err
	}

	//手动处理类型
	res.Id = int64(a.Id)

	return &res, nil
}

// 添加
func (s *AdminArticleService) CreateArticle(ctx context.Context, req *adminv1.CreateArticleReq) (*adminv1.CreateArticleReply, error) {

	user := auth.GetUser(ctx) //从 ctx 里拿用户
	var a biz.Article

	// 自动拷贝
	if err := copier.Copy(&a, req); err != nil {
		return nil, err
	}
	a.CreateUser = user

	id, err := s.uc.Create(ctx, &a)
	if err != nil {
		//写入日志
		s.log.Errorw(
            "msg", "文章添加失败",
            "id", id,
            "err", err,
        )
		return nil, err
	}

	return &adminv1.CreateArticleReply{
		Id: id,
	}, nil
}

// 修改
func (s *AdminArticleService) UpdateArticle(ctx context.Context, req *adminv1.UpdateArticleReq) (*adminv1.UpdateArticleReply, error) {

	user := auth.GetUser(ctx) //从 ctx 里拿用户
	var a biz.Article

	// 自动拷贝
	if err := copier.Copy(&a, req); err != nil {
		return nil, err
	}
	a.UpdateUser = user 

	err := s.uc.Update(ctx, &a)
	if err != nil {
		//写入日志
		s.log.Errorw(
            "msg", "文章更新失败",
            "id", req.Id,
            "err", err,
        )
		return nil, err
	}

	return &adminv1.UpdateArticleReply{}, nil
}

// 删除
func (s *AdminArticleService) DeleteArticle(ctx context.Context, req *adminv1.DeleteArticleReq) (*adminv1.DeleteArticleReply, error) {
	user := auth.GetUser(ctx) //从 ctx 里拿用户

	err := s.uc.Delete(ctx, req.Id, user)
	if err != nil {
		//写入日志
		s.log.Errorw(
            "msg", "文章删除失败",
            "id", req.Id,
            "err", err,
        )
		return nil, err
	}

	return &adminv1.DeleteArticleReply{}, nil
}