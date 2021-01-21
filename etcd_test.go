package getcd_test

import (
	getcd "go-etcd"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Endpoints   []string
	DialTimeout time.Duration
}

var config *Config

func loadYml(path string, t interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, t)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	loadYml("conf/config.yml", &config)
}

func TestEtcd(t *testing.T) {
	e, err := getcd.New(config.Endpoints...)
	if err != nil {
		panic(err)
	}
	putResp, err := e.Set("/hello/world", "test123456")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("set /hello/world:", putResp)

	getResp, err := e.Get("/hello/world")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("get /hello/world:", getResp)
}

// go test --run="TestSetConfig"
func TestSetConfig(t *testing.T) {
	e, err := getcd.New(config.Endpoints...)
	if err != nil {
		panic(err)
	}

	e.Set("/db/user", "root")
	e.Set("/db/pass", "123456")
	e.Set("/db/host", "127.0.0.1")
	e.Set("/db/port", "3306")
	e.Set("/db/name", "test")

	showConfig()
}

func showConfig() {
	e, err := getcd.New(config.Endpoints...)
	if err != nil {
		panic(err)
	}
	db_user, _ := e.Get("/db/user")
	db_pass, _ := e.Get("/db/pass")
	db_host, _ := e.Get("/db/host")
	db_port, _ := e.Get("/db/port")
	db_name, _ := e.Get("/db/name")

	log.Println("db user\t\t:", db_user)
	log.Println("db pass\t\t:", db_pass)
	log.Println("db host\t\t:", db_host)
	log.Println("db port\t\t:", db_port)
	log.Println("db name\t\t:", db_name)
}
