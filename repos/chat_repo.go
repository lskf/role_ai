package repos

import (
	"context"
	"encoding/json"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/redis"
	"github.com/leor-w/kid/database/repos"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/where"
	"role_ai/models"
	"time"
)

type IChatRepository interface {
	repos.IBasicRepository
	AddChatShortTermHistory(chatId int64, history *models.ChatHistory) error
	GetChatShortTermHistory(chatId int64) ([]*models.ChatHistory, error)
	DelChatShortTermHistory(chatId int64) error
}
type ChatRepository struct {
	*mysql.Repository      `inject:""`
	*redis.RedisRepository `inject:""`
}

func (repo *ChatRepository) Provide(context.Context) any {
	return &ChatRepository{}
}

var (
	chatShortTermHistoryKey = "history.short.chat.%d"
	chatLongTermHistoryKey  = "history.long.chat.%d"
)

func (repo *ChatRepository) AddChatShortTermHistory(chatId int64, history *models.ChatHistory) error {
	key := fmt.Sprintf(chatShortTermHistoryKey, chatId)
	historyStr, err := json.Marshal(history)
	if err != nil {
		return err
	}
	llen, err := repo.RDB.LLen(key).Result()
	if err != nil {
		return err
	}
	if _, err = repo.RDB.Pipelined(func(pipe redis2.Pipeliner) error {
		popCount := int(llen) - 30
		for i := 0; i < popCount; i++ {
			if err = repo.RDB.LPop(key).Err(); err != nil {
				return err
			}
		}
		if err = repo.RDB.RPush(key, historyStr).Err(); err != nil {
			return err
		}
		//添加过期时间
		repo.RDB.Expire(key, time.Minute*60*24*7) //7天过期时间
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (repo *ChatRepository) GetChatShortTermHistory(chatId int64) ([]*models.ChatHistory, error) {
	var histories []*models.ChatHistory
	key := fmt.Sprintf(chatShortTermHistoryKey, chatId)
	res, err := repo.RDB.LRange(key, 0, -1).Result()
	if err != nil || len(res) == 0 {
		//查数据库
		err = repo.Find(&finder.Finder{
			Model:          new(models.ChatHistory),
			Wheres:         where.New().And(where.Eq("chat_id", chatId)),
			Recipient:      &histories,
			OrderBy:        "id desc",
			Num:            1,
			Size:           30, //短期记忆只取最新的30条
			IgnoreNotFound: true,
		})
		//将记录添加回redis中
		if err != nil {
			return nil, err
		}
		value := make([]interface{}, 0)
		for _, v := range histories {
			history, _ := json.Marshal(v)
			value = append(value, string(history))
		}
		repo.RDB.LPush(key, value...)
		repo.RDB.Expire(key, time.Minute*60*24*7) //7天过期时间
		return histories, err
	}
	for _, v := range res {
		var history models.ChatHistory
		err = json.Unmarshal([]byte(v), &history)
		if err != nil {
			continue
		}
		histories = append(histories, &history)
	}
	return histories, nil
}

func (repo *ChatRepository) DelChatShortTermHistory(chatId int64) error {
	key := fmt.Sprintf(chatShortTermHistoryKey, chatId)
	err := repo.RDB.LTrim(key, 1, 0).Err()
	return err
}
