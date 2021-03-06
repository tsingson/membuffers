// Copyright 2018 the membuffers authors
// This file is part of the membuffers library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package e2e

import (
	"bytes"
	"github.com/orbs-network/membuffers/go/membufc/e2e/protos"
	"github.com/orbs-network/membuffers/go/membufc/e2e/protos/crypto"
	"testing"
)

func TestReadWriteWithInline(t *testing.T) {
	// write
	builder := &types.FileRecordBuilder{
		Data:        []byte{0x01, 0x02, 0x03},
		Hash:        crypto.Sha256{0x04, 0x05, 0x06},
		AnotherHash: []crypto.Md5{{0x07, 0x08}, {0x09}},
		BlockHeight: 0x22,
	}
	buf := make([]byte, builder.CalcRequiredSize())
	err := builder.Write(buf)
	if err != nil {
		t.Fatalf("error while writing in builder")
	}

	// read
	record := types.FileRecordReader(buf)
	if !bytes.Equal(record.Hash(), crypto.Sha256{0x04, 0x05, 0x06}) {
		t.Fatalf("FileRecord.Hash is not as expected")
	}
	iter := record.AnotherHashIterator()
	if !bytes.Equal(iter.NextAnotherHash(), crypto.Md5{0x07, 0x08}) {
		t.Fatalf("FileRecord.AnotherHash[0] is not as expected")
	}
	if !bytes.Equal(iter.NextAnotherHash(), crypto.Md5{0x09}) {
		t.Fatalf("FileRecord.AnotherHash[0] is not as expected")
	}
	if record.BlockHeight() != 0x22 {
		t.Fatalf("FileRecord.BlockHeight is not as expected")
	}
	if record.Hash().String() != "040506" {
		t.Fatalf("String: FileRecord.Hash is not as expected")
	}
	if record.BlockHeight().String() != "22" {
		t.Fatalf("String: FileRecord.BlockHeight is not as expected")
	}
	if !record.Hash().Equal(crypto.Sha256{0x04, 0x05, 0x06}) {
		t.Fatalf("Equal == : FileRecord.Hash is not as expected")
	}
	if record.Hash().Equal(crypto.Sha256{0x05, 0x06, 0x07}) {
		t.Fatalf("Equal != : FileRecord.Hash is not as expected")
	}
	if !record.BlockHeight().Equal(0x22) {
		t.Fatalf("Equal == : FileRecord.BlockHeight is not as expected")
	}
	if record.BlockHeight().Equal(0x23) {
		t.Fatalf("Equal != : FileRecord.BlockHeight is not as expected")
	}
	m1 := make(map[string]bool)
	m1[record.Hash().KeyForMap()] = true
	if !m1[record.Hash().KeyForMap()] {
		t.Fatalf("FileRecord.Hash cannot be used as a key in a map")
	}
	m2 := make(map[uint64]bool)
	m2[record.BlockHeight().KeyForMap()] = true
	if !m2[record.BlockHeight().KeyForMap()] {
		t.Fatalf("FileRecord.BlockHeight cannot be used as a key in a map")
	}
}
