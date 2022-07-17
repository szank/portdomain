package loader

import (
	"bytes"
	"encoding/json"
	"errors"
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
	// There's no point in using RwMutex unless we make it concurrent.
	mu         sync.Mutex
	fileHandle *os.File
	mmap       mmap.MMap

	onceClose  sync.Once
	closed     bool
	closeError error

	filePath string
	decoder  *json.Decoder

	initialBracketFound bool
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

	handle, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening the file %s: %w", filePath, err)
	}

	mmap, err := mmap.Map(handle, unix.PROT_READ, mmap.RDONLY)
	if err != nil {
		return nil, fmt.Errorf("errror mmaping file %s: %w", filePath, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(mmap))

	return &File{
		fileHandle: handle,
		mmap:       mmap,
		decoder:    decoder,
		closed:     false,
		filePath:   filePath,
	}, nil
}

// Close unmaps and closed the opened input file. It returns a multierror error.
// Close can be called multiple times, but the file is closed only once.
func (f *File) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.onceClose.Do(func() {

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
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return service.Port{}, errors.New("the input file has been clossed")
	}

	if !f.initialBracketFound {
		err := f.findInitialBracket()
		if err != nil {
			return service.Port{}, err
		}
	}

	err := f.skipJsonKey()
	if err != nil {
		return service.Port{}, err
	}

	startOffset, err := f.skipJsonObject()
	if err != nil {
		return service.Port{}, err
	}

	finishOffset := f.decoder.InputOffset()

	// Input validation is misssing here. I assume name should always be there for ex.
	// Or coordinates should be valid if they are set. And so on.
	port := &service.Port{}
	err = json.Unmarshal(f.mmap[startOffset:finishOffset], port)
	if err != nil {
		return service.Port{}, nil
	}

	return *port, nil
}

func (f *File) findInitialBracket() error {
	for {
		t, err := f.decoder.Token()
		if err != nil {
			return err
		}

		delim, ok := t.(json.Delim)
		if !ok {
			return fmt.Errorf("expected an opening bracket, got %s", t)
		}
		if delim.String() != "{" {
			return fmt.Errorf("expected an opening bracket, got %s", t)
		} else {
			// not a fan of unnecessary else statements, but this should silencce the linter
			f.initialBracketFound = true
			break
		}
	}

	return nil
}

func (f *File) skipJsonKey() error {
	for {
		t, err := f.decoder.Token()
		if err != nil {
			return err
		}

		switch token := t.(type) {
		case string:
			return nil
		case json.Delim:
			if token.String() == "}" {
				// this is a closing bracket, assume there's nothing left to read
				return io.EOF
			}
			return fmt.Errorf("expected JSON key, got %s", token)
		default:
			return fmt.Errorf("expected JSON key, got %s", token)
		}
	}

}

// we need to open and close equal number of brackets, assuming first bracket is
// an opening one.
func (f *File) skipJsonObject() (int64, error) {
	openingBracketsCount := 0
	closingBracketsCount := 0

	t, err := f.decoder.Token()
	if err != nil {
		return 0, err
	}

	inputOffset := f.decoder.InputOffset() - 1

	delim, ok := t.(json.Delim)
	if !ok {
		return 0, fmt.Errorf("expected an opening bracket, got %s", t)
	}
	if delim.String() != "{" {
		return 0, fmt.Errorf("expected an opening bracket, got %s", t)
	}
	openingBracketsCount++

	for {
		t, err := f.decoder.Token()
		if err != nil {
			return 0, err
		}

		delim, ok := t.(json.Delim)
		if !ok {
			// we don't care about anything besides brackets
			continue
		}

		switch delim.String() {
		case "{":
			openingBracketsCount++
		case "}":
			closingBracketsCount++
		default:
		}

		if openingBracketsCount == closingBracketsCount {
			return inputOffset, nil
		}
	}
}
