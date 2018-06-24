package miyam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	h "net/http"
	"net/url"
	"regexp"
	"strings"
)

// VideoData data struct of video info
type VideoData struct {
	Site  string
	Title string
	Type  string
	// each format has it's own URLs and Quality
	Formats map[string]FormatData
}

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
	var id string
	id, e = matchOneOf(
		ur,
		`watch\?v=([^/&]+)`,
		`youtu\.be/([^?/]+)`,
		`embed/([^/?]+)`,
		`v/([^/?]+)`,
	)
	var html string
	if e == nil {
		vurl := fmt.Sprintf(
			"https://www.youtube.com/watch?v=%s&gl=US&hl=en"+
				"&has_verified=1&bpctr=9999999999",
			id,
		)
		html, e = y.page(vurl)
	}
	var ytplayer string
	if e == nil {
		ytplayer, e = matchOneOf(html,
			`;ytplayer\.config\s*=\s*({.+?});`)
	}
	var yd *youtubeData
	if e == nil {
		yd = new(youtubeData)
		e = json.Unmarshal([]byte(ytplayer), yd)
	}
	if e == nil {
		d = &VideoData{
			Site:  "YouTube",
			Title: yd.Args.Title,
		}
		d.Formats, e = y.extractURLs(yd, ur)
	}
	return
}

func (y *YouTube) extractURLs(yd *youtubeData,
	ur string) (fs map[string]FormatData, e error) {
	// TODO extract audio URLs
	streams := append(strings.Split(yd.Args.Stream, ","),
		strings.Split(yd.Args.Audio, ",")...)

	fs = make(map[string]FormatData)
	for i := 0; e == nil && i != len(streams); i++ {
		var stream url.Values
		stream, e = url.ParseQuery(streams[i])
		var fm FormatData
		if e == nil {
			fm.Quality = stream.Get("quality")
			fm.Ext, e = matchOneOf(stream.Get("type"),
				`video/([[:word:]]+);`, `audio/([[:word:]]+)`)
		}
		if e == nil {
			streamURL := stream.Get("url")
			fm.URL, e = y.genSignedURL(streamURL, yd.Assets.JS, stream)
		}
		if e == nil {
			fm.Size, e = y.size(fm.URL)
		}
		if e == nil {
			itag := stream.Get("itag")
			fs[itag] = fm
		}
	}
	return
}

const (
	referrer      = "https://youtube.com"
	headReferer   = "Referer"
	contentLength = "Content-Length"
)

func (y *YouTube) size(ur string) (r uint64, e error) {
	var req *h.Request
	req, e = h.NewRequest(h.MethodHead, ur, nil)
	var res *h.Response
	if e == nil {
		req.Header.Set(headReferer, referrer)
		res, e = y.Cl.Do(req)
	}
	if e == nil {
		r = uint64(res.ContentLength)
	}
	return
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

func (y *YouTube) page(ur string) (r string, e error) {
	var rq *h.Request
	rq, e = h.NewRequest(h.MethodGet, ur, nil)
	var res *h.Response
	if e == nil {
		rq.Header.Set(headReferer, referrer)
		res, e = y.Cl.Do(rq)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(res.Body)
		res.Body.Close()
	}
	if e == nil {
		r = string(bs)
	}
	return
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
