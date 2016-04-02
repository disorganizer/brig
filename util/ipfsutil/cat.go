package ipfsutil

import (
	"io"

	log "github.com/Sirupsen/logrus"
	coreunix "github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/jbenet/go-multihash"
)

// Reader accumulates all the interfaces that a brig stream needs to support.
type Reader interface {
	io.Reader
	io.Seeker
	io.Closer
	io.WriterTo
}

// Cat returns an io.Reader that reads from ipfs.
func Cat(node *Node, hash multihash.Multihash) (Reader, error) {
	nd, err := node.proc()
	if err != nil {
		return nil, err
	}

	reader, err := coreunix.Cat(node.Context, nd, hash.B58String())
	if err != nil {
		log.Warningf("ipfs cat: %v", err)
		return nil, err
	}

	return reader, nil
}
