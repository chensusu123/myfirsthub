module maze_game_server

go 1.24.2

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/asynkron/protoactor-go v0.0.0-20240822202345-3c0e61ca19c9
	github.com/bwmarrin/snowflake v0.3.0
	github.com/cloudwego/hertz v0.10.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.5.4
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/gomodule/redigo v1.8.9
	github.com/google/uuid v1.6.0
	github.com/gorilla/schema v1.4.1
	github.com/gorilla/websocket v1.5.3
	github.com/hertz-contrib/websocket v0.2.0
	github.com/json-iterator/go v1.1.12
	github.com/pingcap/check v0.0.0-20211026125417-57bd13f7b5f0
	github.com/pingcap/errors v0.11.4
	github.com/polarismesh/polaris-go v1.6.1
	github.com/samber/slog-zap/v2 v2.6.2
	github.com/spf13/cast v1.9.2
	github.com/urfave/cli v1.22.16
	github.com/xuri/excelize/v2 v2.9.0
	gitlab.ifreetalk.com/maze-plate/freetk v1.0.1-0.20250624033723-994216aa64aa
	go.uber.org/atomic v1.11.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.36.5
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/gorm v1.30.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Workiva/go-datastructures v1.1.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/gopkg v0.1.2 // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/gopkg v0.1.4 // indirect
	github.com/cloudwego/netpoll v0.7.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250607225305-033d6d78b36a // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hertz-contrib/logger/zap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/lestrrat-go/strftime v1.0.1 // indirect
	github.com/lithammer/shortuuid/v4 v4.0.0 // indirect
	github.com/lmittmann/tint v1.0.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/nyaruka/phonenumbers v1.0.55 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.37.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/orcaman/concurrent-map v1.0.0 // indirect
	github.com/pingcap/log v0.0.0-20191012051959-b742a5d432e9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/polarismesh/specification v1.5.5-alpha.1 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/samber/lo v1.47.0 // indirect
	github.com/samber/slog-common v0.18.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	github.com/xuri/efp v0.0.0-20240408161823-9ad904a10d6d // indirect
	github.com/xuri/nfp v0.0.0-20240318013403-ab9948c2c4a7 // indirect
	gitlab.ifreetalk.com/maze-plate/extra v1.0.1-0.20250401060651-722653168d37 // indirect
	gitlab.ifreetalk.com/nano-ecosystem/fklog v0.0.0-20250609141356-5698d091b098 // indirect
	gitlab.ifreetalk.com/nano-ecosystem/nlog v0.0.0-20250621054907-e3e89b67b4be // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.22.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.2.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
)

//replace gitlab.ifreetalk.com/maze-plate/freetk => /Users/majiange/data/dev/go_work/maze/maze-plate/freetk
