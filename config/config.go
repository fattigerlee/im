package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type RpcAddr struct {
	LogicServerAddr   string
	ConnectServerAddr string
}

func GetRpcAddr() *RpcAddr {
	return &RpcAddr{
		LogicServerAddr:   viper.GetString("rpc_addr.logic_server_addr"),
		ConnectServerAddr: viper.GetString("rpc_addr.connect_server_addr"),
	}
}

type FileServer struct {
	LocalAddr    string // 内网监听地址
	WideAddr     string // 外网监听地址
	ResourcePath string // 资源路径
	LogFilePath  string // 日志文件路径
	LogLevel     string // 日志等级
	LogTarget    string // 日志输出目标
}

func GetFileServer() *FileServer {
	info := &FileServer{
		LocalAddr:    viper.GetString("file_server.local_addr"),
		WideAddr:     viper.GetString("file_server.wide_addr"),
		ResourcePath: viper.GetString("file_server.resource_path"),
		LogFilePath:  viper.GetString("file_server.log_file_path"),
		LogLevel:     viper.GetString("file_server.log_level"),
		LogTarget:    viper.GetString("file_server.log_target"),
	}
	return info
}

type LogicServer struct {
	LocalAddr   string // 内网监听地址
	LogFilePath string // 日志文件路径
	LogLevel    string // 日志等级
	LogTarget   string // 日志输出目标
}

func GetLogicServer() *LogicServer {
	info := &LogicServer{
		LocalAddr:   viper.GetString("logic_server.local_addr"),
		LogFilePath: viper.GetString("logic_server.log_file_path"),
		LogLevel:    viper.GetString("logic_server.log_level"),
		LogTarget:   viper.GetString("logic_server.log_target"),
	}
	return info
}

type ConnectServer struct {
	LocalAddr    string // 内网监听地址
	LocalWsAddr  string // 监听websocket地址
	SubscribeNum int    // 消息订阅数量
	LogFilePath  string // 日志文件路径
	LogLevel     string // 日志等级
	LogTarget    string // 日志输出目标
}

func GetConnectServer() *ConnectServer {
	return &ConnectServer{
		LocalAddr:    viper.GetString("connect_server.local_addr"),
		LocalWsAddr:  viper.GetString("connect_server.local_ws_addr"),
		SubscribeNum: viper.GetInt("connect_server.subscribe_num"),
		LogFilePath:  viper.GetString("connect_server.log_file_path"),
		LogLevel:     viper.GetString("connect_server.log_level"),
		LogTarget:    viper.GetString("connect_server.log_target"),
	}
}

type Mysql struct {
	Dsn              string
	Host             string
	Port             uint
	DBName           string
	Username         string
	Password         string
	AutoCreateDB     bool
	AutoMigrateTable bool
	Debug            bool
}

func GetMysql() *Mysql {
	mysql := &Mysql{
		Host:             viper.GetString("mysql.host"),
		Port:             viper.GetUint("mysql.port"),
		DBName:           viper.GetString("mysql.db_name"),
		Username:         viper.GetString("mysql.username"),
		Password:         viper.GetString("mysql.password"),
		AutoCreateDB:     viper.GetBool("mysql.auto_create_db"),
		AutoMigrateTable: viper.GetBool("mysql.auto_migrate_table"),
		Debug:            viper.GetBool("mysql.debug"),
	}
	mysql.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=PRC", mysql.Username, mysql.Password, mysql.Host, mysql.Port, mysql.DBName)
	return mysql
}

type Redis struct {
	Addr     string
	Password string
}

func GetRedis() *Redis {
	return &Redis{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
	}
}

func Init(filePath string) {
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("read config error: %v", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("config change, reload failed.")
			return
		}
		fmt.Println("config change, reload success.")
	})

	fmt.Printf("using config file: %s\n", viper.ConfigFileUsed())
}
