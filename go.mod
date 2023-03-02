module dknight/go-todoapp

go 1.19

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dknight/go-todoapp-sandbox/controllers v0.0.0-00010101000000-000000000000
	github.com/dknight/go-todoapp-sandbox/models v0.0.0-20230302165719-08316196e434
	github.com/gofiber/fiber/v2 v2.42.0
	github.com/gofiber/template v1.7.5
	github.com/lib/pq v1.10.7
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/blockloop/scan v1.3.0 // indirect
	github.com/gofiber/fiber v1.14.6 // indirect
	github.com/gofiber/utils v0.0.10 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/schema v1.1.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.44.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
)

replace (
	github.com/dknight/go-todoapp-sandbox/controllers => ./controllers
	github.com/dknight/go-todoapp-sandbox/models => ./models
)
