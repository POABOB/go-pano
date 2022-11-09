package config

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// 配置文件，Mysql、Mode...等設定
var ConfigFile = "./config.yml"

// GlobalConfig is the global config
type GlobalConfig struct {
	Server           ServerConfig   `yaml:"server"`
	Database         DatabaseConfig `yaml:"database"`
	DatabaseInDocker DatabaseConfig `yaml:"database-in-docker"`
	Python           PythonConfig   `yaml:"python"`
}

// ServerConfig is the server config
type ServerConfig struct {
	Addr      string
	Mode      string
	Version   string
	StaticDir string `yaml:"static_dir"`
	// ViewDir            string `yaml:"view_dir"`
	// UploadDir          string `yaml:"upload_dir"`
	MaxMultipartMemory int64 `yaml:"max_multipart_memory"`
}

// DatabaseConfig is the database config
type DatabaseConfig struct {
	DSN          string `yaml:"datasource"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type PythonConfig struct {
	DevHost  string `yaml:"dev_host"`
	ProdHost string `yaml:"prod_host"`
}

// global configs
var (
	Global     GlobalConfig
	Server     ServerConfig
	Database   DatabaseConfig
	PythonHost string
)

// Load config from file
func Load(file string) (GlobalConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("%v", err)
		return Global, err
	}

	err = yaml.Unmarshal(data, &Global)
	if err != nil {
		log.Printf("%v", err)
		return Global, err
	}

	Server = Global.Server
	if Global.Server.Mode == "prod" {
		Database = Global.DatabaseInDocker
		PythonHost = Global.Python.ProdHost
	} else {
		Database = Global.Database
		PythonHost = Global.Python.DevHost
	}

	return Global, nil
}

// loads configs
func init() {
	if os.Getenv("config") != "" {
		ConfigFile = os.Getenv("config")
	}
	Load(ConfigFile)
}
