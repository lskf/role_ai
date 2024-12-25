package repos

import (
	"context"
	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/repos"
)

type IVoiceRepository interface {
	repos.IBasicRepository
}
type VoiceRepository struct {
	*mysql.Repository `inject:""`
}

func (repo *VoiceRepository) Provide(context.Context) any {
	return &VoiceRepository{}
}
