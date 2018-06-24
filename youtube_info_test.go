package miyam

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInfo(t *testing.T) {
	bs, e := ioutil.ReadFile("youtube_player_line.json")
	require.NoError(t, e)
	testHTML := string(bs)
	var ytplayer string
	ytplayer, e = matchOneOf(testHTML,
		`;ytplayer\.config\s*=\s*({.+?});`)
	require.NoError(t, e)
	yd := new(youtubeData)
	e = json.Unmarshal([]byte(ytplayer), yd)
	require.NoError(t, e)
	d := &VideoData{
		Site:  "YouTube",
		Title: yd.Args.Title,
	}
	require.Equal(t, d.Title, "Eliécer Ávila   El régimen no puede"+
		" controlar a ese potro salvaje que es Internet  720p")

	fs := make(map[string]FormatData)
	streams := append(strings.Split(yd.Args.Stream, ","),
		strings.Split(yd.Args.Audio, ",")...)
	for i := 0; e == nil && i != len(streams); i++ {
		var stream url.Values
		stream, e = url.ParseQuery(streams[i])
		var fm FormatData
		if e == nil {
			fm.Quality = stream.Get("quality")
			fm.Ext, e = matchOneOf(stream.Get("type"),
				`video/([[:word:]]+);`, `audio/([[:word:]]+)`)
		}
		if e == nil {
			fm.URL = stream.Get("url")
			itag := stream.Get("itag")
			fs[itag] = fm
		}
	}
	require.Equal(t, 5, len(fs))
}
