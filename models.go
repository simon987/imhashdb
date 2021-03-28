package imhashdb

import "github.com/simon987/fastimagehash-go"

type HashType string

const (
	DHash8  HashType = "dhash8"
	DHash16 HashType = "dhash16"
	DHash32 HashType = "dhash32"

	MHash8  HashType = "mhash8"
	MHash16 HashType = "mhash16"
	MHash32 HashType = "mhash32"

	PHash8  HashType = "phash8"
	PHash16 HashType = "phash16"
	PHash32 HashType = "phash32"

	WHash8Haar  HashType = "whash8haar"
	WHash16Haar HashType = "whash16haar"
	WHash32Haar HashType = "whash32haar"
)

var HashTypes = []HashType{
	DHash8, DHash16, DHash32,
	MHash8, MHash16, MHash32,
	PHash8, PHash16, PHash32,
	WHash8Haar, WHash16Haar, WHash32Haar,
}

func (h HashType) HashLength() int {
	switch h {
	case DHash8:
		fallthrough
	case MHash8:
		fallthrough
	case PHash8:
		fallthrough
	case WHash8Haar:
		return 8

	case DHash16:
		fallthrough
	case MHash16:
		fallthrough
	case PHash16:
		fallthrough
	case WHash16Haar:
		return 32

	case DHash32:
		fallthrough
	case MHash32:
		fallthrough
	case PHash32:
		fallthrough
	case WHash32Haar:
		return 128
	default:
		panic("Invalid invalid hash")
	}
}

type HashReq struct {
	Data []byte `json:"data"`
}
type Hashes struct {
	DHash8  *fastimagehash.Hash `json:"dhash8"`
	DHash16 *fastimagehash.Hash `json:"dhash16"`
	DHash32 *fastimagehash.Hash `json:"dhash32"`

	MHash8  *fastimagehash.Hash `json:"mhash8"`
	MHash16 *fastimagehash.Hash `json:"mhash16"`
	MHash32 *fastimagehash.Hash `json:"mhash32"`

	PHash8  *fastimagehash.Hash `json:"phash8"`
	PHash16 *fastimagehash.Hash `json:"phash16"`
	PHash32 *fastimagehash.Hash `json:"phash32"`

	WHash8  *fastimagehash.Hash `json:"whash8haar"`
	WHash16 *fastimagehash.Hash `json:"whash16haar"`
	WHash32 *fastimagehash.Hash `json:"whash32haar"`
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
	Err string `json:"err,omitempty"`
}

type Meta struct {
	RetrievedAt int64       `json:"retrieved_at"`
	Id          string      `json:"id"`
	Meta        interface{} `json:"meta"`
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
