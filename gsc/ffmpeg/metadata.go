package ffmpeg

// Metadata ...
type Metadata struct {
	Format  Format    `json:"format"`
	Streams []Streams `json:"streams"`
}

// Format ...
type Format struct {
	Filename       string `json:"filename"`
	NbStreams      int    `json:"nb_streams"`
	NbPrograms     int    `json:"nb_programs"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	Duration       string `json:"duration"`
	StartTime      string `json:"start_time"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
	Tags           Tags   `json:"tags"`
}

// Streams ...
type Streams struct {
	Index              int         `json:"index"`
	ID                 string      `json:"id"`
	CodecName          string      `json:"codec_name"`
	CodecLongName      string      `json:"codec_long_name"`
	Profile            string      `json:"profile"`
	CodecType          string      `json:"codec_type"`
	CodecTimeBase      string      `json:"codec_time_base"`
	CodecTagString     string      `json:"codec_tag_string"`
	CodecTag           string      `json:"codec_tag"`
	Width              int         `json:"width"`
	Height             int         `json:"height"`
	CodedWidth         int         `json:"coded_width"`
	CodedHeight        int         `json:"coded_height"`
	HasBFrames         int         `json:"has_b_frames"`
	SampleAspectRatio  string      `json:"sample_aspect_ratio"`
	DisplayAspectRatio string      `json:"display_aspect_ratio"`
	PixFmt             string      `json:"pix_fmt"`
	Level              int         `json:"level"`
	ChromaLocation     string      `json:"chroma_location"`
	Refs               int         `json:"refs"`
	QuarterSample      string      `json:"quarter_sample"`
	DivxPacked         string      `json:"divx_packed"`
	RFrameRate         string      `json:"r_frame_rate"`
	AvgFrameRate       string      `json:"avg_frame_rate"`
	TimeBase           string      `json:"time_base"`
	DurationTs         int         `json:"duration_ts"`
	Duration           string      `json:"duration"`
	BitRate            string      `json:"bit_rate"`
	IsAVC              string      `json:"is_avc"`
	NalLengthSize      string      `json:"nal_length_size"`
	StartPts           int         `json:"start_pts"`
	StartTime          string      `json:"start_time"`
	BitsPerRawSample   string      `json:"bits_per_raw_sample"`
	NbFrames           string      `json:"nb_frames"`
	Disposition        Disposition `json:"disposition"`
	// Tags               Tags        `json:"tags"`
}

// Tags ...
type Tags struct {
	DisplayWidth  string `json:displayWidth`
	DisplayHeight string `json:displayHeight`
	Fps           string `json:"fps"`
	Encoder       string `json:"encoder"`
	Language      string `json:"language"`
	HandlerName   string `json:"handler_name"`
}

// Disposition ...
type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}

// GetFormat ...
func (m Metadata) GetFormat() Format {
	return m.Format
}

// GetStreams ...
func (m Metadata) GetStreams() (streams []Streams) {
	for _, element := range m.Streams {
		streams = append(streams, element)
	}
	return streams
}

// GetFilename ...
func (f Format) GetFilename() string {
	return f.Filename
}

// GetNbStreams ...
func (f Format) GetNbStreams() int {
	return f.NbStreams
}

// GetNbPrograms ...
func (f Format) GetNbPrograms() int {
	return f.NbPrograms
}

// GetFormatName ...
func (f Format) GetFormatName() string {
	return f.FormatName
}

// GetFormatLongName ...
func (f Format) GetFormatLongName() string {
	return f.FormatLongName
}

// GetDuration ...
func (f Format) GetDuration() string {
	return f.Duration
}

// GetSize ...
func (f Format) GetSize() string {
	return f.Size
}

// GetBitRate ...
func (f Format) GetBitRate() string {
	return f.BitRate
}

// GetProbeScore ...
func (f Format) GetProbeScore() int {
	return f.ProbeScore
}

// GetTags ...
func (f Format) GetTags() Tags {
	return f.Tags
}

// GetEncoder ...
func (t Tags) GetEncoder() string {
	return t.Encoder
}

//GetIndex ...
func (s Streams) GetIndex() int {
	return s.Index
}

//GetID ...
func (s Streams) GetID() string {
	return s.ID
}

//GetCodecName ...
func (s Streams) GetCodecName() string {
	return s.CodecName
}

//GetCodecLongName ...
func (s Streams) GetCodecLongName() string {
	return s.CodecLongName
}

//GetProfile ...
func (s Streams) GetProfile() string {
	return s.Profile
}

//GetCodecType ...
func (s Streams) GetCodecType() string {
	return s.CodecType
}

