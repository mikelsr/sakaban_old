package peer

import (
	"errors"
	"log"
	"path/filepath"

	"bitbucket.org/mikelsr/sakaban/fs"
	"bitbucket.org/mikelsr/sakaban/peer/comm"
	net "github.com/libp2p/go-libp2p-net"
)

func (p *Peer) handleRequest(s net.Stream, msgType comm.MessageType, msg []byte) error {
	switch msgType {
	case comm.MTBlockContent:
		break
	case comm.MTBlockRequest:
		br := comm.BlockRequest{}
		if err := br.Load(msg); err != nil {
			return errors.New("Error unmarshalling BlockRequest")
		}
		return p.handleRequestMTBlockRequest(s, br)
	case comm.MTIndexContent:
		break
	case comm.MTIndexRequest:
		ir := comm.IndexRequest{}
		if err := ir.Load(msg); err != nil {
			return errors.New("Error unmarshalling IndexRequest")
		}
		return p.handleRequestMTIndexRequest(s, ir)
	}
	return nil
}

func (p *Peer) handleRequestMTBlockRequest(s net.Stream, br comm.BlockRequest) error {
	absPath := filepath.Join(p.RootDir, br.FilePath)
	prettyID := p.Host.ID().Pretty()
	prettyID = prettyID[len(prettyID)-4:]
	if s, found := p.RootIndex.Files[absPath]; !found || s.ID != br.FileID.String() {
		return errors.New("File not found")
	}
	f, err := fs.MakeFile(absPath)
	if err != nil {
		return errors.New("Error loading file")
	}
	if len(f.Blocks) <= int(br.BlockN) {
		return errors.New("Invalid block number")
	}
	log.Printf("[P_%s]\tFile loaded: %s", prettyID, absPath)
	blockSize := fs.BlockSize / 1024
	bc := comm.BlockContent{
		BlockN:    br.BlockN,
		BlockSize: uint16(blockSize),
		Content:   f.Blocks[br.BlockN].Content,
		FileID:    f.ID,
	}
	raw := bc.Dump()
	log.Printf("[P_%s]\tSending block %d of file: %s", prettyID, bc.BlockN, absPath)
	if n, err := s.Write(raw); n != len(raw) || err != nil {
		return errors.New("Error writing to steam")
	}
	return nil
}

func (p *Peer) handleRequestMTIndexRequest(s net.Stream, ir comm.IndexRequest) error {
	// TODO: ReloadIndex as a background routine
	ic := comm.IndexContent{Index: p.RootIndex}
	raw := ic.Dump()
	if n, err := s.Write(raw); n != len(raw) || err != nil {
		return errors.New("Error writing to steam")
	}
	return nil
}
