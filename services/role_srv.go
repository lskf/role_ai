package services

import (
	"context"
	"encoding/json"
	"github.com/leor-w/kid/database/repos/creator"
	"github.com/leor-w/kid/database/repos/deleter"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/updater"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/models"
	"role_ai/repos"
	"strings"
	"time"
)

type RoleService struct {
	roleRepo repos.IRoleRepository `inject:""`
}

func (srv *RoleService) Provide(_ context.Context) any {
	return srv
}

func (srv *RoleService) GetRoleList(para dto.RoleListReq) (*dto.RoleListResp, error) {
	var (
		roleList  []models.Role
		total     int64
		listWhere where.Wheres
		sort      string
		res       dto.RoleListResp
	)

	if para.Uid != 0 {
		listWhere = listWhere.And(where.Eq("uid", para.Uid))
	}
	if para.Name != "" {
		listWhere = listWhere.And(where.Like("name", "%"+para.Name+"%"))
	}

	if para.Sort != 0 {
		sort += "chat_num desc,"
	}
	sort += "id desc"

	err := srv.roleRepo.Find(&finder.Finder{
		Model:          models.Role{},
		Wheres:         listWhere,
		OrderBy:        sort,
		Num:            para.PageNum,
		Size:           para.PageSize,
		Recipient:      &roleList,
		Total:          &total,
		IgnoreNotFound: true,
	})
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}

	err = models.Copy(&res.List, &roleList)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	res.Total = total
	return &res, nil
}

func (srv *RoleService) GetDetailById(id int64) (*dto.Role, error) {
	//获取角色详情
	role := models.Role{}
	err := srv.roleRepo.GetOne(&finder.Finder{
		Model:     new(models.Role),
		Wheres:    where.New().And(where.Eq("id", id)),
		Recipient: &role,
	})
	if err != nil {
		return nil, errors.New(ecode.RoleNotExistErr, err)
	}
	data := dto.Role{}
	err = models.Copy(&data, &role)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	//获取角色风格
	roleStyle := models.RoleStyle{}
	err = srv.roleRepo.GetOne(&finder.Finder{
		Model:     new(models.RoleStyle),
		Wheres:    where.New().And(where.Eq("role_id", id)),
		Recipient: &roleStyle,
	})
	if err != nil {
		return nil, errors.New(ecode.RoleNotExistErr, err)
	}
	speechStyleList := make([]dto.SpeechStyleObj, 0)
	err = json.Unmarshal([]byte(roleStyle.Content), &speechStyleList)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	data.StyleArray = speechStyleList
	//标签
	tagArr := strings.Split(role.Tag, ";")
	data.TagArray = tagArr
	//游戏化
	gamificationObj := dto.GamificationObj{}
	err = json.Unmarshal([]byte(role.Gamification), &gamificationObj)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	data.GamificationObj = gamificationObj
	return &data, nil
}

func (srv *RoleService) CreateRole(user *models.User, data dto.CreateRoleReq) (int64, error) {
	if data.IsPublic == models.PrivateRole && user.Menber == models.UserMenberNormal {
		return 0, errors.New(ecode.MenberPermissionErr)
	}

	//判断声音是否存在
	if data.VoiceId > 0 {
		voice := models.Voice{}
		err := srv.roleRepo.GetOne(&finder.Finder{
			Model:     new(models.Voice),
			Wheres:    where.New().And(where.Eq("id", data.VoiceId)),
			Recipient: &voice,
		})
		if err != nil {
			return 0, errors.New(ecode.VoiceNotFound, err)
		}
	}
	role := models.Role{}
	err := models.Copy(&role, &data.Role)
	if err != nil {
		return 0, errors.New(ecode.DataProcessingErr, err)
	}
	role.Uid = user.Id
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	//构建roleStyle
	roleStyle := models.RoleStyle{}
	styleStr, err := json.Marshal(data.StyleArray)
	if err != nil {
		return 0, errors.New(ecode.DataProcessingErr, err)
	}
	roleStyle.Content = string(styleStr)

	//开启事务
	err = srv.roleRepo.Transaction(func(ctx context.Context) error {
		//添加角色
		err = srv.roleRepo.Create(&creator.Creator{
			Tx:   ctx,
			Data: &role,
		})
		if err != nil {
			return errors.New(ecode.DatabaseErr, err)
		}
		//添加角色风格
		roleStyle.RoleId = role.Id
		err = srv.roleRepo.Create(&creator.Creator{
			Tx:   ctx,
			Data: &roleStyle,
		})
		if err != nil {
			return errors.New(ecode.DatabaseErr, err)
		}
		return nil
	})
	if err != nil {
		return 0, errors.New(ecode.DatabaseErr, err)
	}
	return role.Id, nil
}

func (srv *RoleService) UpdateRole(user models.User, data dto.UpdateRoleResp) error {
	roleDetail := models.Role{}
	err := srv.roleRepo.GetOne(&finder.Finder{
		Model:     new(models.Role),
		Wheres:    where.New().And(where.Eq("id", data.Id), where.Eq("uid", user.Uid)),
		Recipient: &roleDetail,
	})
	if err != nil {
		return errors.New(ecode.RoleNotExistErr, err)
	}

	//普通用户不能设置私密角色
	if data.IsPublic == models.PrivateRole && user.Menber == models.UserMenberNormal {
		return errors.New(ecode.MenberPermissionErr)
	}
	//公开角色不能修改为私密
	if roleDetail.IsPublic == 1 && data.IsPublic == 2 {
		return errors.New(ecode.PublicChangeErr)
	}

	err = models.Copy(&roleDetail, &data)
	if err != nil {
		return errors.New(ecode.DataProcessingErr, err)
	}
	roleDetail.UpdatedAt = time.Now()

	//构建roleStyle
	roleStyle := models.RoleStyle{}
	roleStyle.RoleId = roleDetail.Id
	styleStr, err := json.Marshal(data.StyleArray)
	if err != nil {
		return errors.New(ecode.DataProcessingErr, err)
	}
	roleStyle.Content = string(styleStr)

	//开启事务
	err = srv.roleRepo.Transaction(func(ctx context.Context) error {
		//添加角色
		err = srv.roleRepo.Update(&updater.Updater{
			Tx:     ctx,
			Model:  new(models.Role),
			Wheres: where.New().And(where.Eq("id", roleDetail.Id)),
			Fields: map[string]interface{}{
				"avatar":       roleDetail.Avatar,
				"role_name":    roleDetail.RoleName,
				"desc":         roleDetail.Desc,
				"remark":       roleDetail.Remark,
				"tag":          roleDetail.Tag,
				"gamification": roleDetail.Gamification,
				"is_public":    roleDetail.IsPublic,
				"voice_id":     roleDetail.VoiceId,
				"updated_at":   roleDetail.UpdatedAt,
			},
		})
		if err != nil {
			return errors.New(ecode.DatabaseErr, err)
		}
		//删除原角色风格
		err = srv.roleRepo.Delete(&deleter.Deleter{
			Tx:     ctx,
			Model:  new(models.RoleStyle),
			Wheres: where.New().And(where.Eq("role_id", roleDetail.Id)),
		})
		if err != nil {
			return errors.New(ecode.DatabaseErr, err)
		}
		//添加角色风格
		err = srv.roleRepo.Create(&creator.Creator{
			Tx:   ctx,
			Data: &roleStyle,
		})
		if err != nil {
			return errors.New(ecode.DatabaseErr, err)
		}
		return nil
	})
	if err != nil {
		return errors.New(ecode.DatabaseErr, err)
	}

	return nil
}

func (srv *RoleService) DeleteRole() {}
