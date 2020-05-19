package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reverse-proxy/service"
	"strings"
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
	hosts = h.InitHosts()

	http.HandleFunc("/", redirect)
	http.HandleFunc("/reload", reload)

	log.Printf("start on port %v\n", os.Getenv("PORT"))
	if err:=http.ListenAndServe(os.Getenv("PORT"), nil); err!= nil {
		 panic(err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	path := getHost(r.URL.Path[1:]) //hosts[r.URL.Path[1:]]
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
	r.RequestURI = "/v3/api-docs"
	r.URL.Path = "/v3/api-docs"
	fmt.Printf("host is ======= %v\n",r.RequestURI)
	p.ServeHTTP(w, r)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getHost(path string) string{
	for k,v := range hosts{
		if strings.Contains(path,k){
			return v
		}
	}
	return ""
}

func reload(w http.ResponseWriter, r *http.Request) {
	h := iniHandler()
	hosts = h.InitHosts()
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