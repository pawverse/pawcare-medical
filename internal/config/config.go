package config

import (
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitConfig() viper.Viper {
	viper := viper.New()

	viper.SetEnvPrefix("PAWCARE")
	viper.AddConfigPath("$HOME/.config/pawcare")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.BindEnv(common.HTTPPortKey)
	viper.BindEnv(common.DBHostKey)
	viper.BindEnv(common.DBPortKey)
	viper.BindEnv(common.DBUserKey)
	viper.BindEnv(common.DBPassKey)
	viper.BindEnv(common.CertsEndpointKey)
	viper.BindEnv(common.InsecureSkipVerifyKey)
	viper.BindEnv(common.KafkaBrokersKey)

	viper.SetDefault(common.HTTPPortKey, "80")
	viper.SetDefault(common.DBHostKey, "localhost")
	viper.SetDefault(common.DBPortKey, "27017")
	viper.SetDefault(common.DBUserKey, "root")
	viper.SetDefault(common.DBPassKey, "root")
	viper.SetDefault(common.InsecureSkipVerifyKey, false)
	viper.SetDefault(common.KafkaBrokersKey, "kafka:9092")

	pflag.String("http-port", "", "HTTP Port")
	pflag.Parse()

	viper.BindPFlag(common.HTTPPortKey, pflag.Lookup("http-port"))
	common.SetConnectionStringConfig(viper)

	return *viper
}
