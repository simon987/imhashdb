package imhashdb

import (
	"go.uber.org/zap"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var ReImgurImg = regexp.MustCompile("(?:https?://)?(?:www\\.|[im]\\.)?imgur\\.com/(\\w{7}|\\w{5})[sbtmlh]?")
var ReImgurAlbum = regexp.MustCompile("(?:https?://)?(?:www\\.|[im]\\.)?imgur\\.com/a/(\\w{7}|\\w{5})")

type ImgurImgResp struct {
	Data ImgurImg `json:"data"`
}

type ImgurImg struct {
	Link string `json:"link"`
}

type ImgurAlbumResp struct {
	Data struct {
		Images []ImgurImg `json:"images"`
	} `json:"data"`
}

func IsImageLink(link string) bool {

	if strings.HasPrefix(link, "https://i.reddituploads.com") {
		return true
	}

	u, err := url.Parse(link)
	if err != nil {
		return false
	}

	path := strings.ToLower(u.Path)
	for _, suffix := range ImageSuffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}

	return false
}

func handleImgurLink(link string, meta *[]Meta) []string {

	if strings.HasPrefix(link, "https://imgur.fun/") {
		link = strings.Replace(link, "imgur.fun", "imgur.com", 1)
	}

	if ReImgurImg.MatchString(link) {
		id := ReImgurImg.FindStringSubmatch(link)[1]

		var img ImgurImgResp
		var rawJson []byte
		err := FetchJson(
			"https://api.imgur.com/3/image/"+id,
			&img, &rawJson,
			[]string{"Authorization", "Client-Id " + Conf.ImgurClientId},
		)
		if err != nil {
			return nil
		}

		Logger.Debug("Got ImgurImgResp", zap.String("id", id))
		*meta = append(*meta, Meta{RetrievedAt: time.Now().Unix(), Id: "imgur.i." + id, Meta: rawJson})

		return []string{img.Data.Link}

	} else if ReImgurAlbum.MatchString(link) {
		id := ReImgurAlbum.FindStringSubmatch(link)[1]

		var album ImgurAlbumResp
		var rawJson []byte
		err := FetchJson(
			"https://api.imgur.com/3/album/"+id,
			&album, &rawJson,
			[]string{"Authorization", "Client-Id " + Conf.ImgurClientId},
		)
		if err != nil {
			return nil
		}
		Logger.Debug(
			"Got ImgurAlbumResp",
			zap.String("id", id),
			zap.Int("count", len(album.Data.Images)),
		)
		*meta = append(*meta, Meta{RetrievedAt: time.Now().Unix(), Id: "imgur.a." + id, Meta: rawJson})

		var links = make([]string, len(album.Data.Images))
		for i, img := range album.Data.Images {
			links[i] = img.Link
		}
		return links
	}

	return nil
}
