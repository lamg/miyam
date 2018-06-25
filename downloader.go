package miyam

import (
	"fmt"
	"io"
	h "net/http"
	"os"

	"github.com/spf13/afero"
)

type downloader struct {
	barWr io.Writer
	// to be taken by UI for setting progress bar status
	client *h.Client
	fs     afero.Fs
}

func (m *downloader) copy(dest io.WriteCloser,
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
			// open for appending
			dest, e = m.fs.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		}
	} else {
		dest, e = m.fs.Create(path)
	}
	return
}
