package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// doesFileExist returns true if the file exists.
// If there is an error, it will be of type *PathError.
func doesFileExist(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if info == nil {
		return false, nil
	}

	if info.IsDir() {
		return false, nil
	}
	return true, nil
}

func GetPGInfo() string {
	databasestring := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.name"),
	)

	log.Println(databasestring)
	return databasestring
}

func ConfigureViper(environment string) error {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("ENTE")

	viper.EnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile("configuration/" + environment + ".yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
