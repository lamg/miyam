package miyam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	h "net/http"
	"net/url"
	"strings"
)

func decorate(ur string) (vurl string, e error) {
	var id string
	id, e = matchOneOf(
		ur,
		`watch\?v=([^/&]+)`,
		`youtu\.be/([^?/]+)`,
		`embed/([^/?]+)`,
		`v/([^/?]+)`,
	)
	if e == nil {
		vurl = fmt.Sprintf(
			"https://www.youtube.com/watch?v=%s&gl=US&hl=en"+
				"&has_verified=1&bpctr=9999999999",
			id,
		)
	}
	return
}

func parse(html string) (d *parsed, e error) {
	var ytplayer string
	ytplayer, e = matchOneOf(html,
		`;ytplayer\.config\s*=\s*({.+?});`)
	var yd *youtubeData
	if e == nil {
		yd = new(youtubeData)
		e = json.Unmarshal([]byte(ytplayer), yd)
	}
	if e == nil {
		d = &parsed{
			title:   yd.Args.Title,
			streams: strings.Split(yd.Args.Stream, ","),
			js:      yd.Assets.JS,
		}
		if yd.Args.Audio != "" {
			audio := strings.Split(yd.Args.Audio, ",")
			d.streams = append(d.streams, audio...)
		}
	}
	return
}

func extract(ps *parsed) (fp map[string]fmPart, e error) {
	fp = make(map[string]fmPart)
	for i := 0; e == nil && i != len(ps.streams); i++ {
		var stream url.Values
		stream, e = url.ParseQuery(ps.streams[i])
		if e == nil {
			var fm fmPart
			fm.quality = stream.Get("quality")
			streamURL := stream.Get("url")
			fm.sign = newSignInfo(streamURL, ps.js, stream)
			fm.ext, e = matchOneOf(stream.Get("type"),
				`video/([[:word:]]+);`, `audio/([[:word:]]+)`)
			fm.itag = stream.Get("itag")
			fp[fm.itag] = fm
			if i == 0 {
				fp[DefaultKey] = fm
			}
		}
	}
	return
}

func newSignInfo(streamURL, assets string,
	v url.Values) (s *signInfo) {
	s = new(signInfo)
	valuesSig := v.Get("sig")
	s.tokensURL = fmt.Sprintf("https://www.youtube.com%s", assets)
	if strings.Contains(streamURL, "signature=") {
		s.signedURL = streamURL
	}
	if valuesSig == "" {
		s.valuesS = v.Get("s")
	} else {
		s.signedURL = fmt.Sprintf("%s&signature=%s", streamURL,
			valuesSig)
	}
	return
}

func (s *signInfo) decryptSign(tokensCache map[string][]string) {
	if s.valuesS != "" {
		tokens, ok := tokensCache[s.tokensURL]
		if !ok {
			tokens = signTokens(s.html)
			tokensCache[s.tokensURL] = tokens
		}
		s.signedURL = decipher(tokens, s.valuesS)
	}
}

func fillFormatInfo(fp map[string]fmPart, tkC map[string][]string,
	c *h.Client) (fm map[string]FormatData,
	e error) {
	fm = make(map[string]FormatData)
	for k, v := range fp {
		if v.sign.signedURL == "" {
			v.sign.html, e = page(v.sign.tokensURL, c)
		}
		var sz uint64
		if e == nil {
			v.sign.decryptSign(tkC)
			sz, e = size(v.sign.signedURL, c)
		}
		fm[k] = FormatData{
			Ext:     v.ext,
			Quality: v.quality,
			Size:    sz,
			URL:     v.sign.signedURL,
		}
	}
	return
}

const (
	referrer    = "https://youtube.com"
	headReferer = "Referer"
	userAgentV  = "Mozilla/5.0 (Linux; Android 7.0; SM-G930VC Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/58.0.3029.83 Mobile Safari/537.36"
	userAgentK  = "User-Agent"
)

func page(ur string, c *h.Client) (r string, e error) {
	var rq *h.Request
	rq, e = h.NewRequest(h.MethodGet, ur, nil)
	var res *h.Response
	if e == nil {
		rq.Header.Set(headReferer, referrer)
		rq.Header.Set(userAgentK, userAgentV)
		res, e = c.Do(rq)
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

func size(ur string, c *h.Client) (r uint64, e error) {
	var req *h.Request
	req, e = h.NewRequest(h.MethodGet, ur, nil)
	var res *h.Response
	if e == nil {
		req.Header.Set(headReferer, referrer)
		req.Header.Set(userAgentK, userAgentV)
		res, e = c.Do(req)
	}
	if e == nil {
		r = uint64(res.ContentLength)
	}
	return
}
