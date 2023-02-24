package conf

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

// any approach to require this configuration into your program.
var defaultSetting = []byte(`
server: 0.0.0.0:3001
name: File Management Service
`)

func init() {
	fmt.Println("Tải cấu hình ...")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf") // thư mục chứa file config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// File config không tồn tại
			fmt.Println("Không tìm thấy file cấu hình!")
			fmt.Println("Sử dụng cấu hình mặc định")
			viper.ReadConfig(bytes.NewBuffer(defaultSetting))
			viper.SafeWriteConfig()
		} else {
			// File có tồn tại nhưng không đọc được
			panic(fmt.Errorf("Lỗi tải cấu hình: %w", err))
		}
	}

	// Hỗ trợ đọc từ biến môi trường
	viper.SetEnvPrefix("fms")
	viper.AutomaticEnv()
}

func GetServerConfig() map[string]string {
	setting := viper.GetStringMapString("server")
	return setting
}
