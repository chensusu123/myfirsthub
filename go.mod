module gitlab.ifreetalk.com/maze/maze_game_server

go 1.18

require (
	github.com/cloudwego/hertz v0.9.7
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gogo/protobuf v1.1.1
	github.com/gomodule/redigo v1.8.9
	github.com/gorilla/websocket v1.5.3
	github.com/hertz-contrib/websocket v0.2.0
	github.com/json-iterator/go v1.1.12
	github.com/lonng/nano v0.5.1
	github.com/xuri/excelize/v2 v2.9.0
	gitlab.ifreetalk.com/maze-plate/excel v0.0.0-20250515065852-efc0a876269a
	gitlab.ifreetalk.com/maze-plate/extra v1.0.1-0.20250401060651-722653168d37
	gitlab.ifreetalk.com/maze-plate/freetk v1.0.1-0.20250515074528-4334f63a4663
	gitlab.ifreetalk.com/maze-plate/io v0.0.0-20250411054516-f3778c3f9a74
	gitlab.ifreetalk.com/maze-plate/protodef v0.0.0-20250515072513-b2746e13cca6
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.24.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/bytedance/gopkg v0.1.0 // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/netpoll v0.6.4 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/lestrrat-go/strftime v1.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/nyaruka/phonenumbers v1.0.55 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/xuri/efp v0.0.0-20240408161823-9ad904a10d6d // indirect
	github.com/xuri/nfp v0.0.0-20240318013403-ab9948c2c4a7 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

//replace gitlab.ifreetalk.com/maze-plate/excel => gitlab.ifreetalk.com/maze-plate/excel v1.0.83-0.20250506053436-ef2c9ea111c0

//replace gitlab.ifreetalk.com/maze-plate/freetk v1.2.61 => /Users/majiange/data/dev/go_work/go_plate/src/gitlab.ifreetalk.com/maze-plate/freetk

replace github.com/lonng/nano => github.com/zsai001/nano-ex v0.0.0-20250513093902-993404574160
