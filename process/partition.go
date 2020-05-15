package process

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

func Partition(dir, filename string, size int64) error {
	fullPath := filepath.Join(dir, filename)
	stat, err := os.Stat(fullPath)

	if err != nil {
		return errors.Wrapf(err, "couldn't stat file for partition %s", filename)
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
		return errors.Wrapf(err, "opening file to partition %s", filename)
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
			return errors.Wrapf(err, "opening partiion %d, %s", i, outfile)
		}

		_, err = out.Write(buf[:n])
		if err != nil {
			return errors.Wrapf(err, "writing partition file %d, %s", i, outfile)
		}

		i++
	}

	return nil
}
