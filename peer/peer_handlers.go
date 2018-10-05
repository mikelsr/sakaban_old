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
		break
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

/* helper functions */

// recvBlockContent reads all the content of a comm.BlockContent, even if it's
// splitted
func recvBlockContent(s net.Stream) (*comm.BlockContent, error) {
	buf := make([]byte, bufferSize)
	n, err := s.Read(buf)
	if err != nil {
		return nil, err
	}
	buf = buf[:n]

	bc := comm.BlockContent{}
	bc.MessageSize = bc.Size(buf)

	for uint64(len(buf)) < bc.MessageSize {
		recv := make([]byte, bufferSize)
		n, err = s.Read(recv)
		if err != nil {
			return nil, err
		}
		recv = recv[:n]
		buf = append(buf, recv...)
	}

	if err = bc.Load(buf); err != nil {
		return nil, err
	}
	return &bc, nil
}
