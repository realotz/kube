package main

import (
	"github.com/go-kratos/kratos/v2/config"
	kubeConfig "github.com/go-kratos/kube/config"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

// 部署在mesh namespace 下configmap
const Ayaml = `database:
  mysql:
    dsn: "root:Test@tcp(mysql.database.svc.cluster.local:3306)/test?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
    active: 20
    idle: 10
    idle_timeout: 3600
  redis:
    addr: "redis-master.redis.svc.cluster.local:6379"
    password: ""
    db: 4`

const Byaml = `application:
  expire: 3600`

func main() {
	conf := config.New(config.WithSource(
		kubeConfig.NewSource(kubeConfig.SourceOption{
			Namespace:     "mesh",
			LabelSelector: "app=test",
			KubeConfig: filepath.Join(homedir.HomeDir(), ".kube", "config"),
		})),
	)
	err := conf.Load()
	if err != nil {
		log.Panic(err)
	}
}
