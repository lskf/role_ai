package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/leor-w/kid/database/repos/creator"
	"github.com/leor-w/kid/database/repos/deleter"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/updater"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/infrastructure/llm"
	"role_ai/models"
	"role_ai/repos"
	"strconv"
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
	role.Uid = user.Uid
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
				"gender":       roleDetail.Gender,
				"desc":         roleDetail.Desc,
				"worldview":    roleDetail.Worldview,
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

func (srv *RoleService) GetRoleAvatarSetting(para *dto.AiCreateRoleReq) (*dto.GetRoleAvatarResq, error) {
	var (
		gender, storyGenre, roleType, personality, interests, preferences, dislike, background, relationships, quirks, artStyle string
	)
	if len(para.Gender) > 0 {
		gender = strings.Join(para.Gender, `","`)
		gender = `"` + gender + `"`
	}
	if len(para.StoryGenre) > 0 {
		storyGenre = strings.Join(para.StoryGenre, `","`)
		storyGenre = `"` + storyGenre + `"`
	}
	if len(para.RoleType) > 0 {
		roleType = strings.Join(para.RoleType, `","`)
		roleType = `"` + roleType + `"`
	}
	if len(para.Personality) > 0 {
		personality = strings.Join(para.Personality, `","`)
		personality = `"` + personality + `"`
	}
	if len(para.Interests) > 0 {
		interests = strings.Join(para.Interests, `","`)
		interests = `"` + interests + `"`
	}
	if len(para.Preferences) > 0 {
		preferences = strings.Join(para.Preferences, `","`)
		preferences = `"` + preferences + `"`
	}
	if len(para.Dislike) > 0 {
		dislike = strings.Join(para.Dislike, `","`)
		dislike = `"` + dislike + `"`
	}
	if len(para.Gender) > 0 {
		background = strings.Join(para.Background, `","`)
		background = `"` + background + `"`
	}
	if len(para.Relationships) > 0 {
		relationships = strings.Join(para.Relationships, `","`)
		relationships = `"` + relationships + `"`
	}
	if len(para.Quirks) > 0 {
		quirks = strings.Join(para.Quirks, `","`)
		quirks = `"` + quirks + `"`
	}

	prompt := `{"Gender":[%s],"Story Genre":[%s],"Role Type":[%s],"Personality":[%s],"Interests":[%s],"Preferences":[%s],"Dislike":[%s],"Background":[%s],"Relationships":[%s],"Quirks":[%s],"Art style":[%s]}
-The content in the "" format is the first-level column, and the content in the [] format is the corresponding label of the first-level column
[Use the above Json format instructions to generate the visual details of the NPC in the prompt form according to the selected label:]
-The reply must only include the NPC appearance description and style
-Appearance description: Provide NPC A detailed visual description of their appearance. Consider their gender, role, and personality when describing their hair, eyes, clothing, and overall style. For example, a "hero" might have a strong, athletic physique, while a "villain" might have more intimidating features.
-Style: Define the style based only on the selected tags, get the tags from "Art Style"

Response requirements:
-Only gender, clothing, facial features, and style need to be generated
-Must use concise words, no paragraphs
-Must not show unnecessary descriptions, such as "natural posture", "dignified posture", "delicate necklace", "soft facial features"
-Must not need adjectives
-Response in the form of a Vincent Prompt
-Must meet the requirements
-The more important features for the Vincent should be placed first
-Phrase phrases related to eyes and hair need to be connected with underscores, such as "black_eyes", "long_hair"
-Descriptions about hair and eyes should be separated. For example, "long brown hair" should be changed to "long hair" and "brown hair"
-Only full-body photos are generated
-Do not use duplicate descriptions in your reply. For example, "miniskirt" and "tights" cannot appear at the same time
-Commas are used to connect phrases
-White background and simple background
-When replying, only the selected tag will be displayed, and no other description is needed. For example, if the "cute" tag is selected in "Art Style", only "cute" will be replied, and there will be no other style descriptions, such as "cute anime", "1990s cute anime"
-No need to reply to negative prompts
-No need to reply to quality words of the image, such as "masterpiece", "high quality"
-No need to reply to the summary of all descriptions, such as "character appearance", "clothing", "facial features"
-If no tag is selected, a random one will be generated`
	prompt = fmt.Sprintf(prompt, gender, storyGenre, roleType, personality, interests, preferences, dislike, background, relationships, quirks, artStyle)

	artStyle = para.ArtStyle
	claude := (&llm.Claude{}).NewClient()
	messagePara := llm.MessageReq{
		Model:    "claude-3-5-sonnet-latest",
		MaxToken: 5000,
		Messages: []llm.MessageObj{
			{User: prompt},
		},
	}
	resq, err := claude.Message(messagePara)
	if err != nil {
		return nil, errors.New(ecode.InternalErr, err)
	}
	if len(resq.Content) == 0 {
		return nil, errors.New(ecode.InternalErr, errors.New(ecode.ClaudeGeneratedContentErr))
	}
	return &dto.GetRoleAvatarResq{ArtStyle: artStyle, Desc: resq.Content[0].Text}, nil
}

