package services

import (
	"context"
	"encoding/json"
	"github.com/leor-w/kid/database/repos/creator"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/updater"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"
	"io"
	"os"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/infrastructure/llm"
	"role_ai/models"
	"role_ai/repos"
	"time"
)

type ChatService struct {
	chatRepo repos.IChatRepository `inject:""`
	roleRepo repos.IRoleRepository `inject:""`
}

func (srv *ChatService) Provide(_ context.Context) any {
	return srv
}

func (srv *ChatService) Chat(uid int64, para dto.ChatReq) (any, error) {
	var (
		role                 models.Role
		roleStyle            models.RoleStyle
		chat                 models.Chat
		chatHistory          models.ChatHistory
		chatShortTermHistory []*models.ChatHistory //短期记忆

		chatRecords   []llm.MessageObj
		systemSetting string
	)
	//获取角色详情
	err := srv.roleRepo.GetOne(&finder.Finder{
		Model:     new(models.Role),
		Wheres:    where.New().And(where.Eq("id", para.RoleId)),
		Recipient: &role,
	})
	if err != nil {
		return nil, errors.New(ecode.RoleNotExistErr, err)
	}
	if role.IsPublic == models.PrivateRole && role.Uid != uid {
		return nil, errors.New(ecode.RoleNotExistErr)
	}
	//获取角色风格
	err = srv.roleRepo.GetOne(&finder.Finder{
		Model:     new(models.RoleStyle),
		Wheres:    where.New().And(where.Eq("role_id", para.RoleId)),
		Recipient: &roleStyle,
	})
	if err != nil {
		return nil, errors.New(ecode.RoleNotExistErr, err)
	}
	speechStyleList := make([]dto.SpeechStyleObj, 0) //todo 添加对话风格
	err = json.Unmarshal([]byte(roleStyle.Content), &speechStyleList)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	//获取对话详情
	err = srv.chatRepo.GetOne(&finder.Finder{
		Model:          new(models.Chat),
		Wheres:         where.New().And(where.Eq("uid", uid), where.Eq("role_id", para.RoleId)),
		IgnoreNotFound: true,
		Recipient:      &chat,
	})
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}

	chatRecords = make([]llm.MessageObj, 0)
	//获取短期记忆
	if chat.Id > 0 {
		chatShortTermHistory, err = srv.chatRepo.GetChatShortTermHistory(chat.Id)
		if err != nil {
			return nil, errors.New(ecode.DatabaseErr, err)
		}
	}
	if len(chatShortTermHistory) <= 0 {
		//第一次对话
		chatRecords = append(chatRecords, llm.MessageObj{Assistant: role.Remark})
	} else {
		for _, history := range chatShortTermHistory {
			chatRecords = append(chatRecords, llm.MessageObj{User: history.Question, Assistant: history.Answer})
		}
	}
	chatRecords = append(chatRecords, llm.MessageObj{User: para.Question})
	//拼systemSetting
	systemSetting, err = srv.spliceSystem()
	if err != nil {
		return nil, err
	}

	//发送到llm
	claude := (&llm.Claude{}).NewClient()
	messagePara := llm.MessageReq{
		Model:         "claude-3-5-sonnet-latest",
		MaxToken:      1024,
		Messages:      chatRecords,
		SystemSetting: systemSetting,
	}
	resp, err := claude.Message(messagePara)
	if err != nil {
		return nil, errors.New(ecode.InternalErr, errors.New(ecode.ClaudeGeneratedContentErr, err))
	}
	if len(resp.Content) <= 0 {
		return nil, errors.New(ecode.InternalErr, errors.New(ecode.ClaudeGeneratedContentErr))
	}
	content := resp.Content[0].Text
	replyContent := models.Reply{}
	err = json.Unmarshal([]byte(content), &replyContent)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	if chat.Id == 0 {
		chat.Uid = uid
		chat.RoleId = para.RoleId
		chat.CreatedAt = time.Now()
		chat.UpdatedAt = time.Now()
	} else {
		chat.UpdatedAt = time.Now()
	}
	chatHistory.Question = para.Question
	chatHistory.Answer = replyContent.Content
	chatHistory.Abstract = replyContent.Details
	chatHistory.Content = content
	chatHistory.CreatedAt = time.Now()
	chatHistory.UpdatedAt = time.Now()
	err = srv.chatRepo.Transaction(func(ctx context.Context) error {
		if chat.Id == 0 {
			//创建对话
			err = srv.chatRepo.Create(&creator.Creator{
				Tx:   ctx,
				Data: &chat,
			})
			if err != nil {
				return errors.New(ecode.DatabaseErr, err)
			}
		} else {
			err = srv.chatRepo.Update(&updater.Updater{
				Tx:     ctx,
				Model:  new(models.Chat),
				Wheres: where.New().And(where.Eq("id", chat.Id)),
				Fields: map[string]interface{}{
					"updated_at": chat.UpdatedAt,
				},
			})
			if err != nil {
				return errors.New(ecode.DatabaseErr, err)
			}
		}
		chatHistory.ChatId = chat.Id
		err = srv.chatRepo.Create(&creator.Creator{
			Tx:   ctx,
			Data: &chatHistory,
		})
		if err != nil {
			return errors.New(ecode.DatabaseErr, err)
		}
		//保存到缓存中
		err = srv.chatRepo.AddChatShortTermHistory(chat.Id, &chatHistory)
		if err != nil {
			//如果失败的话，删除缓存
			_ = srv.chatRepo.DelChatShortTermHistory(chat.Id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	//更新短期记忆 todo
	//更新长期记忆
	return replyContent.Content, err
}

func (srv *ChatService) GetList(uid int64, para dto.ChatListReq) (any, error) {
	var (
		chatList  []models.Chat
		total     int64
		listWhere where.Wheres
		sort      string
		//res       dto.ChatListResp
	)

	listWhere = listWhere.And(where.Eq("uid", uid))
	if para.Name != "" {
		listWhere = listWhere.And(where.Like("role_name", "%"+para.Name+"%"))
	}

	sort += "id desc"

	err := srv.chatRepo.Find(&finder.Finder{
		Model:          models.Chat{},
		Wheres:         listWhere,
		OrderBy:        sort,
		Num:            para.PageNum,
		Size:           para.PageSize,
		Recipient:      &chatList,
		Total:          &total,
		IgnoreNotFound: true,
	})
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}
	return nil, nil
}

func (srv *ChatService) splicePrompt(input string) (string, error) {
	return input, nil
}

func (srv *ChatService) spliceSystem() (string, error) {
	type SystemContent struct {
		SystemContent   string `json:"systemContent"`
		Setting         string `json:"setting"`
		RoleSetting     string `json:"roleSetting"`
		StatusBlockRule string `json:"statusBlockRule"`
		Repeat          string `json:"repeat"`
		NSFW            string `json:"NSFW"`
		Request         string `json:"request"`
		SecondPerson    string `json:"secondPerson"`
		TimeFormat      string `json:"timeFormat"`
		Psychology      string `json:"psychology"`
		Style           string `json:"style"`
		SummaryRules    string `json:"summaryRules"`
		Response        string `json:"response"`
		Config          string `json:"config"`
		Jailbreak       string `json:"jailbreak"`
		ReplyFormat     string `json:"replyFormat"`
	}
	file, err := os.Open("./files/chat/systemContent.json")
	if err != nil {
		return "", errors.New(ecode.DataProcessingErr, err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New(ecode.DataProcessingErr, err)
	}
	var systemContent SystemContent
	err = json.Unmarshal(data, &systemContent)
	if err != nil {
		return "", errors.New(ecode.DataProcessingErr, err)
	}
	systemSetting := systemContent.SystemContent
	systemSetting += systemContent.Setting
	systemSetting += systemContent.RoleSetting
	//systemSetting += systemContent.StatusBlockRule
	systemSetting += systemContent.Repeat
	systemSetting += systemContent.NSFW
	systemSetting += systemContent.Request
	systemSetting += systemContent.SecondPerson
	systemSetting += systemContent.TimeFormat
	systemSetting += systemContent.Psychology
	systemSetting += systemContent.Style
	systemSetting += systemContent.SummaryRules
	systemSetting += systemContent.Response
	systemSetting += systemContent.Config
	systemSetting += systemContent.Jailbreak

	//返回格式
	systemSetting += systemContent.ReplyFormat

	return systemSetting, nil
}
