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
	AddChatShortTermMemory(chatId int64, histories []models.ChatHistory) error
	GetChatShortTermMemory(chatId int64) ([]*models.ChatHistory, error)
	DelChatShortTermMemory(chatId int64) error
	AddChatHistory(chatId int64, histories []models.ChatHistory) error
	GetChatHistory(chatId int64) ([]*models.ChatHistory, error)
	GetLatestChatHistory(chatId int64) ([]*models.ChatHistory, error)
	DelChatHistory(chatId int64) error
}
type ChatRepository struct {
	*mysql.Repository      `inject:""`
	*redis.RedisRepository `inject:""`
}

func (repo *ChatRepository) Provide(context.Context) any {
	return &ChatRepository{}
}

var (
	chatShortTermMemoryKey = "chat.shortTerm.memory.%d"
	chatHistoryKey         = "chat.history.%d"
)

func (repo *ChatRepository) AddChatShortTermMemory(chatId int64, histories []models.ChatHistory) error {
	key := fmt.Sprintf(chatShortTermMemoryKey, chatId)
	historyCount := len(histories)
	value := make([]interface{}, 0)
	for _, history := range histories {
		historyStr, err := json.Marshal(history)
		if err != nil {
			return err
		}
		value = append(value, string(historyStr))
	}

	llen, err := repo.RDB.LLen(key).Result()
	if err != nil {
		return err
	}
	if _, err = repo.RDB.Pipelined(func(pipe redis2.Pipeliner) error {
		popCount := int(llen) - (60 - historyCount)
		for i := 0; i < popCount; i++ {
			if err = repo.RDB.LPop(key).Err(); err != nil {
				return err
			}
		}
		if err = repo.RDB.RPush(key, value...).Err(); err != nil {
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

func (repo *ChatRepository) GetChatShortTermMemory(chatId int64) ([]*models.ChatHistory, error) {
	var histories []*models.ChatHistory
	key := fmt.Sprintf(chatShortTermMemoryKey, chatId)
	res, err := repo.RDB.LRange(key, 0, -1).Result()
	if err != nil || len(res) == 0 {
		//查数据库
		var newHistories []*models.ChatHistory //需要进行反转顺序的
		err = repo.Find(&finder.Finder{
			Model:          new(models.ChatHistory),
			Wheres:         where.New().And(where.Eq("chat_id", chatId), where.Eq("type", 1)),
			Recipient:      &newHistories,
			OrderBy:        "id desc",
			Num:            1,
			Size:           60, //短期记忆只取最新的60条（30组对话）
			IgnoreNotFound: true,
		})
		//将记录添加回redis中
		if err != nil {
			return nil, err
		}
		historyCount := len(newHistories)
		histories = make([]*models.ChatHistory, historyCount)
		value := make([]interface{}, 0)
		for k, v := range newHistories {
			history, _ := json.Marshal(v)
			value = append(value, string(history))
			histories[historyCount-k-1] = v
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

func (repo *ChatRepository) DelChatShortTermMemory(chatId int64) error {
	key := fmt.Sprintf(chatShortTermMemoryKey, chatId)
	err := repo.RDB.Del(key).Err()
	return err
}

func (repo *ChatRepository) AddChatHistory(chatId int64, histories []models.ChatHistory) error {
	key := fmt.Sprintf(chatHistoryKey, chatId)
	historyCount := len(histories)
	value := make([]interface{}, 0)
	for _, history := range histories {
		historyStr, err := json.Marshal(history)
		if err != nil {
			return err
		}
		value = append(value, string(historyStr))
	}

	llen, err := repo.RDB.LLen(key).Result()
	if err != nil {
		return err
	}
	if _, err = repo.RDB.Pipelined(func(pipe redis2.Pipeliner) error {
		popCount := int(llen) - (20 - historyCount)
		for i := 0; i < popCount; i++ {
			if err = repo.RDB.RPop(key).Err(); err != nil {
				return err
			}
		}
		if err = repo.RDB.LPush(key, value...).Err(); err != nil {
			return err
		}
		//添加过期时间
		repo.RDB.Expire(key, time.Minute*60*24) //1天过期时间
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (repo *ChatRepository) GetChatHistory(chatId int64) ([]*models.ChatHistory, error) {
	var histories []*models.ChatHistory
	key := fmt.Sprintf(chatHistoryKey, chatId)
	res, err := repo.RDB.LRange(key, 0, -1).Result()
	if err != nil || len(res) == 0 {
		//查数据库
		err = repo.Find(&finder.Finder{
			Model:          new(models.ChatHistory),
			Wheres:         where.New().And(where.Eq("chat_id", chatId)),
			Recipient:      &histories,
			OrderBy:        "id desc",
			Num:            1,
			Size:           20, //缓存最新的20条消息
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
		_, _ = repo.RDB.Pipelined(func(pipe redis2.Pipeliner) error {
			repo.RDB.Del(key)
			repo.RDB.RPush(key, value...)
			repo.RDB.Expire(key, time.Minute*60*24) //1天过期时间
			return nil
		})
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

func (repo *ChatRepository) GetLatestChatHistory(chatId int64) ([]*models.ChatHistory, error) {
	var histories []*models.ChatHistory
	key := fmt.Sprintf(chatHistoryKey, chatId)
	res, err := repo.RDB.LIndex(key, 0).Result()
	if err != nil || len(res) == 0 {
		//查数据库
		err = repo.Find(&finder.Finder{
			Model:          new(models.ChatHistory),
			Wheres:         where.New().And(where.Eq("chat_id", chatId)),
			Recipient:      &histories,
			OrderBy:        "id desc",
			Num:            1,
			Size:           20, //缓存最新的20条消息
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
		_, _ = repo.RDB.Pipelined(func(pipe redis2.Pipeliner) error {
			repo.RDB.Del(key)
			repo.RDB.RPush(key, value...)
			repo.RDB.Expire(key, time.Minute*60*24) //1天过期时间
			return nil
		})
		if len(histories) > 0 {
			return histories[0:1], err
		}
		return nil, err
	}
	var history models.ChatHistory
	_ = json.Unmarshal([]byte(res), &history)
	histories = append(histories, &history)
	return histories, nil
}

func (repo *ChatRepository) DelChatHistory(chatId int64) error {
	key := fmt.Sprintf(chatHistoryKey, chatId)
	err := repo.RDB.Del(key).Err()
	return err
}
