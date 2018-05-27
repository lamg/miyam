package miyam

import (
	"fmt"
	"path"
)

type infoExtractor struct {
	exts        []*siteExtractor
	downloadDir string
}

type downloadInfo struct {
	url  string
	path string
}

type videoInfo struct {
	title string
	urls  []string
}

func NoTitleFound(url string) (e error) {
	e = fmt.Errorf("No title found for %s", url)
	return
}

func NoURLFound(url string) (e error) {
	e = fmt.Errorf("No URL found for %s", url)
	return
}

func (n *infoExtractor) extract(url string) (di []downloadInfo,
	e error) {
	ok, i := false, 0
	for !ok && i != len(n.exts) {
		ok, i = n.exts[i].match(url), i+1
	}
	var vi *videoInfo
	if ok {
		vi, e = n.exts[i].extract(url)
	} else {
		e = NoMatch(url)
	}
	if e == nil {
		di = make([]downloadInfo, len(vi.urls))
		for i, j := range vi.urls {
			di[i] = downloadInfo{
				path: path.Join(n.downloadDir, vi.title),
				url:  j,
			}
		}
	}
	return
}

func NoMatch(url string) (e error) {
	e = fmt.Errorf("No matcher found for %s", url)
	return
}
