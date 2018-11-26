package peer

import "bitbucket.org/mikelsr/sakaban/fs"

// fileStack is used to store the files to be retrieved from another peer
type fileStack struct {
	files []*fs.Summary
}

func newFileStack() *fileStack {
	return &fileStack{files: make([]*fs.Summary, 0)}
}

func (f *fileStack) pop() (*fs.Summary, int) {
	lenght := len(f.files) - 1
	if lenght == -1 {
		return nil, -1
	}
	s := f.files[lenght]
	f.files = f.files[:lenght]
	return s, lenght
}

func (f *fileStack) push(s *fs.Summary) {
	f.files = append(f.files, s)
}
