module github.com/cgoder/gss

go 1.14

require (
	github.com/cgoder/gsc v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro/v2 v2.9.1
	github.com/sirupsen/logrus v1.7.0
)

replace (
	github.com/cgoder/gsc => ../gsc
	github.com/cgoder/gss => ./
)