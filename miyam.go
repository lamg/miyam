package miyam

import (
	"fmt"
	"io"
	h "net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/spf13/afero"

	"github.com/cheggaaa/pb"
)

// Miyam program interface
type Miyam struct {
	y   *YouTube
	dwn *downloader
}

// NewMiyam returns a new Miyam instance
func NewMiyam(proxy string, timeout time.Duration) (m *Miyam) {
	m = new(Miyam)
	tr := &h.Transport{
		TLSHandshakeTimeout: timeout,
	}
	if proxy != "" {
		ur, _ := url.Parse(proxy)
		tr.Proxy = func(r *h.Request) (u *url.URL, e error) {
			u = ur
			return
		}
	}
	m.y = &YouTube{
		tokensCache: make(map[string][]string),
		Cl: &h.Client{
			Timeout:   timeout,
			Transport: tr,
		},
	}
	m.dwn = &downloader{
		client: m.y.Cl,
		fs:     afero.NewOsFs(),
	}
	return
}

// DownloadTerm downloads the url in the current working directory
// with a terminal interface
func (m *Miyam) DownloadTerm(ur, itag string) (e error) {
	var vd *VideoData
	vd, e = m.y.Get(ur)
	var fm FormatData
	if e == nil {
		var ok bool
		if itag == "" {
			fm = vd.Formats[DefaultKey]
		} else {
			fm, ok = vd.Formats[itag]
			if !ok {
				e = noItag(itag)
			}
		}
	}

	var dest io.WriteCloser
	var offset, total uint64
	if e == nil {
		ds := destFile(vd.Title, fm.Ext)
		dest, offset, e = m.dwn.storer(ds)
	}
	var src io.ReadCloser
	if e == nil {
		src, total, e = m.dwn.get(fm.URL, offset)
	}
	var bar *pb.ProgressBar
	if e == nil {
		bar = pb.New64(int64(total))
		bar.SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
		bar.Set64(int64(offset))
		bar.ShowSpeed = true
		bar.ShowFinalTime = true
		bar.SetMaxWidth(1000)
		bar.Start()
		m.dwn.barWr = bar
		e = m.dwn.copy(dest, src)
	}
	return
}

// InfoTerm shows video information in terminal
func (m *Miyam) InfoTerm(ur string) (e error) {
	var vd *VideoData
	vd, e = m.y.Get(ur)
	if e == nil {
		fmt.Printf("Site: %s\n", vd.Site)
		fmt.Printf("Title: %s\n", vd.Title)
		fmt.Println("Streams:")
		for k, v := range vd.Formats {
			fmt.Printf("%s: \n", k)
			fmt.Printf("\tQuality: %s\n", v.Quality)
			fmt.Printf("\tExtension: %s\n", v.Ext)
			sz := datasize.ByteSize(v.Size)
			fmt.Printf("\tSize: %s\n", sz.HumanReadable())
			fmt.Println()
		}
	}
	return
}

func noItag(itag string) (e error) {
	e = fmt.Errorf("No itag %s", itag)
	return
}

func destFile(title, ext string) (f string) {
	var repChs []string
	if runtime.GOOS == "windows" {
		repChs = []string{"\"", " ", "?", " ", "*", " ",
			"\\", " ", "<", " ", ">", " "}
	} else {
		repChs = []string{"/", " ", "|", "-", ": ", "：", ":",
			"："}
	}
	rep := strings.NewReplacer(repChs...)
	f = rep.Replace(title + "." + ext)
	return
}
