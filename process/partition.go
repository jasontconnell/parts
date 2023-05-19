package process

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func Partition(dir, filename string, size int64) error {
	fullPath := filepath.Join(dir, filename)
	stat, err := os.Stat(fullPath)

	if err != nil {
		return fmt.Errorf("couldn't stat file for partition %s. %w", filename, err)
	}

	chunks := stat.Size()/int64(size) + 1
	padlen := len(strconv.FormatInt(chunks, 10))

	// format like filename.0001 filename.0002 etc
	outfmt := filename + ".%0" + strconv.Itoa(padlen) + "d"

	return doPartition(dir, fullPath, outfmt, chunks, size)
}

func doPartition(dir, filename, outfmt string, chunks, size int64) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening file to partition %s. %w", filename, err)
	}
	defer f.Close()

	buf := make([]byte, size)
	var i int64
	for i < chunks {
		n, err := f.Read(buf)

		if err == io.EOF {
			break
		}

		outfile := filepath.Join(dir, fmt.Sprintf(outfmt, i))

		out, err := os.OpenFile(outfile, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return fmt.Errorf("opening partiion %d, %s. %w", i, outfile, err)
		}

		_, err = out.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("writing partition file %d, %s. %w", i, outfile, err)
		}

		i++
	}

	return nil
}
