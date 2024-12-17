package constant

import (
	"fmt"
	"time"
)

func GetChatHistoryTableName(date time.Time) string {
	return fmt.Sprintf("chat_histories_%s", date.Format("2006_01"))
}
