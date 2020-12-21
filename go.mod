module gsf

go 1.14

require (
	github.com/sirupsen/logrus v1.7.0
	gsf/common v0.0.0-00010101000000-000000000000 // indirect
	gsf/gsc v0.0.0-00010101000000-000000000000
	gsf/gsc/ffmpeg v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	gsf/common => ./common
	gsf/gsc => ./gsc
	gsf/gsc/ffmpeg => ./gsc/ffmpeg
)
