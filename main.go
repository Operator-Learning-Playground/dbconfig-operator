package main

import (
	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
	"github.com/myoperator/dbconfigoperator/pkg/controller"
	"github.com/myoperator/dbconfigoperator/pkg/k8sconfig"
	"github.com/myoperator/dbconfigoperator/pkg/sysconfig"
	_ "k8s.io/code-generator"
	"k8s.io/klog/v2"
	"log"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

/*
	manager 主要用来管理Controller Admission Webhook 包括：
	访问资源对象的client cache scheme 并提供依赖注入机制 优雅关闭机制

	operator = crd + controller + webhook
*/



func main() {

	logf.SetLogger(zap.New())
	// 1. 管理器初始化
	mgr, err := manager.New(k8sconfig.K8sRestConfig(), manager.Options{
		Logger:  logf.Log.WithName("dbconfig-operator"),
	})
	if err != nil {
		mgr.GetLogger().Error(err, "unable to set up manager")
		os.Exit(1)
	}

	// 2. ++ 注册进入序列化表
	err = dbconfigv1alpha1.SchemeBuilder.AddToScheme(mgr.GetScheme())
	if err != nil {
		klog.Error(err, "unable add schema")
		os.Exit(1)
	}



	// 3. 控制器相关
	dbConfigCtl := controller.NewDbConfigController()

	err = builder.ControllerManagedBy(mgr).
		For(&dbconfigv1alpha1.DbConfig{}).
		Complete(dbConfigCtl)

	// 4. 载入业务配置
	if err = sysconfig.InitConfig(); err != nil {
		klog.Error(err, "unable to load sysconfig")
		os.Exit(1)
	}
	errC := make(chan error)

	// 5. 启动controller管理器
	go func() {
		klog.Info("controller start!! ")
		if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
			errC <-err
		}
	}()



	// 这里会阻塞，两种常驻进程可以使用这个方法
	getError := <-errC
	log.Println(getError.Error())

}


