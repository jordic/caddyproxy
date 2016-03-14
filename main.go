package main

import (
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/fsouza/go-dockerclient"
)

var (
	endpoint = flag.String("socket", "unix:///var/run/docker.sock", "Docker ro socket")
	notls    = flag.Bool("notls", false, "Disable caddy tls")
)

// Image is a running container
type Image struct {
	docker.APIContainers
}

func (i *Image) Domain() string {
	return i.Labels["tempo_domain"]
}

func (i *Image) Vars() map[string]interface{} {

	out := make(map[string]interface{})
	out["Name"] = i.Labels["tempo_proxyto"]
	out["Domains"] = i.Labels["tempo_domain"]
	if val, ok := i.Labels["tempo_root"]; ok {
		out["Root"] = val
	} else {
		out["Root"] = "/serve"
	}

	if val, ok := i.Labels["tempo_statics"]; ok {
		out["Statics"] = val
	}
	if *notls == true {
		out["Tls"] = "tls off"
	}

	return out
}

func main() {

	flag.Parse()

	client, _ := docker.NewClient(*endpoint)

	opts := docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": []string{"tempo_proxy=true"},
		},
	}

	res, err := client.ListContainers(opts)
	if err != nil {
		log.Fatal(err)
	}

	var images []*Image

	for _, item := range res {
		i := &Image{item}
		images = append(images, i)
	}

	//f, err := os.Open(*output)
	//w := bufio.NewWriter(f)
	//defer f.Close()
	//if err != nil {
	//	log.Fatalf("Error %s", err)
	//}
	w := os.Stdout

	tpl, err := template.New("conf").Parse(caddyTpl)
	if err != nil {
		log.Fatalf("tpl error: %s", err)
	}

	for _, it := range images {
		out := it.Vars()
		err := tpl.Execute(w, out)
		if err != nil {
			log.Printf("Err %s", err)
		}
		//w.Flush()
	}

}

var caddyTpl = `
{{.Domains}} {
    {{ if .Tls }}{{.Tls}}{{ end }}
    proxy / {{ .Name }} {
                proxy_header Host {host}
                proxy_header X-Forwarded-Proto {scheme}
                {{ if .Statics }}except {{.Statics}}{{end}}
    }

   root {{ .Root }}
}

`
