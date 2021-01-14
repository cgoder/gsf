package core

type GscOptions struct {
	Input    string
	Output   string
	Opts     Options
	OptSlice []string
}

type Options struct {
	// Input                 string `json:"-i"`
	// Output                string `json:"-o"`
	Aspect                string `json:"-aspect"`
	Resolution            string `json:"-s"`
	VideoBitRate          string `json:"-b:v"`
	VideoBitRateTolerance string `json:"-bt"`
	VideoMaxBitRate       string `json:"-maxrate"`
	VideoMinBitrate       string `json:"-minrate"`
	VideoCodec            string `json:"-c:v"`
	Vframes               string `json:"-vframes"`
	FrameRate             string `json:"-r"`
	AudioRate             string `json:"-ar"`
	KeyframeInterval      string `json:"-g"`
	AudioCodec            string `json:"-c:a"`
	AudioBitrate          string `json:"-ab"`
	AudioChannels         string `json:"-ac"`
	AudioVariableBitrate  string `json:"-q:a"`
	BufferSize            string `json:"-bufsize"`
	Threadset             string `json:"-threads"`
	Threads               string `json:"-threads"`
	Preset                string `json:"-preset"`
	Tune                  string `json:"-tune"`
	AudioProfile          string `json:"-profile:a"`
	VideoProfile          string `json:"-profile:v"`
	Target                string `json:"-target"`
	Duration              string `json:"-t"`
	Qscale                string `json:"-qscale"`
	Crf                   string `json:"-crf"`
	Strict                string `json:"-strict"`
	MuxDelay              string `json:"-muxdelay"`
	SeekTime              string `json:"-ss"`
	SeekUsingTimestamp    string `json:"-seek_timestamp"`
	MovFlags              string `json:"-movflags"`
	HideBanner            string `json:"-hide_banner"`
	OutputFormat          string `json:"-f"`
	SegmentTime           string `json:"-segment_time"`
	ResetTimestamps       string `json:"-reset_timestamps"`
	CopyTs                string `json:"-copyts"`
	NativeFramerateInput  string `json:"-re"`
	InputInitialOffset    string `json:"-itsoffset"`
	RtmpLive              string `json:"-rtmp_live"`
	HlsPlaylistType       string `json:"-hls_playlist_type"`
	HlsListSize           string `json:"-hls_list_size"`
	HlsSegmentDuration    string `json:"-hls_time"`
	HlsMasterPlaylistName string `json:"-master_pl_name"`
	HlsSegmentFilename    string `json:"-hls_segment_filename"`
	HTTPMethod            string `json:"-method"`
	HTTPKeepAlive         string `json:"-multiple_requests"`
	Hwaccel               string `json:"-hwaccel"`
	VideoFilter           string `json:"-vf"`
	AudioFilter           string `json:"-af"`
	SkipVideo             string `json:"-vn"`
	SkipAudio             string `json:"-an"`
	Map                   string `json:"-map"`
	CompressionLevel      string `json:"-compression_level"`
	MapMetadata           string `json:"-map_metadata"`
	EncryptionKey         string `json:"-hls_key_info_file"`
	Bframe                string `json:"-bf"`
	PixFmt                string `json:"-pix_fmt"`
	Overwrite             string `json:"-y"`
}

type OptionFunc func(*GscOptions)

func SetInput(src string) OptionFunc {
	return func(o *GscOptions) {
		o.Input = src
	}
}

func SetOutput(dst string) OptionFunc {
	return func(o *GscOptions) {
		o.Output = dst
	}
}

// // transcode
// // para := []string{
// // 	"-vf", "scale=-2:960",
// // 	"-c:v", "libx264",
// // 	"-profile:v", "main",
// // 	"-level:v", "3.1",
// // 	"-x264opts", "scenecut=0:open_gop=0:min-keyint=72:keyint=72",
// // 	"-minrate", "1000k",
// // 	"-maxrate", "1000k",
// // 	"-bufsize", "1000k",
// // 	"-b:v", "1000k",
// // 	"-y",
// // 	"-i", srcPath + srcFile,
// // 	destPath + destFile,
// // }
// func NewOptionsEncode(opts ...OptionFunc) Options {
// 	opt := Options{
// 		VideoFilter:     "scale=-2:960",
// 		VideoCodec:      "libx264",
// 		VideoProfile:    "main",
// 		VideoMinBitrate: "1000k",
// 		VideoMaxBitRate: "1000k",
// 		BufferSize:      "1000k",
// 		VideoBitRate:    "1000k",
// 		Overwrite:       "",
// 	}

// 	for _, o := range opts {
// 		o(&opt)
// 	}

// 	return opt
// }

// // split
// // para := []string{
// // 	"-c", "copy",
// // 	"-f", "segment",
// // 	"-segment_time", "5",
// // 	"-reset_timestamps", "1",
// // 	"-map", "0:0",
// // 	"-map", "0:1",
// // 	"-y",
// // }
// // destFile = "%d.mp4"
// func NewOptionsSplit(opts ...OptionFunc) Options {
// 	opt := Options{
// 		VideoCodec:      "copy",
// 		AudioCodec:      "copy",
// 		OutputFormat:    "segment",
// 		SegmentTime:     "5",
// 		ResetTimestamps: "1",
// 		Map:             "0:0",
// 		// Map:             "0:1",
// 		Overwrite: "",
// 	}

// 	for _, o := range opts {
// 		o(&opt)
// 	}

// 	return opt
// }

// // remux
// // para := []string{
// // 	"-c", "copy",
// // }
// func NewOptionsRemux(opts ...OptionFunc) Options {
// 	opt := Options{
// 		VideoCodec: "copy",
// 		AudioCodec: "copy",
// 		Overwrite:  "",
// 	}

// 	for _, o := range opts {
// 		o(&opt)
// 	}

// 	return opt
// }
