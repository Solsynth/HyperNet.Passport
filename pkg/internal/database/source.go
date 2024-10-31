package database

import (
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/cruda"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var C *gorm.DB

func NewGorm() error {
	dsn, err := cruda.NewCrudaConn(gap.Nx).AllocDatabase("passport")
	if err != nil {
		return fmt.Errorf("failed to alloc database from nexus: %v", err)
	}

	C, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.New(&log.Logger, logger.Config{
		Colorful:                  true,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  lo.Ternary(viper.GetBool("debug.database"), logger.Info, logger.Silent),
	})})

	return err
}
