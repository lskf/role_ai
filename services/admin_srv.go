package services

import (
	"context"
	"github.com/leor-w/kid/database/repos/creator"
	"github.com/leor-w/kid/database/repos/deleter"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/updater"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"
	"github.com/leor-w/kid/guard"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/models"
	"role_ai/repos"
)

type AdminService struct {
	repo  repos.IAdminRepository `inject:""`
	guard guard.Guard            `inject:""`
}

func (srv *AdminService) Provide(context.Context) any {
	return srv
}

func (srv *AdminService) Create(admin *models.Admin, para *dto.AdminCreateReq) (id int64, err error) {
	var (
		adminDetail models.Admin
	)
	exist, err := srv.repo.PhoneOrAccountExisted(para.Phone, para.Account)
	if err != nil {
		return 0, errors.New(ecode.DatabaseErr, err)
	}
	if exist {
		return 0, errors.New(ecode.AdminNotFound)
	}
	err = models.Copy(&adminDetail, &para)
	if err != nil {
		return 0, errors.New(ecode.DataProcessingErr, err)
	}
	adminDetail.AdminId = srv.repo.GenerateUid()

	//加密密码
	err = (&adminDetail).EncodePassword()
	if err != nil {
		return 0, errors.New(ecode.DataProcessingErr, err)
	}

	adminDetail.CreatorField = models.CreatorField{
		CreatorId:   admin.AdminId,
		CreatorName: admin.RealName,
	}

	err = srv.repo.Create(&creator.Creator{
		Data: &adminDetail,
	})
	if err != nil {
		return 0, errors.New(ecode.DatabaseErr, err)
	}
	return adminDetail.AdminId, nil
}

func (srv *AdminService) Update(admin *models.Admin, para *dto.AdminUpdateReq) (ok bool, err error) {
	var adminDetail models.Admin
	if err = srv.repo.GetOne(&finder.Finder{
		Model:          &models.Admin{},
		Wheres:         where.New().And(where.Eq("uid", para.Uid)),
		Recipient:      &adminDetail,
		IgnoreNotFound: true,
	}); err != nil {
		return false, errors.New(ecode.DatabaseErr, err)
	}
	if adminDetail.Id <= 0 {
		return false, errors.New(ecode.AdminNotFound)
	}

	err = models.Copy(&adminDetail, &para)
	if err != nil {
		return false, errors.New(ecode.DataProcessingErr, err)
	}
	adminDetail.UpdaterId = admin.AdminId
	adminDetail.UpdaterName = admin.RealName

	err = srv.repo.Update(&updater.Updater{
		Model:   &models.Admin{},
		Selects: []string{"phone", "account", "avatar", "real_name", "role", "updater_id", "updated_at", "updater_name"},
		Update:  &adminDetail,
	})
	if err != nil {
		return false, errors.New(ecode.DatabaseErr, err)
	}
	return true, nil
}

func (srv *AdminService) UpdatePassword(admin *models.Admin, para *dto.AdminUpdatePwdReq) (ok bool, err error) {
	var adminDetail models.Admin
	if err = srv.repo.GetOne(&finder.Finder{
		Model:          &models.Admin{},
		Wheres:         where.New().And(where.Eq("uid", para.Uid)),
		Recipient:      &adminDetail,
		IgnoreNotFound: true,
	}); err != nil {
		return false, errors.New(ecode.DatabaseErr, err)
	}
	if adminDetail.Id <= 0 {
		return false, errors.New(ecode.AdminNotFound)
	}
	//密码加密
	adminDetail.Password = para.NewPwd
	err = (&adminDetail).EncodePassword()
	if err != nil {
		return false, errors.New(ecode.DataProcessingErr, err)
	}
	adminDetail.UpdaterId = admin.AdminId
	adminDetail.UpdaterName = admin.RealName
	err = srv.repo.Update(&updater.Updater{
		Model:   &models.Admin{},
		Selects: []string{"password", "salt", "updater_id", "updated_at", "updater_name"},
		Update:  &adminDetail,
	})
	if err != nil {
		return false, errors.New(ecode.DatabaseErr, err)
	}
	return true, nil
}

