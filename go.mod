module adblock

go 1.15

replace local/github.com/victorlpgazolli/lightdns => ../../opensource/lightdns

require (
	github.com/miekg/dns v1.1.35
	github.com/openmohan/lightdns v0.0.0-20181005121551-25aa6453d4ed
	local/github.com/victorlpgazolli/lightdns v0.0.0-00010101000000-000000000000
)
