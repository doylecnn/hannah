package provider

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/lyric"
	"github.com/missdeer/hannah/util"
	"github.com/missdeer/hannah/util/cryptography"
)

const (
	kuwoAPISearch         = `https://www.kuwo.cn/api/www/search/searchMusicBykeyWord?key=%s&pn=%d&rn=%d`
	kuwoAPIToken          = `https://www.kuwo.cn/search/list?key=`
	kuwoAPIConvertURL     = `https://antiserver.kuwo.cn/anti.s?type=convert_url&format=mp3|aac|wma&response=url&rid=%s`
	kuwoAPIGetLossless    = "https://mobi.kuwo.cn/mobi.s?f=kuwo&q="
	kuwoAPIHot            = `https://www.kuwo.cn/www/categoryNew/getPlayListInfoUnderCategory?type=taglist&digest=10000&id=37&start=%d&count=%d`
	kuwoAPIPlaylistDetail = `https://nplserver.kuwo.cn/pl.svc?op=getlistinfo&pn=0&rn=200&encode=utf-8&keyset=pl2012&pcmp4=1&pid=%s&vipver=MUSIC_9.0.2.0_W1&newver=1`
	kuwoAPIGetLyric       = `https://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId=%s`
	kuwoAPIPreArtistSongs = `https://www.kuwo.cn/api/www/artist/artist?artistid=%s`
	kuwoAPIArtistSongs    = `https://www.kuwo.cn/api/www/artist/artistMusic?artistid=%s&pn=%d&rn=%d`
	kuwoAPIAlbumSongs     = `https://www.kuwo.cn/api/www/album/albumInfo?albumId=%s&pn=%d&rn=%d`
)

var (
	ErrEmptyTrackList = errors.New("empty track list")
	ErrEmptyKuwoToken = errors.New("empty kuwo token")
	ErrEmptyKuwoLRC   = errors.New("empty kuwo lyric")
)

type kuwo struct {
}

type kuwoSearchResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total string `json:"total"`
		List  []struct {
			MusicRID string `json:"musicrid"`
			Artist   string `json:"artist"`
			Pic      string `json:"pic"`
			RID      int    `json:"rid"`
			Album    string `json:"album"`
			AlbumID  string `json:"albumid"`
			AlbumPic string `json:"albumpic"`
			Pic120   string `json:"pic120"`
			Name     string `json:"name"`
		} `json:"list"`
	} `json:"data"`
}

func (p *kuwo) getToken() (string, error) {
	req, err := http.NewRequest("GET", kuwoAPIToken, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kuwo.cn/")
	req.Header.Set("Origin", "http://www.kuwo.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", ErrStatusNotOK
	}

	parsedURL, _ := url.Parse(kuwoAPIToken)
	c := httpClient.Jar.Cookies(parsedURL)
	const kuwoToken = "kw_token"
	for _, cookie := range c {
		if cookie.Name == kuwoToken {
			return cookie.Value, nil
		}
	}
	return "", ErrEmptyKuwoToken
}

func (p *kuwo) SearchSongs(keyword string, page int, limit int) (SearchResult, error) {
	token, err := p.getToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf(kuwoAPISearch, url.QueryEscape(keyword), page, limit)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kuwo.cn/")
	req.Header.Set("Origin", "http://www.kuwo.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("csrf", token)

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var sr kuwoSearchResult
	err = json.Unmarshal(content, &sr)
	if err != nil {
		return nil, err
	}

	var res SearchResult
	for _, r := range sr.Data.List {
		id := r.MusicRID
		if strings.HasPrefix(id, "MUSIC_") {
			id = id[len("MUSIC_"):]
		}
		res = append(res, Song{
			ID:       id,
			Title:    r.Name,
			Image:    r.Pic120,
			Artist:   r.Artist,
			Provider: "kuwo",
		})
	}

	return res, nil
}

func (p *kuwo) ResolveSongURL(song Song) (Song, error) {
	token, err := p.getToken()
	u := kuwoAPIGetLossless + base64.StdEncoding.EncodeToString(cryptography.DESEncrypt([]byte("corp=kuwo&p2p=1&type=convert_url2&sig=0&format=flac|mp3&rid="+song.ID)))
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "http://www.kuwo.cn/search/list?key=The+Call")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("csrf", token)
	req.Header.Set("cookie", "kw_token="+token)

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return song, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return song, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return song, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		b := scanner.Bytes()
		if bytes.HasPrefix(b, []byte("url=")) {
			song.URL = string(b[len(`url=`):])
			break
		}
	}
	song.Provider = "kuwo"
	return song, nil
}

