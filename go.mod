module timerhodoks

go 1.20

replace go.etcd.io/etcd/server/v3 v3.5.9 => ../go.etcd.io/etcd/server

replace go.etcd.io/etcd/api/v3 v3.6.0-alpha.0 => ../go.etcd.io/etcd/api

replace go.etcd.io/etcd/client/pkg/v3 v3.6.0-alpha.0 => ../go.etcd.io/etcd/client/pkg

require (
	go.etcd.io/etcd/client/pkg/v3 v3.6.0-alpha.0
	go.uber.org/zap v1.24.0
	google.golang.org/grpc v1.57.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.43.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/etcd/api/v3 v3.6.0-alpha.0 // indirect
	go.etcd.io/etcd/pkg/v3 v3.6.0-alpha.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20230526161137-0005af68ea54 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230525234035-dd9d682886f9 // indirect
)

require (
	github.com/bits-and-blooms/bitset v1.8.0
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/xhit/go-str2duration/v2 v2.1.0
	go.etcd.io/etcd/server/v3 v3.5.9
	go.etcd.io/raft/v3 v3.0.0-20230725081940-e87bcc2c5b0f
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sync v0.3.0
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.31.0
)

replace proto => ./etcd/timerhodoks/proto

// replace go.etcd.io/etcd/server/v3/etcdserver/api/snap => ../go.etcd.io/etcd/server/etcdserver/api/snap
