package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"io"
	"log"
	"sync"
)

func MD5FromReaderAsync(r io.Reader, ss *sync.WaitGroup, errChan chan error, outChan chan []byte) {
	defer ss.Done()

	buf := make([]byte, 1*1024)
	s := md5.New()
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, err := s.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	errChan <- nil
	outChan <- s.Sum(nil)
}

func Sha1FromReaderAsync(r io.Reader, ss *sync.WaitGroup, errChan chan error, outChan chan []byte) {
	defer ss.Done()

	buf := make([]byte, 1*1024)
	s := sha1.New()
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, err := s.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	errChan <- nil
	outChan <- s.Sum(nil)
}

func Sha256FromReaderAsync(r io.Reader, ss *sync.WaitGroup, errChan chan error, outChan chan []byte) {
	defer ss.Done()

	buf := make([]byte, 1*1024)
	s := sha256.New()
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, err := s.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	errChan <- nil
	outChan <- s.Sum(nil)
}

func Sha512FromReaderAsync(r io.Reader, ss *sync.WaitGroup, errChan chan error, outChan chan []byte) {
	defer ss.Done()

	buf := make([]byte, 1*1024)
	s := sha512.New()
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, err := s.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	errChan <- nil
	outChan <- s.Sum(nil)
}
