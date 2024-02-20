package config

import "github.com/spf13/viper"

var C *viper.Viper

func InitConfig() error {
	C = viper.New()
	C.AddConfigPath("./")
	C.SetConfigFile("config.yaml")
	C.SetConfigType("yaml")
	if err := C.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
