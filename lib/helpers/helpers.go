package helpers

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetRootMysqlConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		viper.GetString("mysql.root-username"),
		viper.GetString("mysql.root-password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		"mysql")
}

func GetReplicaMysqlConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"))
}
