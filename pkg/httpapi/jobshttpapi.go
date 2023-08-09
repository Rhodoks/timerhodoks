package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"timerhodoks/pkg/job"
	"timerhodoks/pkg/raftstore"
)

type JobsHttpAPI struct {
	Store *raftstore.RaftStore
}

// layui表格的默认json格式
type JobTableResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Count   int             `json:"count"`
	Data    []*job.JobEntry `json:"data"`
}

// 返回表单指定页的条目
func (s *JobsHttpAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	page, err := strconv.Atoi(r.URL.Query()["page"][0])
	limit, err := strconv.Atoi(r.URL.Query()["limit"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.Store.Mutex.RLock()
	defer s.Store.Mutex.RUnlock()
	start := page*limit - limit
	end := page*limit - 1
	response := JobTableResponse{}
	response.Data = s.Store.Jobs.GetJobList(start, end)
	response.Count = s.Store.Jobs.Size()
	res, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
