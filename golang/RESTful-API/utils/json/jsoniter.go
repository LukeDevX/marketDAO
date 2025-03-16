// Copyright 2017 Bo-Yi Wu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package json

import (
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

var (
	json                = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal             = json.Marshal
	MarshalToString     = json.MarshalToString
	Unmarshal           = json.Unmarshal
	UnmarshalFromString = json.UnmarshalFromString
	MarshalIndent       = json.MarshalIndent
	NewDecoder          = json.NewDecoder
	Valid               = json.Valid
)
