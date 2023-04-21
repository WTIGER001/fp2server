package common

import "os"

func init() {

	s := &LocalFSStorage{
		rootDir: "E:\\fp2",
	}
	Storage = s
	os.MkdirAll(s.rootDir, 0600)

	References = newReferenceManager()
	Text = NewTextManager()
}
