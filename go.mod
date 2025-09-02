module github.com/go-musicfox/netease-music

go 1.23.0

require (
	github.com/buger/jsonparser v1.1.1
	github.com/cnsilvan/UnblockNeteaseMusic v0.0.0-20240731043907-91afd9361e8b
	github.com/forgoer/openssl v1.8.0
	github.com/go-musicfox/requests v0.2.3
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
)

require (
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	golang.org/x/text v0.28.0 // indirect
)

replace github.com/cnsilvan/UnblockNeteaseMusic => github.com/go-musicfox/UnblockNeteaseMusic v0.1.6
