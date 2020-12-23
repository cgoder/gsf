package ffmpeg

import (
	"encoding/json"

	"github.com/cgoder/gsc/common"
	log "github.com/sirupsen/logrus"
)

// type FFopt struct {
// 	// input file
// 	Input string
// 	// output file
// 	Output string
// 	// ffmpeg cmmmon option
// 	CmdOpt string

// 	Opt Options
// }

type Options struct {
	// input file
	// Input string
	// // output file
	// Output string
	// ffmpeg cmmmon option
	CmdSlice []string
	// all parsed args from CmdOpt
	CmdOpt    string
	arguments map[string]string
}

// Options defines allowed FFmpeg arguments
type CmdArgs struct {
	Input                 string `json:"-i"`
	Output                string `json:"-o"`
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
	CompressionLevel      string `json:"-compression_level"`
	MapMetadata           string `json:"-map_metadata"`
	EncryptionKey         string `json:"-hls_key_info_file"`
	Bframe                string `json:"-bf"`
	PixFmt                string `json:"-pix_fmt"`
	Overwrite             string `json:"-y"`
}

// type CmdArgs struct {
// 	Input                 string `json:"-i"`
// 	Output                string `json:"-o,omitempty"`
// 	Aspect                string `json:"-aspect,omitempty"`
// 	Resolution            string `json:"-s,omitempty"`
// 	VideoBitRate          string `json:"-b:v,omitempty"`
// 	VideoBitRateTolerance string `json:"-bt,omitempty"`
// 	VideoMaxBitRate       string `json:"-maxrate,omitempty"`
// 	VideoMinBitrate       string `json:"-minrate,omitempty"`
// 	VideoCodec            string `json:"-c:v,omitempty"`
// 	Vframes               string `json:"-vframes,omitempty"`
// 	FrameRate             string `json:"-r,omitempty"`
// 	AudioRate             string `json:"-ar,omitempty"`
// 	KeyframeInterval      string `json:"-g,omitempty"`
// 	AudioCodec            string `json:"-c:a,omitempty"`
// 	AudioBitrate          string `json:"-ab,omitempty"`
// 	AudioChannels         string `json:"-ac,omitempty"`
// 	AudioVariableBitrate  string `json:"-q:a,omitempty"`
// 	BufferSize            string `json:"-bufsize,omitempty"`
// 	Threadset             string `json:"-threads,omitempty"`
// 	Threads               string `json:"-threads"`
// 	Preset                string `json:"-preset,omitempty"`
// 	Tune                  string `json:"-tune,omitempty"`
// 	AudioProfile          string `json:"-profile:a,omitempty"`
// 	VideoProfile          string `json:"-profile:v,omitempty"`
// 	Target                string `json:"-target,omitempty"`
// 	Duration              string `json:"-t,omitempty"`
// 	Qscale                string `json:"-qscale,omitempty"`
// 	Crf                   string `json:"-crf,omitempty"`
// 	Strict                string `json:"-strict,omitempty"`
// 	MuxDelay              string `json:"-muxdelay,omitempty"`
// 	SeekTime              string `json:"-ss,omitempty"`
// 	SeekUsingTimestamp    string `json:"-seek_timestamp,omitempty"`
// 	MovFlags              string `json:"-movflags,omitempty"`
// 	HideBanner            string `json:"-hide_banner,omitempty"`
// 	OutputFormat          string `json:"-f,omitempty"`
// 	CopyTs                string `json:"-copyts,omitempty"`
// 	NativeFramerateInput  string `json:"-re,omitempty"`
// 	InputInitialOffset    string `json:"-itsoffset,omitempty"`
// 	RtmpLive              string `json:"-rtmp_live,omitempty"`
// 	HlsPlaylistType       string `json:"-hls_playlist_type,omitempty"`
// 	HlsListSize           string `json:"-hls_list_size,omitempty"`
// 	HlsSegmentDuration    string `json:"-hls_time,omitempty"`
// 	HlsMasterPlaylistName string `json:"-master_pl_name,omitempty"`
// 	HlsSegmentFilename    string `json:"-hls_segment_filename,omitempty"`
// 	HTTPMethod            string `json:"-method,omitempty"`
// 	HTTPKeepAlive         string `json:"-multiple_requests,omitempty"`
// 	Hwaccel               string `json:"-hwaccel,omitempty"`
// 	VideoFilter           string `json:"-vf,omitempty"`
// 	AudioFilter           string `json:"-af,omitempty"`
// 	SkipVideo             string `json:"-vn,omitempty"`
// 	SkipAudio             string `json:"-an,omitempty"`
// 	CompressionLevel      string `json:"-compression_level,omitempty"`
// 	MapMetadata           string `json:"-map_metadata,omitempty"`
// 	EncryptionKey         string `json:"-hls_key_info_file,omitempty"`
// 	Bframe                string `json:"-bf,omitempty"`
// 	PixFmt                string `json:"-pix_fmt,omitempty"`
// 	Overwrite             string `json:"-y"`
// 	// StreamIds             map[string]string `json:"-streamid,omitempty"`
// 	// Metadata              map[string]string `json:"-metadata,omitempty"`
// 	// WhiteListProtocols []string `json:"-protocol_whitelist,omitempty"`
// 	// reuse for output
// 	// ExtraArgs map[string]interface{} `json:",omitempty"`
// }

