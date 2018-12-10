package main

import (
	"io"
	"net/url"
	"time"
)

type miyam struct {
	alts []extractor
	// ae is the *exCfg parameter for each element in alts
	nc  *netConf
	dc  *dwCfg
	nrc netRC
	fwc fileWC
}

// (nc, url)
type cfgExtr func(*netConf, string) (*info, error)

type extractor func(string) (*info, error)

func clsExtr(nc *netConf, ex cfgExtr) (r extractor) {
	r = func(url string) (i *info, e error) {
		i, e = ex(nc, url)
		return
	}
	return
}

// (nc, url, offset) (rc, e)
type cfgNetRC func(*netConf, string, uint64) (io.ReadCloser, error)

type netRC func(string, uint64) (io.ReadCloser, error)

func clsNetRC(nc *netConf, nrc cfgNetRC) (r netRC) {
	r = func(url string, offset uint64) (rc io.ReadCloser, e error) {
		rc, e = nrc(nc, url, offset)
		return
	}
	return
}

// (matcher, s) (offset, w, e)
type fileWC func(string, *stream) (uint64, io.WriteCloser, error)

type netConf struct {
	timeout    time.Duration
	httpProxy  *url.URL
	socksProxy string
}

type info struct {
	matcher string
	formats map[string]stream
}

type stream struct {
	url   string
	total uint64
	descr string
	ext   string
	file  string
}

type dwCfg struct {
	skip    bool
	formats []string
}

func dwInfo(alts []extractor, url string) (inf *info, e error) {
	for i := 0; e == nil && inf == nil && i != len(alts); i++ {
		inf, e = alts[i](url)
	}
	return
}

func dwnStreams(dc *dwCfg, inf *info, fwc fileWC,
	nrc netRC) (fs []string, e error) {
	if !dc.skip {
		fs = make([]string, len(dc.formats))
		for i := 0; e == nil && i != len(dc.formats); i++ {
			str := inf.formats[dc.formats[i]]
			fs[i] = str.file
			var offset uint64
			var wc io.WriteCloser
			offset, wc, e = fwc(inf.matcher, &str)
			var rc io.ReadCloser
			if e == nil {
				rc, e = nrc(str.url, offset)
			}
			if e == nil {
				_, e = io.Copy(wc, rc)
			}
		}
	}
	return
}

func mux(paths []string) (e error) {
	// TODO search Go library for merging the streams
	// in a Matroska file
	return
}
