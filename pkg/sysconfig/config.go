package sysconfig

import (
	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
)

var SysConfig1 = new(SysConfig)

type SysConfig struct {
	Dsn         string    `yaml:"dsn"`
	MaxIdleConn int       `yaml:"maxIdleConn"`
	MaxOpenConn int       `yaml:"maxOpenConn"`
	Services    []Service `yaml:"services"`
}

type Service struct {
	Dbname   string `yaml:"dbname"`
	Tables   string `yaml:"tables"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	ReBuild  bool   `yaml:"rebuild"`
}

//type Services struct {
//	Service Service `yaml:"service"`
//}

/*
	NOTE: 上一版使用 app.yaml 全局配置文件来管理 cr 的配置，
	目前更新后，集群可接受创建多个 cr 资源，并把配置记录在 global_config 目录下
*/

//func CleanConfig() error {
//
//	// 1. 把SysConfig1中的都删除
//	// 清零后需要先更新app.yaml文件
//	SysConfig1.Services = make([]Services, 0)
//	if err := saveConfigToFile(); err != nil {
//		return err
//	}
//
//	return ReloadConfig()
//}

//// AppConfig 刷新配置文件
//func AppConfig(dbconfig *dbconfigv1alpha1.DbConfig) error {
//
//	// 比较当前db有的 user 与 dbname
//
//	// 如果数量不同，直接全部重新赋值
//	if len(SysConfig1.Services) != len(dbconfig.Spec.Services) {
//		// 清零后需要先更新 app.yaml 文件
//		SysConfig1.Services = make([]Services, len(dbconfig.Spec.Services))
//		if err := saveConfigToFile(); err != nil {
//			return err
//		}
//	}
//
//	// 2. 更新内存的配置
//	SysConfig1.Dsn = dbconfig.Spec.Dsn
//	SysConfig1.MaxIdleConn = dbconfig.Spec.MaxIdleConn
//	SysConfig1.MaxOpenConn = dbconfig.Spec.MaxOpenConn
//	for i, service := range dbconfig.Spec.Services {
//		SysConfig1.Services[i].Service.Dbname = service.Service.Dbname
//		SysConfig1.Services[i].Service.Tables = service.Service.Tables
//		SysConfig1.Services[i].Service.User = service.Service.User
//		SysConfig1.Services[i].Service.Password = service.Service.Password
//		SysConfig1.Services[i].Service.ReBuild = service.Service.ReBuild
//	}
//
//	// 保存配置文件
//	if err := saveConfigToFile(); err != nil {
//		return err
//	}
//
//	return ReloadConfig()
//}

// ReloadConfig 重载配置
//func ReloadConfig() error {
//	return InitConfig()
//}
//
//func InitConfig() error {
//	// 读取 yaml 配置
//	config, err := ioutil.ReadFile("./app.yaml")
//	if err != nil {
//		return err
//	}
//
//	err = yaml.Unmarshal(config, SysConfig1)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// saveConfigToFile 把 config 配置放入文件中
//func saveConfigToFile() error {
//
//	b, err := yaml.Marshal(SysConfig1)
//	if err != nil {
//		return err
//	}
//	// 读取文件
//	path := common.GetWd()
//	filePath := path + "/app.yaml"
//	appYamlFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 644)
//	if err != nil {
//		return err
//	}
//
//	defer appYamlFile.Close()
//	_, err = appYamlFile.Write(b)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// CompareNeedToDelete 比较原本的config与新增的资源对象config，如果在新的资源中没有找到，代码需要删除此表与用户
func CompareNeedToDelete(dbconfig *dbconfigv1alpha1.DbConfig, sysconfig *SysConfig) ([]string, []string) {
	needDeleteDb := make([]string, 0)
	needDeleteUser := make([]string, 0)

	// 挑出需要删除的资源
	for _, v := range sysconfig.Services {
		isDb := searchDbNotInList(v.Dbname, dbconfig.Spec.Services)
		isUser := searchUserNotInList(v.User, dbconfig.Spec.Services)
		if isDb {
			needDeleteDb = append(needDeleteDb, v.Dbname)
		}
		if isUser {
			needDeleteUser = append(needDeleteUser, v.User)
		}
	}

	return needDeleteDb, needDeleteUser
}

// searchDbNotInList 比较 dbconfig 中是否有此 dbname，
// 如果没有代表需要在 app.yaml 中挑出来，准备删除的 dbname
func searchDbNotInList(target string, services []dbconfigv1alpha1.Service) bool {
	for _, v := range services {
		if v.Dbname == target {
			return false
		}
	}
	return true
}

func searchUserNotInList(target string, services []dbconfigv1alpha1.Service) bool {
	for _, v := range services {
		if v.User == target {
			return false
		}
	}
	return true
}
