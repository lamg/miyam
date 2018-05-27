package miyam

import (
	"fmt"
	"io"
	h "net/http"
)

type Miyam struct {
	barWr  io.Writer
	client *h.Client
	// filesystem interface
}

func (m *Miyam) download(dest io.WriteCloser,
	src io.ReadCloser) (e error) {
	wr := io.MultiWriter(dest, m.barWr)
	_, e = io.Copy(wr, src)
	src.Close()
	dest.Close()
	return
}

func (m *Miyam) get(url string, offset uint64) (src io.ReadCloser,
	e error) {
	var rq *h.Request
	rq, e = h.NewRequest(h.MethodGet, url, nil)
	if e == nil && offset != 0 {
		rq.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}
	var r *h.Response
	r, e = m.client.Do(rq)
	if e == nil {
		src = r.Body
	}
	return
}

func (m *Miyam) storer(path string) (dest io.WriteCloser,
	offset uint64, e error) {
	// filesystem interface
	// open temporary file for appending or creating
	return
}

func (m *Miyam) downloadURL(url, path string) (e error) {
	var dest io.WriteCloser
	var offset uint64
	dest, offset, e = m.storer(path)
	var src io.ReadCloser
	if e == nil {
		src, e = m.get(url, offset)
	}
	if e == nil {
		e = m.download(dest, src)
	}
	return
}
