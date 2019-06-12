package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"html/template"
	"net"
	"net/http"
	"os"
	"time"
)

type Data struct {
	Addr    string
	Time    float64
	Version []string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		domain := r.FormValue("domain")
		if domain != "" {
			part1 := Data{}

			// Pinger
			p := fastping.NewPinger()

			ra, err := net.ResolveIPAddr("ip4", domain)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			p.AddIPAddr(ra)
			p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {

				part1.Addr = addr.String()
				part1.Time = rtt.Seconds()
			}
			err = p.Run()
			//
			// request for php version
			var jsonStr = []byte(`{"":""}`)
			req, err := http.NewRequest("POST", "http://"+domain, bytes.NewBuffer(jsonStr))
			req.Header.Set("X-Custom-Header", "myvalue")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			part1.Version = resp.Header["X-Powered-By"]

			jsonData, err := json.Marshal(part1)
			tmpl, _ := template.ParseFiles("templates/data.html")
			_, _ = fmt.Fprintln(w, string(jsonData))
			_ = tmpl.Execute(w, string(jsonData))
		}
	}

}

func main() {

	http.HandleFunc("/", IndexHandler)

	_ = http.ListenAndServe(":8182", nil)
}
