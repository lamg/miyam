package miyam

import (
	"fmt"
	h "net/http"
	"regexp"
)

// VideoData data struct of video info
type VideoData struct {
	Site  string
	Title string
	Type  string
	// each format has it's own URLs and Quality
	Formats map[string]FormatData
}

const (
	// DefaultKey is the VideoData.Formats key associated to
	// the defaul FormatData
	DefaultKey = "default"
)

// FormatData data struct of every format
type FormatData struct {
	// [URLData: {URL, Size, Ext}, ...]
	// Some video files have multiple fragments
	// and support for downloading multiple image files at once
	URL     string
	Ext     string
	Quality string
	Size    uint64 // total size of all urls
}

// YouTube video information getter
type YouTube struct {
	Cl          *h.Client
	tokensCache map[string][]string
}

// Get gets the download information of a YouTube video
func (y *YouTube) Get(ur string) (d *VideoData, e error) {
	var vurl string
	vurl, e = decorate(ur)
	var html string
	if e == nil {
		html, e = page(vurl, y.Cl)
	}
	var ps *parsed
	if e == nil {
		ps, e = parse(html)
	}
	var fp map[string]fmPart
	if e == nil {
		fp, e = extract(ps)
	}
	if e == nil {
		d.Site = "YouTube"
		d.Title = ps.title
		d.Type = "video"
		d.Formats, e = fillFormatInfo(fp, y.tokensCache, y.Cl)
	}
	return
}

type parsed struct {
	title   string
	streams []string
	js      string
}

type fmPart struct {
	quality string
	ext     string
	itag    string
	sign    *signInfo
	// data structure preprocess genSigned
}

type signInfo struct {
	// shared between pre-execution and post-execution
	// of input/output
	valuesS string
	// input/output parameter
	tokensURL string
	// input/output result
	html string
	// return value
	signedURL string
}

type youtubeData struct {
	Args   args   `json:"args"`
	Assets assets `json:"assets"`
}

type args struct {
	Title  string `json:"title"`
	Stream string `json:"url_encoded_fmt_stream_map"`
	Audio  string `json:"adaptive_fmts"`
}

type assets struct {
	JS string `json:"js"`
}

// matchOneOf match one of the regular expressions in rs
func matchOneOf(text string, rs ...string) (r string, e error) {
	var sm []string
	for i := 0; len(sm) == 0 && i != len(rs); i++ {
		re := regexp.MustCompile(rs[i])
		sm = re.FindStringSubmatch(text)
	}
	if sm == nil || len(sm) < 2 {
		e = noMatch(text)
	} else {
		r = sm[1]
	}
	return
}

func noMatch(text string) (e error) {
	e = fmt.Errorf("No match for %s", text)
	return
}
