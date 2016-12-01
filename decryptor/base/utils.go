// Copyright 2016, Cossack Labs Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package base

import (
	"bytes"
	"encoding/binary"
	"github.com/cossacklabs/acra/utils"
	"github.com/cossacklabs/themis/gothemis/cell"
	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/cossacklabs/themis/gothemis/message"
)

const (
	LENGTH_SIZE = 8
)

func DecryptAcrastruct(data []byte, private_key *keys.PrivateKey, zone []byte) ([]byte, error) {
	inner_data := data[len(TAG_BEGIN):]
	pubkey := &keys.PublicKey{Value: inner_data[:PUBLIC_KEY_LENGTH]}
	smessage := message.New(private_key, pubkey)
	symmetric_key, err := smessage.Unwrap(inner_data[PUBLIC_KEY_LENGTH:KEY_BLOCK_LENGTH])
	if err != nil {
		return []byte{}, err
	}
	//
	var length uint64
	// convert from little endian
	err = binary.Read(bytes.NewReader(inner_data[KEY_BLOCK_LENGTH:KEY_BLOCK_LENGTH+LENGTH_SIZE]), binary.LittleEndian, &length)
	if err != nil {
		return []byte{}, err
	}
	scell := cell.New(symmetric_key, cell.CELL_MODE_SEAL)
	decrypted, err := scell.Unprotect(inner_data[KEY_BLOCK_LENGTH+LENGTH_SIZE:], nil, zone)
	// fill zero symmetric_key
	utils.FillSlice(byte(0), symmetric_key)
	if err != nil {
		return []byte{}, err
	}
	return decrypted, nil
}
