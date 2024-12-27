package services

import (
	"context"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/errors"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/models"
	"role_ai/repos"
	"strings"
)

type VoiceService struct {
	voiceRepo repos.IVoiceRepository `inject:""`
}

func (srv *VoiceService) Provide(_ context.Context) any {
	return srv
}

func (srv *RoleService) GetVoiceList(para dto.VoiceListReq) (*dto.VoiceListResp, error) {
	var (
		list  []models.Voice
		total int64
		//listWhere where.Wheres
		sort string
		res  dto.VoiceListResp
	)

	//if para.Uid != 0 {
	//	listWhere = listWhere.And(where.Eq("uid", para.Uid))
	//}
	//if para.Name != "" {
	//	listWhere = listWhere.And(where.Like("name", "%"+para.Name+"%"))
	//}
	//
	//if para.Sort != 0 {
	//	sort += "chat_num desc,"
	//}
	sort += "id desc"

	err := srv.roleRepo.Find(&finder.Finder{
		Model: models.Voice{},
		//Wheres:         listWhere,
		OrderBy:        sort,
		Num:            para.PageNum,
		Size:           para.PageSize,
		Recipient:      &list,
		Total:          &total,
		IgnoreNotFound: true,
	})
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}

	err = models.Copy(&res.List, &list)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	for k, v := range res.List {
		if v.Desc != "" {
			res.List[k].DescArr = strings.Split(v.Desc, ";")
		}
	}
	res.Total = total
	return &res, nil
}