type kuwoLyric struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		LRCList []struct {
			LineLyric string `json:"lineLyric"`
			Time      string `json:"time"`
		} `json:"lrclist"`
		SongInfo struct {
			SongName string `json:"songName"`
			Album    string `json:"album"`
			ID       string `json:"id"`
			Artist   string `json:"artist"`
			Pic      string `json:"pic"`
		} `json:"songinfo"`
	} `json:"data"`
}

func (p *kuwo) ResolveSongLyric(song Song, format string) (Song, error) {
	u := fmt.Sprintf(kuwoAPIGetLyric, song.ID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kuwo.cn/")
	req.Header.Set("Origin", "http://www.kuwo.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return song, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return song, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return song, err
	}

	var lrc kuwoLyric
	err = json.Unmarshal(content, &lrc)
	if err != nil {
		return song, err
	}
	if len(lrc.Data.LRCList) == 0 {
		return song, ErrEmptyKuwoLRC
	}
	var lines []string
	for _, l := range lrc.Data.LRCList {
		tt := strings.Split(l.Time, ".")
		var timestamp string
		if len(tt) == 2 {
			seconds, err := strconv.Atoi(tt[0])
			if err == nil {
				minutes := seconds / 60
				seconds = seconds % 60
				millisecond, _ := strconv.Atoi(tt[1])
				timestamp = fmt.Sprintf("%02d:%02d.%02d", minutes, seconds, millisecond%100)
			}
		}
		lines = append(lines, fmt.Sprintf(`[%s]%s`, timestamp, l.LineLyric))
	}
	song.Lyric = lyric.LyricConvert("lrc", format, strings.Join(lines, "\n"))
	return song, nil
}

type kuwoHotPlaylists struct {
	Data []struct {
		Img   string `json:"img"`
		Total string `json:"total"`
		Data  []struct {
			Img   string `json:"img"`
			UName string `json:"uname"`
			Name  string `json:"name"`
			UID   string `json:"uid"`
			Total string `json:"total"`
			ID    string `json:"id"`
		} `json:"data"`
		Start string `json:"start"`
		Count string `json:"count"`
		Name  string `json:"name"`
		ID    string `json:"id"`
		Type  string `json:"type"`
	} `json:"data"`
	Msg    string `json:"msg"`
	RegID  string `json:"regid"`
	Status int    `json:"status"`
}

func (p *kuwo) HotPlaylist(page int, limit int) (res Playlists, err error) {
	u := fmt.Sprintf(kuwoAPIHot, (page-1)*limit+1, limit)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kuwo.cn/")
	req.Header.Set("Origin", "http://www.kuwo.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var hot kuwoHotPlaylists
	if err = json.Unmarshal(content, &hot); err != nil {
		return nil, err
	}

	if len(hot.Data) == 0 {
		return nil, errors.New("empty playlist")
	}

	for _, d := range hot.Data[0].Data {
		res = append(res, Playlist{
			ID:       d.ID,
			Title:    d.Name,
			Provider: "kuwo",
			URL:      fmt.Sprintf(`http://kuwo.cn/playlist_detail/%s`, d.ID),
		})
	}

	return res, nil
}

type kuwoPlaylistDetail struct {
	ID        int    `json:"id"`
	Info      string `json:"info"`
	Pic       string `json:"pic"`
	Title     string `json:"title"`
	Total     int    `json:"total"`
	MusicList []struct {
		ID     string `json:"id"`
		Format string `json:"format"`
		Artist string `json:"artist"`
		Name   string `json:"name"`
	} `json:"musiclist"`
}

func (p *kuwo) PlaylistDetail(pl Playlist) (Songs, error) {
	u := fmt.Sprintf(kuwoAPIPlaylistDetail, pl.ID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kuwo.cn/")
	req.Header.Set("Origin", "http://www.kuwo.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var pld kuwoPlaylistDetail
	if err = json.Unmarshal(content, &pld); err != nil {
		return nil, err
	}

	var songs Songs
	for _, p := range pld.MusicList {
		songs = append(songs, Song{
			ID:       p.ID,
			Title:    p.Name,
			Artist:   p.Artist,
			Provider: "kuwo",
		})
	}
	if len(songs) == 0 {
		return nil, ErrEmptyTrackList
	}
	return songs, nil
}

type kuwoArtistSongs struct {
	Code int `json:"code"`
	Data struct {
		Total int `json:"total"`
		List  []struct {
			MusicRID string `json:"musicrid"`
			Artist   string `json:"artist"`
			ArtistID int    `json:"artistid"`
			Pic      string `json:"pic"`
			RID      int    `json:"rid"`
			Name     string `json:"name"`
		} `json:"list"`
	} `json:"data"`
}

func (p *kuwo) ArtistSongs(id string) (res Songs, err error) {
	token, err := p.getToken()
	for page := config.Page; ; page++ {
		u := fmt.Sprintf(kuwoAPIArtistSongs, id, page, config.Limit)
		log.Println("kuwo artist songs page", page)
		req, e := http.NewRequest("GET", u, nil)
		if e != nil {
			log.Println("new request failed", e)
			err = e
			break
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:79.0) Gecko/20100101 Firefox/79.0")
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Referer", "http://www.kuwo.cn/")
		req.Header.Set("Origin", "http://www.kuwo.cn/")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("csrf", token)
		req.Header.Set("cookie", "kw_token="+token)

		httpClient := util.GetHttpClient()
		resp, e := httpClient.Do(req)
		if e != nil {
			log.Println("http client do request failed", e)
			err = e
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println("http response code != 200", resp.StatusCode)
			if resp.StatusCode == 504 {
				second := rand.Intn(8) + 3
				time.Sleep(time.Duration(second * int(time.Second)))
				page--
				// try again
				continue
			}
			err = ErrStatusNotOK
			break
		}

		content, e := util.ReadHttpResponseBody(resp)
		if e != nil {
			log.Println("read http response body failed", e)
			err = e
			break
		}

		var artistSongs kuwoArtistSongs
		err = json.Unmarshal(content, &artistSongs)
		if err != nil {
			log.Println("unmarshal content failed", e, string(content))
			break
		}
		if len(artistSongs.Data.List) == 0 {
			log.Println("empty music list")
			break
		}
		for _, song := range artistSongs.Data.List {
			res = append(res, Song{
				ID:       strconv.Itoa(song.RID),
				Title:    song.Name,
				Artist:   song.Artist,
				Image:    song.Pic,
				Provider: "kuwo",
			})
		}
	}

	if len(res) == 0 {
		return nil, ErrEmptyTrackList
	} else {
		log.Printf("got %d songs totally\n", len(res))
		err = nil
	}
	return
}

type kuwoAlbumSongs struct {
	Code int `json:"code"`
	Data struct {
		Artist    string `json:"artist"`
		ArtistID  int    `json:"artistid"`
		Album     string `json:"album"`
		AlbumID   int    `json:"albumid"`
		Pic       string `json:"pic"`
		MusicList []struct {
			MusicRID    string `json:"musicrid"`
			RID         int    `json:"rid"`
			Artist      string `json:"artist"`
			ArtistID    int    `json:"artistid"`
			HasLossless bool   `json:"hasLossless"`
			Album       string `json:"album"`
			AlbumID     int    `json:"albumid"`
			Name        string `json:"name"`
			Pic         string `json:"pic120"`
		} `json:"musicList"`
	} `json:"data"`
}

func (p *kuwo) AlbumSongs(id string) (res Songs, err error) {
	token, err := p.getToken()
	for page := config.Page; ; page++ {
		u := fmt.Sprintf(kuwoAPIAlbumSongs, id, page, config.Limit)
		log.Println("kuwo album songs page", page)
		req, e := http.NewRequest("GET", u, nil)
		if e != nil {
			err = e
			break
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:79.0) Gecko/20100101 Firefox/79.0")
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Referer", "http://www.kuwo.cn/")
		req.Header.Set("Origin", "http://www.kuwo.cn/")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("csrf", token)
		req.Header.Set("cookie", "kw_token="+token)

		httpClient := util.GetHttpClient()
		resp, e := httpClient.Do(req)
		if e != nil {
			err = e
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			err = ErrStatusNotOK
			break
		}

		content, e := util.ReadHttpResponseBody(resp)
		if e != nil {
			err = e
			break
		}

		var albumSongs kuwoAlbumSongs
		err = json.Unmarshal(content, &albumSongs)
		if err != nil {
			break
		}
		if len(albumSongs.Data.MusicList) == 0 {
			break
		}
		for _, song := range albumSongs.Data.MusicList {
			res = append(res, Song{
				ID:       strconv.Itoa(song.RID),
				Title:    song.Name,
				Artist:   song.Artist,
				Image:    song.Pic,
				Provider: "kuwo",
			})
		}
	}
	if len(res) == 0 {
		return nil, ErrEmptyTrackList
	} else {
		err = nil
	}
	return
}

func (p *kuwo) Login() error {
	return ErrNotImplemented
}

func (p *kuwo) RefreshToken() error {
	return ErrNotImplemented
}

func (p *kuwo) Name() string {
	return "kuwo"
}
