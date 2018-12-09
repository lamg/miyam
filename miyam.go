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
type extractor func(*netConf, string) (*info, error)

// (nc, url, offset) (rc, e)
type netRC func(*netConf, string, uint64) (io.ReadCloser, error)

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

func (m *miyam) dwInfo(url string) (inf *info, e error) {
	for i := 0; e == nil && inf == nil && i != len(m.alts); i++ {
		inf, e = m.alts[i](m.nc, url)
	}
	return
}

func (m *miyam) dwnStreams(inf *info) (fs []string, e error) {
	if !m.dc.skip {
		fs = make([]string, len(m.dc.formats))
		for i := 0; e == nil && i != len(m.dc.formats); i++ {
			str := inf.formats[m.dc.formats[i]]
			fs[i] = str.file
			var offset uint64
			var wc io.WriteCloser
			offset, wc, e = m.fwc(inf.matcher, &str)
			var rc io.ReadCloser
			if e == nil {
				rc, e = m.nrc(m.nc, str.url, offset)
			}
			if e == nil {
				_, e = io.Copy(wc, rc)
			}
		}
	}
	return
}

func (m *miyam) mux(paths []string) (e error) {
	// TODO search Go library for merging the streams
	// in a Matroska file
	return
}
