package sysconfig

import (
	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
	"io/ioutil"
	"os"
	"sigs.k8s.io/yaml"
)

// TODO: 把配置文件搞成项目某一个目录中

// CheckFileExists 检查文件是否存在
func CheckFileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func CleanConfig(sysconfig *SysConfig, fileName string) error {

	// 1. 把SysConfig1中的都删除
	// 清零后需要先更新app.yaml文件
	sysconfig.Services = make([]Services, 0)
	sysconfig.Dsn = ""
	sysconfig.MaxOpenConn = 0
	sysconfig.MaxIdleConn = 0
	var sys *SysConfig
	var err error
	if sys, err = saveConfigToFile(sysconfig, fileName); err != nil {
		return err
	}

	return ReloadConfig(sys, fileName)
}

// CreateFile 创建文件
func CreateFile(fileName string) error {
	filePath := fileName
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// GetContentFromFile 获取原配置文件的实例
func GetContentFromFile(fileName string) (*SysConfig, error) {
	filePath := fileName

	// 读取文件内容
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var sysConfig *SysConfig
	err = yaml.Unmarshal(fileContent, &sysConfig)
	if err != nil {
		return nil, err
	}

	return sysConfig, nil
}

// CreateAppConfig 创建配置文件实例
func CreateAppConfig(dbconfig *dbconfigv1alpha1.DbConfig, fileName string) error {

	// 比较当前db有的 user 与 dbname
	sysconfig := &SysConfig{
		Services: make([]Services, len(dbconfig.Spec.Services)),
	}

	// 2. 更新内存的配置
	sysconfig.Dsn = dbconfig.Spec.Dsn
	sysconfig.MaxIdleConn = dbconfig.Spec.MaxIdleConn
	sysconfig.MaxOpenConn = dbconfig.Spec.MaxOpenConn
	for i, service := range dbconfig.Spec.Services {
		sysconfig.Services[i].Service.Dbname = service.Service.Dbname
		sysconfig.Services[i].Service.Tables = service.Service.Tables
		sysconfig.Services[i].Service.User = service.Service.User
		sysconfig.Services[i].Service.Password = service.Service.Password
		sysconfig.Services[i].Service.ReBuild = service.Service.ReBuild
	}

	// 保存配置文件
	var sys *SysConfig
	var err error
	if sys, err = saveConfigToFile(sysconfig, fileName); err != nil {
		return err
	}

	return ReloadConfig(sys, fileName)
}

func AppConfig(dbconfig *dbconfigv1alpha1.DbConfig, sysconfig *SysConfig, fileName string) error {

	if sysconfig == nil {
		sysconfig = &SysConfig{
			Services: make([]Services, 0),
		}
	}

	// 比较当前db有的 user 与 dbname

	// 如果数量不同，直接全部重新赋值
	if len(sysconfig.Services) != len(dbconfig.Spec.Services) {
		// 清零后需要先更新 app.yaml 文件
		sysconfig.Services = make([]Services, len(dbconfig.Spec.Services))
		if _, err := saveConfigToFile(sysconfig, fileName); err != nil {
			return err
		}
	}

	// 2. 更新内存的配置
	sysconfig.Dsn = dbconfig.Spec.Dsn
	sysconfig.MaxIdleConn = dbconfig.Spec.MaxIdleConn
	sysconfig.MaxOpenConn = dbconfig.Spec.MaxOpenConn
	for i, service := range dbconfig.Spec.Services {
		sysconfig.Services[i].Service.Dbname = service.Service.Dbname
		sysconfig.Services[i].Service.Tables = service.Service.Tables
		sysconfig.Services[i].Service.User = service.Service.User
		sysconfig.Services[i].Service.Password = service.Service.Password
		sysconfig.Services[i].Service.ReBuild = service.Service.ReBuild
	}

	// 保存配置文件
	var sys *SysConfig
	var err error
	if sys, err = saveConfigToFile(sysconfig, fileName); err != nil {
		return err
	}

	return ReloadConfig(sys, fileName)
}

// saveConfigToFile 把 config 配置放入文件中
func saveConfigToFile(sysconfig *SysConfig, fileName string) (*SysConfig, error) {

	b, err := yaml.Marshal(sysconfig)
	if err != nil {
		return nil, err
	}
	filePath := fileName
	appYamlFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 644)
	if err != nil {
		return nil, err
	}

	defer appYamlFile.Close()
	_, err = appYamlFile.Write(b)
	if err != nil {
		return nil, err
	}

	return sysconfig, nil
}

// ReloadConfig 重载配置
func ReloadConfig(sysconfig *SysConfig, fileName string) error {
	return InitConfig(sysconfig, fileName)
}

func InitConfig(sysconfig *SysConfig, fileName string) error {
	// 读取 yaml 配置
	config, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(config, sysconfig)
	if err != nil {
		return err
	}
	return nil
}
