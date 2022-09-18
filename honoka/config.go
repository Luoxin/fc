package honoka

import (
	"os"

	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(utils.GetExecPath())
	viper.AddConfigPath(utils.GetPwd())

	err := viper.ReadRemoteConfig()
	if err != nil {
		switch e := err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warn("not found conf file, use default")
		case *os.PathError:
			log.Warnf("not find conf file in %s", e.Path)
		default:
			log.Debugf("load config fail:%v", err)
		}
	}

	return nil
}

func ConfigGet(key string) any {
	return viper.Get(key)
}
