package util

import (
	"io"
	"log"
)

func Sync(source1 io.ReadWriteCloser, source2 io.ReadWriteCloser) {
	go func() {
		defer closeAll(source2, source1)

		_, err := io.Copy(source2, source1)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	_, err := io.Copy(source1, source2)
	if err != nil {
		log.Fatalln(err)
	}
}

func closeAll(sources ...io.Closer) {
	for _, source := range sources {
		source.Close()
	}
	log.Println("Closed all connections")
}
