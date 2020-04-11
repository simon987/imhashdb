package imhashdb

type HashType string

const (
	AHash12    HashType = "ahash:12"
	DHash12    HashType = "dhash:12"
	MHash12    HashType = "mhash:12"
	PHash12    HashType = "phaash:12:4"
	WHash8Haar HashType = "whash:8:haar"
)

type HashReq struct {
	Data []byte `json:"data"`
}

type HashResp struct {
	AHash []byte `json:"ahash:12"`
	DHash []byte `json:"dhash:12"`
	MHash []byte `json:"mhash:12"`
	PHash []byte `json:"phash:12:4"`
	WHash []byte `json:"whash:18:haar"`
}

type QueryReq struct {
	HashType HashType `json:"type"`
	Hash     []byte   `json:"hash"`
	Distance uint     `json:"distance"`
	Limit    uint     `json:"limit"`
	Offset   uint     `json:"offset"`
}

type ImageList struct {
	Images []*Image `json:"images"`
}

type QueryResp struct {
	Err      string `json:"err,omitempty"`
}

type Meta struct {
	RetrievedAt int64  `json:"retrieved_at"`
	Id          string `json:"id"`
	Meta        []byte `json:"meta"`
}

type ImageHasMeta struct {
	Url  string `json:"url"`
	Meta Meta   `json:"meta"`
}

type Image struct {
	id     int64
	Size   int    `json:"size"`
	Sha1   []byte `json:"sha1"`
	Md5    []byte `json:"md5"`
	Sha256 []byte `json:"sha256"`
	Crc32  uint32 `json:"crc32"`

	Meta []ImageHasMeta `json:"meta"`
}