// GetArgument
func (opts *Options) GetArgument(arg string) (string, bool) {
	if opts.CmdSlice == nil {
		return "", false
	}
	// cmd := strings.Join(opts.CmdSlice, " ")
	// log.Info(common.JsonFormat(cmd))
	if opts.arguments == nil {
		if err := json.Unmarshal([]byte(opts.CmdOpt), &opts.arguments); err != nil {
			log.Error("cmd args encode to map fail")
			return "", false
		}

		// if opts.CmdOpt == "" {
		// 	log.Error("parse args first!")
		// 	return "", false
		// } else {
		// 	opts.SetStrArguments(opts.CmdSlice)
		// }
	}
	log.Info(common.JsonFormat(opts.arguments))

	if v, ok := opts.arguments[arg]; ok {
		return v, true
	}

	log.Error("can not find arg-> .", arg)
	return "", false
}

// GetStrArguments ...
func (opts *Options) GetStrArguments() map[string]string {
	return opts.arguments
}

// ParseStrArguments
func (opts *Options) SetStrArguments(cmdSlice []byte) {

	// if opts.CmdOpt == "" {
	// 	opts.CmdOpt = cmd
	// 	// opts.arguments = getCmdArguments(opts.CmdOpt)
	// } else {
	// 	if opts.CmdOpt != cmd {
	// 		opts.CmdOpt = cmd
	// 		// opts.arguments = getCmdArguments(opts.CmdOpt)
	// 	}
	// }

	// if err := json.Unmarshal([]byte(cmd), &opts.arguments); err != nil {
	// 	log.Error("cmd encode to map fail")

	// }

	// log.Info("opts---> ", opts)
	// log.Info("opts.CmdOpt---> ", JsonFormat(opts.CmdOpt))
	// log.Info("opts.arguments--->", JsonFormat(opts.arguments))
}

// func getCmdArguments(cmd string) map[string]string {
// 	var mapResult map[string]string
// 	if err := json.Unmarshal([]byte(cmd), &mapResult); err != nil {
// 		log.Error("cmd encode to map fail")
// 		return nil
// 	}
// 	log.Info("mapResult---> ", JsonFormat(mapResult))

// 	return mapResult
// }

// func getStrArguments(args CmdArgs) map[string]string {
// 	f := reflect.TypeOf(args)
// 	v := reflect.ValueOf(args)

// 	strArgs := make(map[string]string)

// 	for i := 0; i < f.NumField(); i++ {
// 		flag := f.Field(i).Tag.Get("json")
// 		// flag := f.Field(i).Interface()
// 		value := v.Field(i).Interface()
// 		// log.Info("flag.... ", flag)
// 		// log.Info("value.... ", value)
// 		if !v.Field(i).IsZero() {
// 			if vs, ok := value.(string); ok {
// 				strArgs[flag] = vs
// 			}
// 		}
// 	}
// 	// log.Info("strArgs.... ", JsonFormat(strArgs))
// 	return strArgs
// }
