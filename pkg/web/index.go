package web

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
	"timerhodoks/pkg/raftnode"
)

type IndexServer struct {
	RaftNode *raftnode.RaftNode
}

// web服务，使用模板替换服务器端时间，状态等
func (s *IndexServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tmpl, _ := template.ParseFiles("./web/index.html")
	info := make(map[string]string)

	serverTime := time.Now()
	info["Time"] = serverTime.Format("2006-01-02 15:04:05")
	info["TimeStamp"] = strconv.Itoa(int(serverTime.Unix()))
	info["Status"] = s.RaftNode.GetState().String()

	tmpl.Execute(w, info)
}
