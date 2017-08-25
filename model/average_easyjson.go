// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

import (
	json "encoding/json"
	"fmt"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC312ac33DecodeGithubComLa0rgHighloadcupModel(in *jlexer.Lexer, out *Avg) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "avg":
			out.Value = float64(in.Float64())
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
func easyjsonC312ac33EncodeGithubComLa0rgHighloadcupModel(out *jwriter.Writer, in Avg) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"avg\":")
	out.RawString(fmt.Sprintf("%.5f", in.Value)) // EDITED
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Avg) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC312ac33EncodeGithubComLa0rgHighloadcupModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Avg) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC312ac33EncodeGithubComLa0rgHighloadcupModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Avg) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC312ac33DecodeGithubComLa0rgHighloadcupModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Avg) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC312ac33DecodeGithubComLa0rgHighloadcupModel(l, v)
}