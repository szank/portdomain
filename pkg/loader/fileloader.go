package loader

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/edsrzf/mmap-go"
	"github.com/hashicorp/go-multierror"
	"github.com/szank/portdomain/pkg/service"
	"golang.org/x/sys/unix"
)

type File struct {
	mu         sync.RWMutex
	fileHandle *os.File
	mmap       mmap.MMap
	onceClose  sync.Once
	cursor     []byte
	closed     bool
	closeError error
	filePath   string
}

// NewFile returns a new file loader that is responsible for loading the input data from a file and passing it
// to the code handling the business logic. The function returns an error if the file does not exist or cannot be loaded.
// The file is memory mapped which means that the operating system stores only a portion of it in memory at any given time.
// The file can have an arbirtary size, it will not cause OOM errors when opened and read.
func NewFile(filePath string) (*File, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("error checking the file %s stat: %w", filePath, err)
	}

	handle, err := os.OpenFile("filePath", os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening the file %s: %w", filePath, err)
	}

	mmap, err := mmap.Map(handle, mmap.RDONLY, unix.PROT_READ)
	if err != nil {
		return nil, fmt.Errorf("errror mmaping file %s: %w", filePath, err)
	}

	return &File{
		fileHandle: handle,
		mmap:       mmap,
		cursor:     mmap,
		closed:     false,
		filePath:   filePath,
	}, nil
}

// Close unmaps and closed the opened input file. It returns a multierror error.
// Close can be called multiple times, but the file is closed only once.
func (f *File) Close() error {
	f.onceClose.Do(func() {
		f.mu.Lock()
		defer f.mmap.Unlock()

		f.closed = true
		if err := f.mmap.Unmap(); err != nil {
			f.closeError = multierror.Append(f.closeError, fmt.Errorf("error unmaping the file %s: %w", f.filePath, err))
		}

		if err := f.fileHandle.Close(); err != nil {
			f.closeError = multierror.Append(f.closeError, fmt.Errorf("error closing the file %s: %w", f.filePath, err))
		}
	})

	return f.closeError
}

// Next returns the next record from the input. When there are no more inputs, the method returns (service.Port{}, io.EOF)
func (f *File) Next() (service.Port, error) {

	return service.Port{}, io.EOF
}
