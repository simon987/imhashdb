// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package imhashdb

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	fastimagehash_go "github.com/simon987/fastimagehash-go"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeGithubComSimon987Imhashdb(in *jlexer.Lexer, out *QueryResp) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "err":
			out.Err = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987Imhashdb(out *jwriter.Writer, in QueryResp) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Err != "" {
		const prefix string = ",\"err\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Err))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v QueryResp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v QueryResp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *QueryResp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *QueryResp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb(l, v)
}
func easyjsonD2b7633eDecodeGithubComSimon987Imhashdb1(in *jlexer.Lexer, out *QueryReq) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			out.HashType = HashType(in.String())
		case "hash":
			if in.IsNull() {
				in.Skip()
				out.Hash = nil
			} else {
				out.Hash = in.Bytes()
			}
		case "distance":
			out.Distance = uint(in.Uint())
		case "limit":
			out.Limit = uint(in.Uint())
		case "offset":
			out.Offset = uint(in.Uint())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987Imhashdb1(out *jwriter.Writer, in QueryReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix[1:])
		out.String(string(in.HashType))
	}
	{
		const prefix string = ",\"hash\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Hash)
	}
	{
		const prefix string = ",\"distance\":"
		out.RawString(prefix)
		out.Uint(uint(in.Distance))
	}
	{
		const prefix string = ",\"limit\":"
		out.RawString(prefix)
		out.Uint(uint(in.Limit))
	}
	{
		const prefix string = ",\"offset\":"
		out.RawString(prefix)
		out.Uint(uint(in.Offset))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v QueryReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v QueryReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *QueryReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *QueryReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb1(l, v)
}
func easyjsonD2b7633eDecodeGithubComSimon987Imhashdb2(in *jlexer.Lexer, out *Meta) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "retrieved_at":
			out.RetrievedAt = int64(in.Int64())
		case "id":
			out.Id = string(in.String())
		case "meta":
			if m, ok := out.Meta.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Meta.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Meta = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987Imhashdb2(out *jwriter.Writer, in Meta) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"retrieved_at\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.RetrievedAt))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"meta\":"
		out.RawString(prefix)
		if m, ok := in.Meta.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Meta.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Meta))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Meta) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Meta) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Meta) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Meta) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb2(l, v)
}
func easyjsonD2b7633eDecodeGithubComSimon987Imhashdb3(in *jlexer.Lexer, out *ImageList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "images":
			if in.IsNull() {
				in.Skip()
				out.Images = nil
			} else {
				in.Delim('[')
				if out.Images == nil {
					if !in.IsDelim(']') {
						out.Images = make([]*Image, 0, 8)
					} else {
						out.Images = []*Image{}
					}
				} else {
					out.Images = (out.Images)[:0]
				}
				for !in.IsDelim(']') {
					var v4 *Image
					if in.IsNull() {
						in.Skip()
						v4 = nil
					} else {
						if v4 == nil {
							v4 = new(Image)
						}
						(*v4).UnmarshalEasyJSON(in)
					}
					out.Images = append(out.Images, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987Imhashdb3(out *jwriter.Writer, in ImageList) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"images\":"
		out.RawString(prefix[1:])
		if in.Images == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Images {
				if v5 > 0 {
					out.RawByte(',')
				}
				if v6 == nil {
					out.RawString("null")
				} else {
					(*v6).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ImageList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ImageList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ImageList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ImageList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb3(l, v)
}
func easyjsonD2b7633eDecodeGithubComSimon987Imhashdb4(in *jlexer.Lexer, out *ImageHasMeta) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "url":
			out.Url = string(in.String())
		case "meta":
			(out.Meta).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987Imhashdb4(out *jwriter.Writer, in ImageHasMeta) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(string(in.Url))
	}
	{
		const prefix string = ",\"meta\":"
		out.RawString(prefix)
		(in.Meta).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ImageHasMeta) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ImageHasMeta) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ImageHasMeta) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ImageHasMeta) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb4(l, v)
}
func easyjsonD2b7633eDecodeGithubComSimon987Imhashdb5(in *jlexer.Lexer, out *Image) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "size":
			out.Size = int(in.Int())
		case "sha1":
			if in.IsNull() {
				in.Skip()
				out.Sha1 = nil
			} else {
				out.Sha1 = in.Bytes()
			}
		case "md5":
			if in.IsNull() {
				in.Skip()
				out.Md5 = nil
			} else {
				out.Md5 = in.Bytes()
			}
		case "sha256":
			if in.IsNull() {
				in.Skip()
				out.Sha256 = nil
			} else {
				out.Sha256 = in.Bytes()
			}
		case "crc32":
			out.Crc32 = uint32(in.Uint32())
		case "meta":
			if in.IsNull() {
				in.Skip()
				out.Meta = nil
			} else {
				in.Delim('[')
				if out.Meta == nil {
					if !in.IsDelim(']') {
						out.Meta = make([]ImageHasMeta, 0, 1)
					} else {
						out.Meta = []ImageHasMeta{}
					}
				} else {
					out.Meta = (out.Meta)[:0]
				}
				for !in.IsDelim(']') {
					var v10 ImageHasMeta
					(v10).UnmarshalEasyJSON(in)
					out.Meta = append(out.Meta, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987Imhashdb5(out *jwriter.Writer, in Image) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"size\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Size))
	}
	{
		const prefix string = ",\"sha1\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Sha1)
	}
	{
		const prefix string = ",\"md5\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Md5)
	}
	{
		const prefix string = ",\"sha256\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Sha256)
	}
	{
		const prefix string = ",\"crc32\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.Crc32))
	}
	{
		const prefix string = ",\"meta\":"
		out.RawString(prefix)
		if in.Meta == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v17, v18 := range in.Meta {
				if v17 > 0 {
					out.RawByte(',')
				}
				(v18).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Image) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Image) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Image) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Image) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb5(l, v)
}
func easyjsonD2b7633eDecodeGithubComSimon987Imhashdb6(in *jlexer.Lexer, out *Hashes) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "dhash8":
			if in.IsNull() {
				in.Skip()
				out.DHash8 = nil
			} else {
				if out.DHash8 == nil {
					out.DHash8 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.DHash8)
			}
		case "dhash16":
			if in.IsNull() {
				in.Skip()
				out.DHash16 = nil
			} else {
				if out.DHash16 == nil {
					out.DHash16 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.DHash16)
			}
		case "dhash32":
			if in.IsNull() {
				in.Skip()
				out.DHash32 = nil
			} else {
				if out.DHash32 == nil {
					out.DHash32 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.DHash32)
			}
		case "mhash8":
			if in.IsNull() {
				in.Skip()
				out.MHash8 = nil
			} else {
				if out.MHash8 == nil {
					out.MHash8 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.MHash8)
			}
		case "mhash16":
			if in.IsNull() {
				in.Skip()
				out.MHash16 = nil
			} else {
				if out.MHash16 == nil {
					out.MHash16 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.MHash16)
			}
		case "mhash32":
			if in.IsNull() {
				in.Skip()
				out.MHash32 = nil
			} else {
				if out.MHash32 == nil {
					out.MHash32 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.MHash32)
			}
		case "phash8":
			if in.IsNull() {
				in.Skip()
				out.PHash8 = nil
			} else {
				if out.PHash8 == nil {
					out.PHash8 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.PHash8)
			}
		case "phash16":
			if in.IsNull() {
				in.Skip()
				out.PHash16 = nil
			} else {
				if out.PHash16 == nil {
					out.PHash16 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.PHash16)
			}
		case "phash32":
			if in.IsNull() {
				in.Skip()
				out.PHash32 = nil
			} else {
				if out.PHash32 == nil {
					out.PHash32 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.PHash32)
			}
		case "whash8haar":
			if in.IsNull() {
				in.Skip()
				out.WHash8 = nil
			} else {
				if out.WHash8 == nil {
					out.WHash8 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.WHash8)
			}
		case "whash16haar":
			if in.IsNull() {
				in.Skip()
				out.WHash16 = nil
			} else {
				if out.WHash16 == nil {
					out.WHash16 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.WHash16)
			}
		case "whash32haar":
			if in.IsNull() {
				in.Skip()
				out.WHash32 = nil
			} else {
				if out.WHash32 == nil {
					out.WHash32 = new(fastimagehash_go.Hash)
				}
				easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in, out.WHash32)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987Imhashdb6(out *jwriter.Writer, in Hashes) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"dhash8\":"
		out.RawString(prefix[1:])
		if in.DHash8 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.DHash8)
		}
	}
	{
		const prefix string = ",\"dhash16\":"
		out.RawString(prefix)
		if in.DHash16 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.DHash16)
		}
	}
	{
		const prefix string = ",\"dhash32\":"
		out.RawString(prefix)
		if in.DHash32 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.DHash32)
		}
	}
	{
		const prefix string = ",\"mhash8\":"
		out.RawString(prefix)
		if in.MHash8 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.MHash8)
		}
	}
	{
		const prefix string = ",\"mhash16\":"
		out.RawString(prefix)
		if in.MHash16 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.MHash16)
		}
	}
	{
		const prefix string = ",\"mhash32\":"
		out.RawString(prefix)
		if in.MHash32 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.MHash32)
		}
	}
	{
		const prefix string = ",\"phash8\":"
		out.RawString(prefix)
		if in.PHash8 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.PHash8)
		}
	}
	{
		const prefix string = ",\"phash16\":"
		out.RawString(prefix)
		if in.PHash16 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.PHash16)
		}
	}
	{
		const prefix string = ",\"phash32\":"
		out.RawString(prefix)
		if in.PHash32 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.PHash32)
		}
	}
	{
		const prefix string = ",\"whash8haar\":"
		out.RawString(prefix)
		if in.WHash8 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.WHash8)
		}
	}
	{
		const prefix string = ",\"whash16haar\":"
		out.RawString(prefix)
		if in.WHash16 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.WHash16)
		}
	}
	{
		const prefix string = ",\"whash32haar\":"
		out.RawString(prefix)
		if in.WHash32 == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out, *in.WHash32)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Hashes) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Hashes) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Hashes) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Hashes) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb6(l, v)
}
func easyjsonD2b7633eDecodeGithubComSimon987FastimagehashGo(in *jlexer.Lexer, out *fastimagehash_go.Hash) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "size":
			out.Size = int(in.Int())
		case "bytes":
			if in.IsNull() {
				in.Skip()
				out.Bytes = nil
			} else {
				out.Bytes = in.Bytes()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987FastimagehashGo(out *jwriter.Writer, in fastimagehash_go.Hash) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"size\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Size))
	}
	{
		const prefix string = ",\"bytes\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Bytes)
	}
	out.RawByte('}')
}
func easyjsonD2b7633eDecodeGithubComSimon987Imhashdb7(in *jlexer.Lexer, out *HashReq) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				out.Data = in.Bytes()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComSimon987Imhashdb7(out *jwriter.Writer, in HashReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		out.Base64Bytes(in.Data)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v HashReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v HashReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComSimon987Imhashdb7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *HashReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *HashReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComSimon987Imhashdb7(l, v)
}
