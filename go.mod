module github.com/mahmoud-shabban/greenlight

go 1.24.5

replace github.com/mahmoud-shabban/greenlight => .

require (
	github.com/felixge/httpsnoop v1.0.4
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.10.9
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	golang.org/x/crypto v0.37.0
	golang.org/x/time v0.14.0
	gopkg.in/mail.v2 v2.3.1
)

require gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
