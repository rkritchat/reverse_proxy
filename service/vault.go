package service

import (
	"github.com/hashicorp/vault/api"
	"log"
	"reflect"
	"reverse-proxy/common"
)

type vaultTask struct {
	vaultAddr string
	vaultToken string
	environment string
	c *api.Logical
	paths []interface{}
	hosts map[string]string
}

func (h *Handler) InitHosts() map[string]string{
	var t vaultTask
	t.initTask(h)

	conf:=t.initConfig()
	if err:=t.initClient(conf);err!=nil {
		log.Panicln(err)
	}

	t.initParentPath()
	t.initMap()
	t.logs()
	return t.hosts
}

func (t *vaultTask) logs(){
	for k,v := range t.hosts {
		log.Printf("key[%v]: %v\n", k, v)
	}
}

func (t *vaultTask) initMap(){
	for _,v := range t.paths {
		if str, ok := v.(string);ok{
			log.Printf("seach on path %v\n", str)
			read, err := t.c.Read("secret/data/"+t.environment+"/"+str+"/jenkins")
			if err!=nil{
				log.Printf("Cannot Read data on path %v\n", str)
				continue
			}
			if read!=nil {
				for k, v := range read.Data {
					log.Printf("Key[%v] : %v\n", k, v)
					t.initHost(v, str)
				}
			}
		}
	}
}

func (t *vaultTask) initHost(i interface{},name string) {
	if t.hosts == nil {
		t.hosts = make(map[string]string)
	}
	m := reflect.ValueOf(i)
	if m.Kind() == reflect.Map{
		for _,key := range m.MapKeys(){
			rs := m.MapIndex(key)
			if key.Interface() == "ip"{
				t.hosts[name] = rs.Interface().(string)
			}
		}
	}
}

func (t *vaultTask) initClient(conf *api.Config) error{
	client, err := api.NewClient(conf)
	if err!= nil {
		return err
	}
	client.SetToken(t.vaultToken)
	t.c = client.Logical()
	return nil
}

func (t *vaultTask) initTask(h *Handler){
	t.vaultAddr = h.VaultAddr
	t.vaultToken = h.VaultToken
	t.environment = h.Environment
}

func (t *vaultTask) initParentPath(){
	list, err := t.c.List("secret/metadata/dev")
	if err!=nil{
		log.Panicln(err)
	}
	var result []interface{}

	if list!=nil{
		for _, v := range list.Data{
			result = common.ConvertToArray(v)
		}
	}
	t.paths = result
}

func (t *vaultTask) initConfig() *api.Config {
	return &api.Config{
		Address: t.vaultAddr,
	}
}