func (srv *RoleService) GetRoleSetting(para *dto.AiCreateRoleReq) (*dto.CreateRoleReq, error) {
	var (
		gender, storyGenre, roleType, personality, interests, preferences, dislike, background, relationships, quirks, artStyle string
	)
	if len(para.Gender) > 0 {
		gender = strings.Join(para.Gender, `","`)
		gender = `"` + gender + `"`
	}
	if len(para.StoryGenre) > 0 {
		storyGenre = strings.Join(para.StoryGenre, `","`)
		storyGenre = `"` + storyGenre + `"`
	}
	if len(para.RoleType) > 0 {
		roleType = strings.Join(para.RoleType, `","`)
		roleType = `"` + roleType + `"`
	}
	if len(para.Personality) > 0 {
		personality = strings.Join(para.Personality, `","`)
		personality = `"` + personality + `"`
	}
	if len(para.Interests) > 0 {
		interests = strings.Join(para.Interests, `","`)
		interests = `"` + interests + `"`
	}
	if len(para.Preferences) > 0 {
		preferences = strings.Join(para.Preferences, `","`)
		preferences = `"` + preferences + `"`
	}
	if len(para.Dislike) > 0 {
		dislike = strings.Join(para.Dislike, `","`)
		dislike = `"` + dislike + `"`
	}
	if len(para.Gender) > 0 {
		background = strings.Join(para.Background, `","`)
		background = `"` + background + `"`
	}
	if len(para.Relationships) > 0 {
		relationships = strings.Join(para.Relationships, `","`)
		relationships = `"` + relationships + `"`
	}
	if len(para.Quirks) > 0 {
		quirks = strings.Join(para.Quirks, `","`)
		quirks = `"` + quirks + `"`
	}

	prompt := `{"Gender":[%s],"Story Genre":[%s],"Role Type":[%s],"Personality":[%s],"Interests":[%s],"Preferences":[%s],"Dislike":[%s],"Background":[%s],"Relationships":[%s],"Quirks":[%s],"Art style":[%s]}
-The content in the "" format is the first-level column, and the content in the [] format is the corresponding label of the first-level column
[Use the above Json format instructions to generate NPC details according to the selected tags:]
-The reply must include role name, gender,  introduction, worldview and opening remarks
-Role Name: Background is NPC according to all selected tags Generate a suitable name. The name should make them feel authentically part of the world they belong to.
-Gender: Generate the gender of the NPC based on the selected gender tag (e.g., "Male", "Female"). Use appropriate pronouns and descriptions based on the selected gender. Gender should affect appearance and interactions with the world.
-Introduction: Provide a short biography about the NPC. This should include:
    - Main personality traits, tagged from "Personality" 
    - A brief biography of their background, tagged from "Background" 
    - The NPC's hobbies, preferences, and dislikes, tagged from "Interests", "Perferences", and "Dislikes"
    - Any unique traits or quirks, tagged from "Quirks" 
    - The entire biography needs to match the selected "Story Genre" and "Role Type"
-Worldview: Describe how the NPC sees the world around them, based on their personality and background. Consider whether they are optimistic, cynical, idealistic, or something else. What are their views on society, justice, or relationships?
-Opening Remarks: Generate an engaging opening or greeting that the NPC might say when they meet the protagonist. The opening should reflect the NPC's personality and situation. For example, a "hero" might say, "I represent justice and I will protect the weak," while a "villain" might say, "From the moment you met me, your fate was sealed."

Response requirements:
-The number of words in the personal introduction should be between 200 and 500 words
-Each sentence must be fluent
-The generated words must meet the requirements of the selected tags
-All selected tags must be displayed during output
-Two examples of the dialogue style of {{user}} and {{char}} must be displayed during output
-The selected tags must be consistent with the tags displayed during output
-The initial favorability and initial sexual desire value must be displayed during output
{valueType}: {char}'s favorability towards {user}
{value1}: {value1}/100 (initial favorability: {value1})
{valueType}: {char}'s sexual desire towards {user}
{value1}: {value1}/100 (initial sexual desire: {value1})
-The generated format must be output in the following Json format:
{"role_name":"testCreate","gender":"","desc":"A brave warrior with unmatched skills.","worldview":"","remark":"Welcome to my world!","tag":["fighter","hero"],"is_public":1,"style":[{"user":"I'm really angry!","role":"What's wrong? If you want, can you tell me? I hope to understand you from your life."},{"user":"I'm so happy now!","role":"Wow!!! What's so happy? I want to hear it too!!!"}],"gamification":{"affection_initial":10,"sexuality_initial":15}}
-If no tag is selected, a random tag is generated according to the Json format described above`
	prompt = fmt.Sprintf(prompt, gender, storyGenre, roleType, personality, interests, preferences, dislike, background, relationships, quirks, artStyle)

	artStyle = para.ArtStyle
	claude := (&llm.Claude{}).NewClient()
	messagePara := llm.MessageReq{
		Model:    "claude-3-5-sonnet-latest",
		MaxToken: 5000,
		Messages: []llm.MessageObj{
			{User: prompt},
		},
	}
	resq, err := claude.Message(messagePara)
	if err != nil {
		return nil, errors.New(ecode.InternalErr, err)
	}
	if len(resq.Content) == 0 {
		return nil, errors.New(ecode.InternalErr, errors.New(ecode.ClaudeGeneratedContentErr))
	}
	createRoleReq := dto.CreateRoleReq{}
	err = json.Unmarshal([]byte(resq.Content[0].Text), &createRoleReq)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	return &createRoleReq, nil
}

