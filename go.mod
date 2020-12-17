module gsf

go 1.14

require (
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1
	gsf/cmd v0.0.0-00010101000000-000000000000 // indirect
	gsf/ffmpeg v0.0.0-00010101000000-000000000000
)

replace (
	gsf/cmd => ./cmd
	gsf/ffmpeg => ./ffmpeg
)
