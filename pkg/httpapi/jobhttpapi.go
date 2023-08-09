package httpapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"timerhodoks/pkg/raftstore"
)

type JobHttpAPI struct {
	Store *raftstore.RaftStore
}

// 操作任务的API
// 增删查改
func (h *JobHttpAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case http.MethodPost:
		v, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read on POST: (%v)\n", err)
			http.Error(w, "Failed on POST", http.StatusBadRequest)
			return
		}
		err = h.Store.ProposeUpdate(v)
		if err != nil {
			log.Printf("Failed to parse POST: (%v)\n", err)
			http.Error(w, "Failed on POST", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodPut:
		v, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read on PUT: (%v)\n", err)
			http.Error(w, "Failed on PUT", http.StatusBadRequest)
			return
		}
		err = h.Store.ProposeInsert(v)
		if err != nil {
			log.Printf("Failed to parse PUT: (%v)\n", err)
			http.Error(w, "Failed on PUT", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		v, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read on DELETE: (%v)\n", err)
			http.Error(w, "Failed on DELETE", http.StatusBadRequest)
			return
		}
		err = h.Store.ProposeDelete(v)
		if err != nil {
			log.Printf("Failed to parse DELETE: (%v)\n", err)
			http.Error(w, "Failed on DELETE", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		v, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read on Get: (%v)\n", err)
			http.Error(w, "Failed on Get", http.StatusBadRequest)
			return
		}

		httpJobGetBody := HttpJobGetBody{Id: 0}
		err = json.Unmarshal(v, &httpJobGetBody)

		if err != nil {
			log.Printf("Failed to parse Get: (%v)\n", err)
			http.Error(w, "Failed on Get", http.StatusBadRequest)
			return
		}
		if httpJobGetBody.Id == 0 {
			log.Printf("Invalid field: do not have Id or Id = 0\n")
			http.Error(w, "Failed on Get", http.StatusBadRequest)
			return
		}

		reply, ok := h.Store.Lookup(httpJobGetBody.Id)

		if ok {
			w.WriteHeader(http.StatusOK)
			reply, _ := json.Marshal(reply)
			w.Write(reply)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}

	default:
		w.Header().Set("Allow", http.MethodPut)
		w.Header().Add("Allow", http.MethodGet)
		w.Header().Add("Allow", http.MethodPost)
		w.Header().Add("Allow", http.MethodDelete)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