func (srv *AdminService) Delete(admin *models.Admin, deleteId int64) (ok bool, err error) {
	adminDetail := models.Admin{}
	if err = srv.repo.GetOne(&finder.Finder{
		Model:          &models.Admin{},
		Wheres:         where.New().And(where.Eq("uid", deleteId), where.IsNull("deleted_at")),
		Recipient:      &adminDetail,
		IgnoreNotFound: true,
	}); err != nil {
		return false, errors.New(ecode.DatabaseErr, err)
	}
	if adminDetail.Id <= 0 {
		return false, errors.New(ecode.AdminNotFound)
	}

	if err = srv.repo.Transaction(func(ctx context.Context) error {
		if err = srv.repo.Update(&updater.Updater{
			Tx:     ctx,
			Model:  new(models.Admin),
			Wheres: where.New().And(where.Eq("uid", deleteId)),
			Fields: map[string]interface{}{
				"phone":        adminDetail.Phone + models.GetDeletedMark(),   //唯一索引添加删除标识
				"account":      adminDetail.Account + models.GetDeletedMark(), //唯一索引添加删除标识
				"deleter_id":   admin.AdminId,
				"deleter_name": admin.RealName,
			},
		}); err != nil {
			return err
		}
		if err = srv.repo.Delete(&deleter.Deleter{
			Tx:     ctx,
			Model:  new(models.Admin),
			Wheres: where.New().And(where.Eq("uid", deleteId)),
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false, errors.New(ecode.DatabaseErr, err)
	}
	return true, nil
}

func (srv *AdminService) GetList(para *dto.AdminListReq) (*dto.AdminListResp, error) {
	var (
		adminList []models.Admin
		total     int64
		listWhere where.Wheres
		res       dto.AdminListResp
	)
	if para.Uid != 0 {
		listWhere = listWhere.And(where.Eq("uid", para.Uid))
	}
	if para.Phone != "" {
		listWhere = listWhere.And(where.Eq("phone", para.Phone))
	}
	if para.Account != "" {
		listWhere = listWhere.And(where.Eq("account", para.Account))
	}
	if para.Role != 0 {
		listWhere = listWhere.And(where.Eq("role", para.Role))
	}
	if para.RealName != "" {
		listWhere = listWhere.And(where.Eq("real_name", para.RealName))
	}
	err := srv.repo.Find(&finder.Finder{
		Model:          models.Admin{},
		Wheres:         listWhere,
		OrderBy:        "id desc",
		Num:            para.PageNum,
		Size:           para.PageSize,
		Recipient:      &adminList,
		Total:          &total,
		IgnoreNotFound: true,
	})
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}

	err = models.Copy(&res.List, &adminList)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	res.Total = total
	return &res, nil
}

func (srv *AdminService) Detail(uid int64) (*dto.Admin, error) {
	var (
		detail models.Admin
		res    dto.Admin
	)
	err := srv.repo.GetOne(&finder.Finder{
		Model:          &models.Admin{},
		Wheres:         where.New().And(where.Eq("uid", uid)),
		Recipient:      &detail,
		IgnoreNotFound: true,
	})
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}
	if detail.Id <= 0 {
		return nil, errors.New(ecode.AdminNotFound)
	}
	err = models.Copy(&res, &detail)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	return &res, nil
}

func (srv *AdminService) GetListByUIds(uids []int64) (list []models.Admin, err error) {
	err = srv.repo.GetOne(&finder.Finder{
		Model:          &models.Admin{},
		Wheres:         where.New().And(where.In("uid", uids)),
		Recipient:      &list,
		IgnoreNotFound: true,
	})
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}
	return
}

func (srv *AdminService) GetCreatorAndUpdaterList(slice interface{}) (adminListM map[int64]models.Admin, err error) {
	//获取创建人的uid
	adminUids := make([]int64, 0)
	creatorUids := GetFieldValues(slice, "CreatorId")
	for _, v := range creatorUids {
		if uid, ok := v.(int64); ok {
			adminUids = append(adminUids, uid)
		}
	}
	//获取编辑人的uid
	updaterUids := GetFieldValues(slice, "UpdaterId")
	for _, v := range updaterUids {
		if uid, ok := v.(int64); ok {
			adminUids = append(adminUids, uid)
		}
	}
	//获取管理员列表
	adminList, err := srv.GetListByUIds(adminUids)
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}
	adminListM = make(map[int64]models.Admin)
	for _, v := range adminList {
		adminListM[v.AdminId] = v
	}
	return
}
