package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/partfile/data"
	"github.com/jasontconnell/partfile/process"
)

func main() {
	start := time.Now()
	size := flag.Int("s", 10_000_000, "size of the partitions")
	file := flag.String("f", "", "the file")
	mode := flag.String("m", "partition", "mode of operation (join partitions or partition a file)")
	flag.Parse()

	m := data.Partition
	if *mode == "join" {
		m = data.Join
	}

	if file == nil || *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var err error
	if m == data.Partition {
		err = process.Partition(*file, *size)
	} else {
		err = process.Join(*file)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("finished", time.Since(start))
}