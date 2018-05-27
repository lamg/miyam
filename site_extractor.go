package miyam

import (
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

func (s *siteExtractor) extract(url string) (v *videoInfo, e error) {

	var fv []string
	fv, e = s.infoF.filter(url)
	if e == nil {
		v = new(videoInfo)
		if len(fv) == 0 {
			e = NoTitleFound(url)
		} else if len(fv) == 1 {
			e = NoURLFound(url)
		} else {
			v.title, v.urls = fv[0], fv[1:]
		}
	}
	return
}

func (s *siteExtractor) proc()