func (srv *RoleService) CreateRoleAvatar(uid int64, para *dto.CreateRoleAvatarReq) (any, error) {
	comfyUi := (&llm.ComfyUi{}).NewComfyUi()
	promptReq := llm.PromptReq{
		ClientId:    strconv.FormatInt(uid, 10),
		CkptName:    "juggernautXL_v9Rundiffusionphoto2.safetensors",
		PictureNum:  strconv.FormatInt(para.PictureNum, 10),
		Prompt:      para.Desc,
		ParaFileUrl: "./files/comfyUi/role_avatar/createRoleAvatarPara.json",
	}
	resp, err := comfyUi.Prompt(promptReq)
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}
	return resp, nil
}

func (srv *RoleService) GetRoleAvatarHistory(promptId string) (any, error) {
	comfyUi := (&llm.ComfyUi{}).NewComfyUi()
	resp, err := comfyUi.GetHistoryDetail(promptId)
	if err != nil {
		return nil, errors.New(ecode.InternalErr, err)
	}
	return resp, nil
}

func (srv *RoleService) GetRoleAvatar(para *dto.GetViewReq) (any, error) {
	viewReq := llm.ViewReq{
		FileName:  para.FileName,
		Type:      para.Type,
		Subfolder: para.Subfolder,
	}
	comfyUi := (&llm.ComfyUi{}).NewComfyUi()
	resp, err := comfyUi.View(viewReq)
	if err != nil {
		return nil, errors.New(ecode.InternalErr, err)
	}
	return resp, nil
}
