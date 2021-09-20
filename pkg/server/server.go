package server

import (
	"fmt"

	"github.com/kmpp/pkg/config"
	"github.com/kmpp/pkg/cron"
	"github.com/kmpp/pkg/data"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/encrypt"
	"github.com/kmpp/pkg/logger"
	"github.com/kmpp/pkg/migrate"
	"github.com/kmpp/pkg/plugin"
	"github.com/kmpp/pkg/plugin/xpack"
	"github.com/kmpp/pkg/router"
	"github.com/kmpp/pkg/server/hook"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

type Phase interface {
	Init() error
	PhaseName() string
}

func Phases() []Phase {
	return []Phase{
		&encrypt.InitEncryptPhase{
			Multilevel: viper.GetStringMap("encrypt.multilevel"),
		},
		&db.InitDBPhase{
			Host:         viper.GetString("db.host"),
			Port:         viper.GetInt("db.port"),
			Name:         viper.GetString("db.name"),
			User:         viper.GetString("db.user"),
			Password:     viper.GetString("db.password"),
			MaxOpenConns: viper.GetInt("db.max_open_conns"),
			MaxIdleConns: viper.GetInt("db.max_idle_conns"),
		},
		&migrate.InitMigrateDBPhase{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Name:     viper.GetString("db.name"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
		},
		&data.InitDataPhase{},
		&plugin.InitPluginDBPhase{},
		&cron.InitCronPhase{
			Enable: viper.GetBool("cron.enable"),
		},
	}
}

func Start() error {
	config.Init()
	logger.Init()
	phases := Phases()
	for _, phase := range phases {
		if err := phase.Init(); err != nil {
			logger.Log.Errorf("start phase [%v] failed reason: %s",
				phase, err.Error())
			return err
		}
		logger.Log.Infof("start phase [%s] success", phase.PhaseName())
	}
	s := router.Server()
	// load xpack plugin must behead router init,so can not create an phase for it.
	if err := xpack.LoadXpackPlugin(); err != nil {
		logger.Log.Error("xpack load failed, xpack can not be registered")
	}
	bind := fmt.Sprintf("%s:%d",
		viper.GetString("bind.host"),
		viper.GetInt("bind.port"))

	if err := hook.BeforeApplicationStart.Run(); err != nil {
		return err
	}
	return s.Run(iris.Addr(bind))
}
