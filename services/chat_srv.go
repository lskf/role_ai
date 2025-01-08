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
	"strconv"
	"strings"
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
		chatHistories        []models.ChatHistory
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
		chatShortTermHistory, err = srv.chatRepo.GetChatShortTermMemory(chat.Id)
		if err != nil {
			return nil, errors.New(ecode.DatabaseErr, err)
		}
	}
	if len(chatShortTermHistory) <= 0 {
		//第一次对话
		chatRecords = append(chatRecords, llm.MessageObj{Assistant: role.Remark})
	} else {
		for _, history := range chatShortTermHistory {
			switch history.RoleType {
			case models.ChatHistoryRoleUser:
				chatRecords = append(chatRecords, llm.MessageObj{User: history.Content})
			case models.ChatHistoryRoleAssistant:
				chatRecords = append(chatRecords, llm.MessageObj{Assistant: history.Content})
			}
		}
	}
	chatRecords = append(chatRecords, llm.MessageObj{User: para.Question})
	//拼systemSetting
	systemSetting, err = srv.spliceSystem(role, chat)
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
	content = strings.Replace(content, "{{char}}", role.RoleName, -1)
	askContent := models.Reply{Content: para.Question}
	askContentStr, _ := json.Marshal(askContent)
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
	chatHistories = make([]models.ChatHistory, 2)
	chatHistories[0].RoleType = models.ChatHistoryRoleUser
	chatHistories[0].Type = models.ChatHistoryTypeChat
	chatHistories[0].Content = para.Question
	chatHistories[0].Info = string(askContentStr)
	chatHistories[0].CreatedAt = time.Now()
	chatHistories[0].UpdatedAt = time.Now()

	chatHistories[1].RoleType = models.ChatHistoryRoleAssistant
	chatHistories[1].Type = models.ChatHistoryTypeChat
	chatHistories[1].Content = replyContent.Content
	chatHistories[1].Info = content
	chatHistories[1].CreatedAt = time.Now()
	chatHistories[1].UpdatedAt = time.Now()
	err = srv.chatRepo.Transaction(func(ctx context.Context) error {
		if chat.Id == 0 {
			//创建对话
			gamification, _ := json.Marshal(role.Gamification)
			chat.Gamification = string(gamification)
			err = srv.chatRepo.Create(&creator.Creator{
				Tx:   ctx,
				Data: &chat,
			})
			if err != nil {
				return errors.New(ecode.DatabaseErr, err)
			}
		} else {
			affection, _ := strconv.Atoi(replyContent.Affection)
			sexuality, _ := strconv.Atoi(replyContent.Sexuality)
			gamificationObj := dto.GamificationObj{
				Affection: int64(affection),
				Sexuality: int64(sexuality),
			}
			gamification, _ := json.Marshal(gamificationObj)
			chat.Gamification = string(gamification)
			err = srv.chatRepo.Update(&updater.Updater{
				Tx:     ctx,
				Model:  new(models.Chat),
				Wheres: where.New().And(where.Eq("id", chat.Id)),
				Fields: map[string]interface{}{
					"gamification": chat.Gamification,
					"updated_at":   chat.UpdatedAt,
				},
			})
			if err != nil {
				return errors.New(ecode.DatabaseErr, err)
			}
		}
		chatHistories[0].ChatId = chat.Id
		chatHistories[1].ChatId = chat.Id
		err = srv.chatRepo.Create(&creator.Creator{
			Tx:   ctx,
			Data: &chatHistories,
		})
		if err != nil {
			return errors.New(ecode.DatabaseErr, err)
		}
		//更新短期记忆,更新缓存
		err = srv.chatRepo.AddChatShortTermMemory(chat.Id, chatHistories)
		_ = srv.chatRepo.AddChatHistory(chat.Id, chatHistories)
		if err != nil {
			//如果失败的话，删除短期记忆，删除缓存,
			_ = srv.chatRepo.DelChatShortTermMemory(chat.Id)
			_ = srv.chatRepo.DelChatHistory(chat.Id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	reply := make([]dto.ChatHistory, 2)
	err = models.Copy(&reply, &chatHistories)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr, err)
	}
	reply[1].Affection = replyContent.Affection
	reply[1].Sexuality = replyContent.Sexuality
	return &reply, err
}

func (srv *ChatService) GetList(uid int64, para dto.ChatListReq) (any, error) {
	var (
		chatList  []models.Chat
		roleList  []models.Role
		total     int64
		listWhere where.Wheres
		sort      string
		roleIds   []int64
		//res       dto.ChatListResp
	)

	listWhere = listWhere.And(where.Eq("uid", uid))
	//if para.Name != "" {
	//	//获取角色id
	//
	//	listWhere = listWhere.And(where.Like("role_name", "%"+para.Name+"%"))
	//}

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
	//获取角色列表
	for _, v := range chatList {
		roleIds = append(roleIds, v.RoleId)
	}
	err = srv.roleRepo.Find(&finder.Finder{
		Model:     models.Role{},
		Wheres:    where.New().And(where.In("id", roleIds)),
		Recipient: &roleList,
		Total:     &total,
	})
	if err != nil {
		return nil, errors.New(ecode.DatabaseErr, err)
	}
	roleM := make(map[int64]models.Role)
	for _, v := range roleList {
		roleM[v.Id] = v
	}
	return nil, nil
}

func (srv *ChatService) splicePrompt(input string) (string, error) {
	return input, nil
}

func (srv *ChatService) spliceSystem(role models.Role, chat models.Chat) (string, error) {
	type SystemContent struct {
		Initialization     string `json:"initialization"`
		Setting            string `json:"setting"`
		Repeat             string `json:"repeat"`
		ConvenientFlirting string `json:"convenientFlirting"`
		RoleSetting        string `json:"roleSetting"`
		InitialState       string `json:"initialState"`
		StatusBlockRule    string `json:"statusBlockRule"`
		RepeatPro          string `json:"repeat_pro"`
		NSFW               string `json:"NSFW"`
		Request            string `json:"request"`
		SecondPerson       string `json:"secondPerson"`
		TimeFormat         string `json:"timeFormat"`
		Psychology         string `json:"psychology"`
		Style              string `json:"style"`
		FeMaleStyle        string `json:"FeMale_style"`
		MaleStyle          string `json:"Male_style"`
		Response           string `json:"response"`
		WordLimit          string `json:"wordLimit"`
		Config             string `json:"config"`
		Jailbreak          string `json:"jailbreak"`
		ReplyFormat        string `json:"replyFormat"`
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

	systemSetting := systemContent.Initialization
	systemSetting += systemContent.Setting
	systemSetting += systemContent.Repeat
	systemSetting += systemContent.ConvenientFlirting
	//角色设置
	roleSetting := systemContent.RoleSetting
	roleSetting = strings.Replace(roleSetting, "{{char_name}}", role.RoleName, -1)
	roleSetting = strings.Replace(roleSetting, "{{char_desc}}", role.Desc, -1)
	systemSetting += roleSetting

	//游戏化设置（角色状态）
	initialState := systemContent.InitialState
	gamification := dto.GamificationObj{}
	if chat.Id > 0 {
		_ = json.Unmarshal([]byte(chat.Gamification), &gamification)
	} else {
		_ = json.Unmarshal([]byte(role.Gamification), &gamification)
	}
	initialState = strings.Replace(initialState, "{{affection}}", strconv.FormatInt(gamification.Affection, 10), -1)
	initialState = strings.Replace(initialState, "{{sexuality}}", strconv.FormatInt(gamification.Sexuality, 10), -1)
	systemSetting += initialState

	systemSetting += systemContent.StatusBlockRule
	systemSetting += systemContent.RepeatPro
	systemSetting += systemContent.NSFW
	systemSetting += systemContent.Request
	systemSetting += systemContent.SecondPerson
	systemSetting += systemContent.TimeFormat
	systemSetting += systemContent.Psychology
	//回复风格
	switch chat.ReplyStyle {
	case models.ReplyStyleLyrical:
		systemSetting += systemContent.Style
	case models.ReplyStyleFemaleNsfw:
		systemSetting += systemContent.FeMaleStyle
	case models.ReplyStyleMaleNsfw:
		systemSetting += systemContent.MaleStyle
	}

	systemSetting += systemContent.Response
	//字数规模
	wordLimit := systemContent.WordLimit
	wordLimit = strings.Replace(wordLimit, "{{word_count}}", strconv.FormatInt(chat.WordCount, 10), -1)
	systemSetting += wordLimit

	systemSetting += systemContent.Config
	systemSetting += systemContent.Jailbreak

	//返回格式
	systemSetting += systemContent.ReplyFormat

	return systemSetting, nil
}