package lib

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	vp := viper.New()
	vp.SetConfigFile("./.env")
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file: %v", err)
	}

	os.Setenv("STORAGE_EMULATOR_HOST", vp.GetString("STORAGE_EMULATOR_HOST"))
	os.Setenv(`PUBSUB_EMULATOR_HOST`, vp.GetString(`PUBSUB_EMULATOR_HOST`))
	vp.AutomaticEnv()

	return vp
}
