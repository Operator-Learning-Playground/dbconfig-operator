package sysconfig

import (
	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
	"github.com/myoperator/dbconfigoperator/pkg/common"
	"io/ioutil"
	"os"
	"sigs.k8s.io/yaml"
)



var SysConfig1 = new(SysConfig)

func InitConfig() error {
	// 读取yaml配置
	config, err := ioutil.ReadFile("./app.yaml")
	if err != nil {
		return err
	}


	err = yaml.Unmarshal(config, SysConfig1)
	if err != nil {
		return err
	}


	return nil

}

type SysConfig struct {
	Dsn    		string  `yaml:"dsn"`
	MaxIdleConn int  	`yaml:"maxIdleConn"`
	Services    []Services `yaml:"services"`
}

type Service struct {
	Dbname   string   `yaml:"dbname"`
	Tables   string   `yaml:"tables"`
	User     string	  `yaml:"user"`
	Password string   `yaml:"password"`
}

type Services struct {
	Service Service   `yaml:"service"`
}




// AppConfig 刷新配置文件
func AppConfig(dbconfig *dbconfigv1alpha1.DbConfig) error {

	// 比较目前db有的 user 与dbname

	if len(SysConfig1.Services) != len(dbconfig.Spec.Services) {
		// 清零后需要先更新app.yaml文件
		SysConfig1.Services = make([]Services, len(dbconfig.Spec.Services))
		if err := saveConfigToFile(); err != nil {
			return err
		}
	}

	// 2. 更新内存的配置
	SysConfig1.Dsn = dbconfig.Spec.Dsn
	SysConfig1.MaxIdleConn = dbconfig.Spec.MaxIdleConn
	for i, service := range dbconfig.Spec.Services {
		SysConfig1.Services[i].Service.Dbname = service.Service.Dbname
		SysConfig1.Services[i].Service.Tables = service.Service.Tables
		SysConfig1.Services[i].Service.User = service.Service.User
		SysConfig1.Services[i].Service.Password = service.Service.Password
	}

	// 保存配置文件
	if err := saveConfigToFile(); err != nil {
		return err
	}

	return ReloadConfig()
}

// ReloadConfig 重载配置
func ReloadConfig() error {
	return InitConfig()
}

//saveConfigToFile 把config配置放入文件中
func saveConfigToFile() error {

	b, err := yaml.Marshal(SysConfig1)
	if err != nil {
		return err
	}
	// 读取文件
	path := common.GetWd()
	filePath := path + "/app.yaml"
	appYamlFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 644)
	if err != nil {
		return err
	}

	defer appYamlFile.Close()
	_, err = appYamlFile.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// CompareNeedToDelete 比较原本的config与新增的资源对象config，如果在新的资源中没有找到，代码需要删除此表与用户
func CompareNeedToDelete(dbconfig *dbconfigv1alpha1.DbConfig, sysconfig *SysConfig) ([]string, []string) {
	needDeleteDb := make([]string, 0)
	needDeleteUser := make([]string, 0)

	for _, v := range sysconfig.Services {
		isDb := searchDbNotInList(v.Service.Dbname, dbconfig.Spec.Services)
		isUser := searchUserNotInList(v.Service.User, dbconfig.Spec.Services)
		if isDb {
			needDeleteDb = append(needDeleteDb, v.Service.Dbname)
		}
		if isUser {
			needDeleteUser = append(needDeleteUser, v.Service.User)
		}
	}

	return needDeleteDb, needDeleteUser

}

func searchDbNotInList(target string, services []dbconfigv1alpha1.Services) bool {
	for _, v := range services {
		if v.Service.Dbname == target {
			return false
		}
	}
	return true
}

func searchUserNotInList(target string, services []dbconfigv1alpha1.Services) bool {
	for _, v := range services {
		if v.Service.User == target {
			return false
		}
	}
	return true
}