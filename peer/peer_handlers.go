package peer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bitbucket.org/mikelsr/sakaban/fs"
	"bitbucket.org/mikelsr/sakaban/peer/comm"
	net "github.com/libp2p/go-libp2p-net"
)

func (p *Peer) handleRequest(s net.Stream, msgType comm.MessageType, msg []byte) error {
	defer s.Close()
	switch msgType {
	case comm.MTBlockContent:
		bc := new(comm.BlockContent)
		if err := bc.Load(msg); err != nil {
			return errors.New("Error unmarshalling BlockContent")
		}
		return p.handleRequestMTBlockContent(s, bc)
	case comm.MTBlockRequest:
		br := comm.BlockRequest{}
		if err := br.Load(msg); err != nil {
			return errors.New("Error unmarshalling BlockRequest")
		}
		return p.handleRequestMTBlockRequest(s, br)
	case comm.MTIndexContent:
		ic := new(comm.IndexContent)
		if err := ic.Load(msg); err != nil {
			return errors.New("Error unmarshalling BlockContent")
		}
		return p.handleRequestMTIndexContent(s, ic)
	case comm.MTIndexRequest:
		ir := comm.IndexRequest{}
		if err := ir.Load(msg); err != nil {
			return errors.New("Error unmarshalling IndexRequest")
		}
		return p.handleRequestMTIndexRequest(s, ir)
	}
	return nil
}

func (p *Peer) handleRequestMTBlockContent(s net.Stream, bc *comm.BlockContent) error {
	if len(p.fileMap) == 0 {
		return errors.New("Didn't expect any blocks")
	}

	eid := bc.FileID.String()
	requestedFile, found := p.fileMap[eid]
	if !found {
		return fmt.Errorf("Didn't expect blocks from file %s", eid)
	}

	contact := requestedFile.contact
	file := requestedFile.file
	summary := requestedFile.summary

	// block comes from expected peer
	if s.Conn().RemotePeer().String() != contact.ID().String() {
		return errors.New("Block from unexpected peer")
	}

	// block belongs to expected file
	if summary.ID != eid {
		return fmt.Errorf("File IDs do not match: got %s expected %s",
			file.ID.String(), eid)
	}

	if bc.BlockN > uint8(len(summary.Blocks)) {
		return fmt.Errorf("Block index out of range: max is %d got %d",
			len(file.Blocks), bc.BlockN)
	}

	if summary.Blocks[bc.BlockN] == 0 {
		return fmt.Errorf("Block %d was unchanged", bc.BlockN)
	}

	file.Blocks[bc.BlockN] = &fs.Block{Content: bc.Content}

	// if file is not complete, return
	for i := range file.Blocks {
		if i != 0 {
			if file.Blocks[i] == nil {
				return nil
			}
		}
	}
	file.Write()
	p.fileMap[file.ID.String()] = nil
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

func (p *Peer) handleRequestMTIndexContent(s net.Stream, ir *comm.IndexContent) error {
	if !p.waiting {
		return errors.New("Unexpected index received")
	}

	var contact *Contact
	for _, c := range p.Contacts {
		if c.MultiAddr().Equal(s.Conn().LocalMultiaddr()) ||
			c.MultiAddr().Equal(s.Conn().RemoteMultiaddr()) {
			contact = &c
			break
		}
	}
	if contact == nil {
		return errors.New("Unknown contact")
	}

	i := p.RootIndex
	ni := &ir.Index
	comparison := i.Compare(ni)
	for _, path := range comparison.Deletions {
		// TODO: delete path
		os.Remove(path)
	}

	for _, sum := range comparison.Additions {
		requestedFile, err := MakeRequestedFile(sum, contact)
		if err != nil {
			return err
		}
		p.fileMap[sum.ID] = requestedFile
	}

	// TODO: fileMap is not empty, request and update files of stack
	p.waiting = false
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
