package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reverse-proxy/service"
)

const(
     forwardedHost = "X-Forwarded-Host"
     host = "host"
)

var (
	hosts map[string]string
)


func main() {
	h := iniHandler()
	hosts = h.InitHosts() //map[string]string{"/test":"http://127.0.0.1:8088"}

	http.HandleFunc("/", redirect)
	http.HandleFunc("/reload", reload)
	if err:=http.ListenAndServe(":9991", nil); err!= nil {
		 panic(err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	path := hosts[r.URL.Path]
	if len(path)==0{
		w.Write([]byte("Not found host on path " + r.URL.Path))
		return
	}
	u,_ := url.Parse(path)
	p := httputil.NewSingleHostReverseProxy(u)
	fmt.Printf("host %v\n", r.Header.Get(host))
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set(forwardedHost, r.Header.Get(host))
	r.Host = u.Host
	p.ServeHTTP(w, r)
}

func reload(w http.ResponseWriter, r *http.Request) {
	hosts = map[string]string{"/test":"http://127.0.0.1:8088", "/test2":"http://127.0.0.1:8088"}   //h.InitHosts()
	if _, err:=w.Write([]byte("Reload config successfully")); err != nil {
		log.Printf("Exception while generate resp %v\n", err.Error())
	}
}

func iniHandler() *service.Handler{
	return &service.Handler{
		VaultAddr: os.Getenv("VAULT_ADDR"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		Environment: os.Getenv("ENVIRONMENT"),
	}
}