//GetCodecTimeBase ...
func (s Streams) GetCodecTimeBase() string {
	return s.CodecTimeBase
}

//GetCodecTagString ...
func (s Streams) GetCodecTagString() string {
	return s.CodecTagString
}

//GetCodecTag ...
func (s Streams) GetCodecTag() string {
	return s.CodecTag
}

//GetWidth ...
func (s Streams) GetWidth() int {
	return s.Width
}

//GetHeight ...
func (s Streams) GetHeight() int {
	return s.Height
}

//GetCodedWidth ...
func (s Streams) GetCodedWidth() int {
	return s.CodedWidth
}

//GetCodedHeight ...
func (s Streams) GetCodedHeight() int {
	return s.CodedHeight
}

//GetHasBFrames ...
func (s Streams) GetHasBFrames() int {
	return s.HasBFrames
}

//GetSampleAspectRatio ...
func (s Streams) GetSampleAspectRatio() string {
	return s.SampleAspectRatio
}

//GetDisplayAspectRatio ...
func (s Streams) GetDisplayAspectRatio() string {
	return s.DisplayAspectRatio
}

//GetPixFmt ...
func (s Streams) GetPixFmt() string {
	return s.PixFmt
}

//GetLevel ...
func (s Streams) GetLevel() int {
	return s.Level
}

//GetChromaLocation ...
func (s Streams) GetChromaLocation() string {
	return s.ChromaLocation
}

//GetRefs ...
func (s Streams) GetRefs() int {
	return s.Refs
}

//GetQuarterSample ...
func (s Streams) GetQuarterSample() string {
	return s.QuarterSample
}

//GetDivxPacked ...
func (s Streams) GetDivxPacked() string {
	return s.DivxPacked
}

//GetRFrameRrate ...
func (s Streams) GetRFrameRrate() string {
	return s.RFrameRate
}

//GetAvgFrameRate ...
func (s Streams) GetAvgFrameRate() string {
	return s.AvgFrameRate
}

//GetTimeBase ...
func (s Streams) GetTimeBase() string {
	return s.TimeBase
}

//GetDurationTs ...
func (s Streams) GetDurationTs() int {
	return s.DurationTs
}

//GetDuration ...
func (s Streams) GetDuration() string {
	return s.Duration
}

//GetDisposition ...
func (s Streams) GetDisposition() Disposition {
	return s.Disposition
}

//GetBitRate ...
func (s Streams) GetBitRate() string {
	return s.BitRate
}

//GetDefault ...
func (d Disposition) GetDefault() int {
	return d.Default
}

//GetDub ...
func (d Disposition) GetDub() int {
	return d.Dub
}

//GetOriginal ...
func (d Disposition) GetOriginal() int {
	return d.Original
}

//GetComment ...
func (d Disposition) GetComment() int {
	return d.Comment
}

//GetLyrics ...
func (d Disposition) GetLyrics() int {
	return d.Lyrics
}

//GetKaraoke ...
func (d Disposition) GetKaraoke() int {
	return d.Karaoke
}

//GetForced ...
func (d Disposition) GetForced() int {
	return d.Forced
}

//GetHearingImpaired ...
func (d Disposition) GetHearingImpaired() int {
	return d.HearingImpaired
}

//GetVisualImpaired ...
func (d Disposition) GetVisualImpaired() int {
	return d.VisualImpaired
}

//GetCleanEffects ...
func (d Disposition) GetCleanEffects() int {
	return d.CleanEffects
}

// GetMetadata Returns metadata for the specified input file
// func GetMetadata() (Metadata, error) {

// 	{
// 		var outb, errb bytes.Buffer

// 		input := t.input

// 		if t.inputPipeReader != nil {
// 			input = "pipe:"
// 		}

// 		args := []string{"-i", input, "-print_format", "json", "-show_format", "-show_streams", "-show_error"}

// 		cmd := exec.Command(ffprobe, args...)
// 		cmd.Stdout = &outb
// 		cmd.Stderr = &errb

// 		err := cmd.Run()
// 		if err != nil {
// 			return nil, fmt.Errorf("error executing (%s) with args (%s) | error: %s | message: %s %s", t.config.FfprobeBinPath, args, err, outb.String(), errb.String())
// 		}

// 		var metadata Metadata

// 		if err = json.Unmarshal([]byte(outb.String()), &metadata); err != nil {
// 			return nil, err
// 		}

// 		t.metadata = metadata

// 		return metadata, nil
// 	}

// 	return nil, errors.New("ffprobe binary not found")
// }
