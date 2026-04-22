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
type AdminAdService struct {
	adminv1.UnimplementedAdminAdServer
	uc *biz.AdUsecase
	log *log.Helper //不要全局用 log.xxx，应该依赖注入 logger
}

func NewAdminAdService(uc *biz.AdUsecase, logger log.Logger) *AdminAdService {
	return &AdminAdService{
		uc: uc,
		log: log.NewHelper(log.With(logger, "module", "admin-Ad-service")),
	}
}

// 获取列表
func (s *AdminAdService) GetAdList(ctx context.Context, req *adminv1.GetAdListReq) (*adminv1.GetAdListReply, error) {

	//把 proto 请求参数转成 biz 结构
	list, total, err := s.uc.List(ctx, &biz.AdQuery{
		Title:		req.Title,
		Catid:		int(req.Catid),
		Page:     	int(req.Page),
		PageSize: 	int(req.PageSize),
	})
	
	if err != nil {
		//不要自己包装
		return nil, err
	}

	items := make([]*adminv1.AdItem, 0, len(list))

	for i := range list {
		var item adminv1.AdItem

		if err := copier.Copy(&item, list[i]); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return &adminv1.GetAdListReply{
		List:  items,
		Total: total,
	}, nil
}

// 获取详情(不走缓存)
func (s *AdminAdService) GetAd(ctx context.Context, req *adminv1.GetAdReq) (*adminv1.GetAdReply, error) {

	a, err := s.uc.AdminGetById(ctx, req.Id)
	if err != nil {
		//写入日志
		s.log.Errorw(
            "msg", "广告不存在",
            "id", req.Id,
            "err", err,
        )
		return nil, err
	}

	var res adminv1.GetAdReply

	// 自动拷贝同名字段
	if err := copier.Copy(&res, a); err != nil {
		return nil, err
	}

	//手动处理类型
	res.Id = int64(a.Id)

	return &res, nil
}

// 添加
func (s *AdminAdService) CreateAd(ctx context.Context, req *adminv1.CreateAdReq) (*adminv1.CreateAdReply, error) {

	user := auth.GetUser(ctx) //从 ctx 里拿用户
	var a biz.Ad

	// 自动拷贝
	if err := copier.Copy(&a, req); err != nil {
		return nil, err
	}
	a.CreateUser = user

	id, err := s.uc.Create(ctx, &a)
	if err != nil {
		//写入日志
		s.log.Errorw(
            "msg", "广告添加失败",
            "id", id,
            "err", err,
        )
		return nil, err
	}

	return &adminv1.CreateAdReply{
		Id: id,
	}, nil
}

// 修改
func (s *AdminAdService) UpdateAd(ctx context.Context, req *adminv1.UpdateAdReq) (*adminv1.UpdateAdReply, error) {

	user := auth.GetUser(ctx) //从 ctx 里拿用户
	var a biz.Ad

	// 自动拷贝
	if err := copier.Copy(&a, req); err != nil {
		return nil, err
	}
	a.UpdateUser = user 

	err := s.uc.Update(ctx, &a, "")
	if err != nil {
		//写入日志
		s.log.Errorw(
            "msg", "广告更新失败",
            "id", req.Id,
            "err", err,
        )
		return nil, err
	}

	return &adminv1.UpdateAdReply{}, nil
}

// 删除
func (s *AdminAdService) DeleteAd(ctx context.Context, req *adminv1.DeleteAdReq) (*adminv1.DeleteAdReply, error) {
	user := auth.GetUser(ctx) //从 ctx 里拿用户

	err := s.uc.Delete(ctx, req.Id, user)
	if err != nil {
		//写入日志
		s.log.Errorw(
            "msg", "广告删除失败",
            "id", req.Id,
            "err", err,
        )
		return nil, err
	}

	return &adminv1.DeleteAdReply{}, nil
}