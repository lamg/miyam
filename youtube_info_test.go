package miyam

import (
	"io/ioutil"
	h "net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInfo(t *testing.T) {
	bs, e := ioutil.ReadFile("youtube_player_line.json")
	require.NoError(t, e)
	testHTML := string(bs)
	var ps *parsed
	ps, e = parse(testHTML)
	require.NoError(t, e)
	require.Equal(t, ps.title, "Eliécer Ávila   El régimen no puede"+
		" controlar a ese potro salvaje que es Internet  720p")

	var fp map[string]fmPart
	fp, e = extract(ps)
	require.NoError(t, e)
	require.Equal(t, 6, len(fp))

	d := &VideoData{
		Site:  "YouTube",
		Title: ps.title,
		Type:  "video",
	}
	d.Formats, e = fillFormatInfo(fp, make(map[string][]string),
		h.DefaultClient)
	require.NotNil(t, d.Formats)
	_, ok := e.(*url.Error)
	require.True(t, ok)
}
