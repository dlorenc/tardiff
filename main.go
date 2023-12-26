package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"maps"
	"os"

	"github.com/google/go-cmp/cmp"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s tarball1 tarball2\n", os.Args[0])
		os.Exit(1)
	}

	tarball1, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening tarball1: %v\n", err)
		os.Exit(1)
	}
	defer tarball1.Close()

	tarball2, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening tarball2: %v\n", err)
		os.Exit(1)
	}
	defer tarball2.Close()

	gz1, err1 := gzip.NewReader(tarball1)

	gz2, err2 := gzip.NewReader(tarball2)

	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Error reading tarball1: %v %v\n", err1, err2)
		os.Exit(1)
	}

	tr1 := tar.NewReader(gz1)
	tr2 := tar.NewReader(gz2)

	for {
		hdr1, err := tr1.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading tarball1: %v\n", err)
			os.Exit(1)
		}

		hdr2, err := tr2.Next()
		if err == io.EOF {
			fmt.Printf("File only in tarball1: %s\n", hdr1.Name)
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading tarball2: %v\n", err)
			os.Exit(1)
		}

		if hdr1.Name != hdr2.Name {
			fmt.Printf("Files differ: %s (tarball1) vs %s (tarball2)\n", hdr1.Name, hdr2.Name)
			continue
		}

		if !maps.Equal(hdr1.PAXRecords, hdr2.PAXRecords) {
			fmt.Printf("PAX headers differ for file %s\n", hdr1.Name)
			fmt.Println(cmp.Diff(hdr1.PAXRecords, hdr2.PAXRecords))
		}

		if hdr1.Size != hdr2.Size {
			fmt.Printf("File sizes differ for file %s: %d (tarball1) vs %d (tarball2)\n", hdr1.Name, hdr1.Size, hdr2.Size)
		} else {
			buf1 := make([]byte, hdr1.Size)
			buf2 := make([]byte, hdr2.Size)

			_, err := io.ReadFull(tr1, buf1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s from tarball1: %v\n", hdr1.Name, err)
				os.Exit(1)
			}

			_, err = io.ReadFull(tr2, buf2)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s from tarball2: %v\n", hdr1.Name, err)
				os.Exit(1)
			}

			if !bytes.Equal(buf1, buf2) {
				fmt.Printf("File contents differ for file %s\n", hdr1.Name)
				fmt.Println(cmp.Diff(buf1, buf2))
			}
		}
	}
}
