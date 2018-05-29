package miyam

import (
	"io/ioutil"
	h "net/http"
	"regexp"
)

// There are 4 fundamental operations made by siteExtractor
// 0 - URL matching
// 1 - HTML downloading
// 2 - Extraction of information from HTML
// 3 - Decoration of the extracted information

// The extracted and decorated information is either a video title
// or a video URL

type siteExtractor struct {
	client     *h.Client
	urlMatcher *regexp.Regexp
	titleP     *htmlProc
	urlP       *htmlProc
}

func (s *siteExtractor) match(url string) (ok bool) {
	ok = s.urlMatcher.MatchString(url)
	return
}

func (s *siteExtractor) extract(url string) (v *videoInfo,
	e error) {
	v = new(videoInfo)
	var ss []string
	ss, e = s.titleP.proc(s.client, url)
	if e == nil {
		if len(ss) == 1 {
			v.title = ss[0]
		} else {
			e = NoTitleFound(url)
		}
	}
	v.urls, e = s.urlP.proc(s.client, url)
	return
}

type htmlProc struct {
	rs []*regexp.Regexp
}

func (p *htmlProc) proc(c *h.Client, url string) (us []string,
	e error) {
	var r *h.Response
	r, e = c.Get(url)
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(r.Body)
	}
	if e == nil {
		html := string(bs)
		i, ok := 0, false
		us = make([]string, 0)
		for !ok && i != len(p.rs) {
			ss := p.rs[i].FindStringSubmatch(html)
			ok = ss != nil
			if ok {

			}
		}

	}
	return
}

/* annie's youtubeDownload
func youtubeDownload(uri string) downloader.VideoData {
vid := utils.MatchOneOf(
		uri,
		`watch\?v=([^/&]+)`,
		`youtu\.be/([^?/]+)`,
		`embed/([^/?]+)`,
		`v/([^/?]+)`,
)
-- regexp on uri
-- vid != nil && len(vid) > 0
videoURL := fmt.Sprintf(
		"https://www.youtube.com/watch?v=%s&gl=US&hl=en&has_verified=1&bpctr=9999999999",
		vid[1],
	)
-- decoration

html := request.Get(videoURL, referer)
ytplayer := utils.MatchOneOf(html, `;ytplayer\.config\s*=\s*({.+?});`)[1]
-- HTTP request and regexp

var youtube youtubeData
	json.Unmarshal([]byte(ytplayer), &youtube)
	title := youtube.Args.Title
	streams := strings.Split(youtube.Args.Stream, ",")
-- regexp


*/

/*
youtubeData grammar

*/
