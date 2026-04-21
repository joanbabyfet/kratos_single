package service

import (
	"context"

	v1 "kratos_single/api/client/v1"
	"kratos_single/internal/biz"

	"github.com/jinzhu/copier"
)

//一套 proto = 一个 service 实现
type AdService struct {
	v1.UnimplementedAdServer
	uc *biz.AdUsecase
}

func NewAdService(uc *biz.AdUsecase) *AdService {
	return &AdService{uc: uc}
}

// 获取首页广告列表(前几条)
func (s *AdService) GetHomeAdList(ctx context.Context, req *v1.GetHomeAdListReq) (*v1.GetHomeAdListReply, error) {

	//把 proto 请求参数转成 biz 结构
	list, _, err := s.uc.List(ctx, &biz.AdQuery{
		Limit:		int(req.Limit),
	})
	
	if err != nil {
		//不要自己包装
		return nil, err
	}

	items := make([]*v1.AdItem, 0, len(list))

	for i := range list {
		var item v1.AdItem

		if err := copier.Copy(&item, list[i]); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return &v1.GetHomeAdListReply{
		List:  items,
	}, nil
}

// 获取详情
func (s *AdService) GetAd(ctx context.Context, req *v1.GetAdReq) (*v1.GetAdReply, error) {
	
	a, err := s.uc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var res v1.GetAdReply

	// 自动拷贝同名字段
	if err := copier.Copy(&res, a); err != nil {
		return nil, err
	}

	//手动处理类型
	res.Id = int64(a.Id)

	return &res, nil
}