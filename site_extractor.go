package miyam

import (
	"fmt"
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
	prs []*dataProc
}

func (s *siteExtractor) extract(url string,
	c *h.Client) (v *videoInfo, e error) {
	s.prs[0].inpTxt = url

	return
}

// YouTube exctractor
// video ID, from matched URL
// decorated ID as URL
// youtube player, from matched content, got from decorated
//  URL
// title matched from youtube player
// streams matched from youtube player

type dataProc struct {
	inpTxt string //textual data
	inpURL string //for getting the data from the web if inpTxt
	// is ""
	rgs       []*regexp.Regexp
	decorator string
	// decorator for matched text
}

// cl: client for getting inpURL content
func (d *dataProc) proc(cl *h.Client) (r []string, e error) {
	if d.inpTxt == "" {
	} else {
		var p *h.Response
		p, e = cl.Get(d.inpURL)
		var bs []byte
		if e == nil {
			bs, e = ioutil.ReadAll(p.Body)
		}
		if e == nil {
			d.inpTxt = string(bs)
		}
	}
	var sm []string // submatch selected result
	if e == nil {
		var sms []string // string submatch results
		for i := 0; sms == nil && i != len(d.rgs); i++ {
			sms = d.rgs[i].FindStringSubmatch(d.inpTxt)
		}
		if len(sms) > 1 {
			sm = sms[1]
			// sm is the content of the first submatch group
		} else {
			e = NoMatch(d.inpTxt)
		}
	}
	if e == nil {
		if d.decorator != "" {
			r = fmt.Sprintf(d.decorator, sm...)
		} else {
			r = sm
		}
	}
	return
}
