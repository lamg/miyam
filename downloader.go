package miyam

import (
	"fmt"
	"io"
	h "net/http"

	"golang.org/x/tools/godoc/vfs"
)

type BarSeeker struct {
	Offset uint64
	Total  uint64
}

type downloader struct {
	barWr  io.Writer
	barSk  *BarSeeker
	client *h.Client
	fs     vfs.FileSystem
}

func (m *downloader) download(dest io.WriteCloser,
	src io.ReadCloser) (e error) {
	wr := io.MultiWriter(dest, m.barWr)
	_, e = io.Copy(wr, src)
	src.Close()
	dest.Close()
	return
}

func (m *downloader) get(url string,
	offset uint64) (src io.ReadCloser, total uint64, e error) {
	var rq *h.Request
	rq, e = h.NewRequest(h.MethodGet, url, nil)
	if e == nil && offset != 0 {
		rq.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}
	var r *h.Response
	r, e = m.client.Do(rq)
	if e == nil {
		src, total = r.Body, uint64(r.ContentLength)
	}
	return
}

func (m *downloader) storer(path string) (dest io.WriteCloser,
	offset uint64, e error) {
	fi, e := m.fs.Stat(path)
	if e == nil {
		offset = uint64(fi.Size())
		if offset != 0 {
			dest, e = m.fs.Open(path) //open for appending
		} else {
			dest, e = m.fs.Create(path)
		}
	}
	return
}

func (m *downloader) downloadURL(url, path string) (e error) {
	var dest io.WriteCloser
	var offset, total uint64
	dest, offset, e = m.storer(path)
	var src io.ReadCloser
	if e == nil {
		src, total, e = m.get(url, offset)
	}
	if e == nil {
		m.barSk.Offset, m.barSk.Total = offset, total
		e = m.download(dest, src)
	}
	return
}
