package peer

import (
	"errors"

	"bitbucket.org/mikelsr/sakaban/fs"
	uuid "github.com/satori/go.uuid"
)

// RequestedFile stores a file requested to a peer and the contact of that peer
type RequestedFile struct {
	contact *Contact
	file    *fs.File
	summary *fs.Summary
}

// MakeRequestedFile creates a RequestFile given a contact and a file summary
func MakeRequestedFile(s *fs.Summary, c *Contact) (*RequestedFile, error) {
	if c == nil || s == nil {
		return nil, errors.New("Nil parameter")
	}

	f, err := fs.MakeFileFromSummary(s)
	// if file doesn't exist
	if err != nil {
		id, err := uuid.FromString(s.ID)
		if err != nil {
			return nil, err
		}
		var parentID uuid.UUID
		if s.Parent == "" {
			parentID = uuid.Nil
		} else {
			parentID, err = uuid.FromString(s.Parent)
			if err != nil {
				return nil, err
			}
		}
		f = &fs.File{
			ID:     id,
			Parent: parentID,
			Path:   s.Path,
			Perm:   s.Perm,
			Blocks: make([]*fs.Block, len(s.Blocks)),
		}
		if err != nil {
			return nil, err
		}
	} else {
		// empty changed blocks
		for i := range f.Blocks {
			if s.Blocks[i] != 0 {
				f.Blocks[i] = nil
			}
		}

	}

	return &RequestedFile{
		contact: c,
		file:    f,
		summary: s,
	}, nil
}
