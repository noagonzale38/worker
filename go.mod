module github.com/TicketsBot-cloud/worker

go 1.23.8

toolchain go1.24.2

// replace github.com/TicketsBot-cloud/database => ../database
// replace github.com/TicketsBot-cloud/gdl => ../gdl
// replace github.com/TicketsBot-cloud/archiverclient => ../archiverclient
// replace github.com/TicketsBot-cloud/logarchiver => ../logarchiver

require (
	cloud.google.com/go/profiler v0.4.2
	github.com/TicketsBot-cloud/analytics-client v0.0.0-20250604180646-6606dfc8fc8c
	github.com/TicketsBot-cloud/archiverclient v0.0.0-20250514201416-cf23f65eb3fc
	github.com/TicketsBot-cloud/common v0.0.0-20250509064208-a2d357175463
	github.com/TicketsBot-cloud/database v0.0.0-20250601094419-c168a11bd60a
	github.com/rxdn/gdl latest
	github.com/caarlos0/env/v10 v10.0.0
	github.com/elliotchance/orderedmap v1.8.0
	github.com/getsentry/sentry-go v0.32.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-redsync/redsync/v4 v4.13.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v4 v4.18.3
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/joho/godotenv v1.5.1
	github.com/json-iterator/go v1.1.12
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.22.0
	github.com/schollz/progressbar/v3 v3.18.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.10.0
	go.uber.org/atomic v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6
	golang.org/x/net v0.40.0
	golang.org/x/sync v0.14.0
	golang.org/x/tools v0.33.0
	gopkg.in/alexcesaro/statsd.v2 v2.0.0
)

require (
	cloud.google.com/go v0.121.0 // indirect
	cloud.google.com/go/auth v0.16.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/ClickHouse/ch-go v0.66.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.36.0 // indirect
	github.com/TicketsBot-cloud/logarchiver v0.0.0-20250514201320-d5141071a6eb // indirect
	github.com/TicketsBot/common v0.0.0-20241117150316-ff54c97b45c1 // indirect
	github.com/TicketsBot/ttlcache v1.6.1-0.20200405150101-acc18e37b261 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/caarlos0/env v3.5.0+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/pprof v0.0.0-20250501235452-c0086092b71a // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/juju/ratelimit v1.0.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/minio/crc64nvme v1.0.2 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.91 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/panjf2000/ants/v2 v2.11.3 // indirect
	github.com/pasztorpisti/qs v0.0.0-20171216220353-8d6c33ee906c // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.63.0 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/tatsuworks/czlib v0.0.0-20190916144400-8a51758ea0d9 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/twmb/franz-go v1.19.0 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.11.2 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.60.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.17.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/api v0.232.0 // indirect
	google.golang.org/genproto v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/grpc v1.72.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	nhooyr.io/websocket v1.8.17 // indirect
)
