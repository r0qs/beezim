// Copyright 2021 Ethersphere. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This code is based on the beekeeper beeclient api

package tarball

import (
	"bytes"
	"io"

	"github.com/ethersphere/bee/pkg/swarm"
)

// File represents Bee file
type File struct {
	address    swarm.Address
	hash       []byte
	name       string
	dataReader io.Reader
	size       int64
}

// NewBufferFile returns new file with specified buffer
func NewBufferFile(name string, buffer *bytes.Buffer) *File {
	return &File{
		name:       name,
		dataReader: buffer,
		size:       int64(buffer.Len()),
	}
}

// Address returns file's address
func (f *File) Address() swarm.Address {
	return f.address
}

// Name returns file's name
func (f *File) Name() string {
	return f.name
}

// Hash returns file's hash
func (f *File) Hash() []byte {
	return f.hash
}

// DataReader returns file's data reader
func (f *File) DataReader() io.Reader {
	return f.dataReader
}

// Size returns file size
func (f *File) Size() int64 {
	return f.size
}

func (f *File) SetAddress(a swarm.Address) {
	f.address = a
}

func (f *File) SetHash(h []byte) {
	f.hash = h
}
