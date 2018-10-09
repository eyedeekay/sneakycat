package torcat

import (
	//"fmt"
	"io"
	"log"
	"os"
	//"strconv"
	"strings"

	"github.com/yawning/bulb"
	//"github.com/yawning/bulb/utils"
)

var Verbose bool

func ioCopy(errchan chan error, w io.Writer, r io.Reader) {
	_, err := io.Copy(w, r)
	errchan <- err
}

func runIO(conn io.ReadWriter) error {
	errchan := make(chan error)
	go ioCopy(errchan, os.Stdout, conn)
	go ioCopy(errchan, conn, os.Stdin)

	select {
	case err := <-errchan:
		return err
	}
}

func Listen(listen *uint, torctl *bulb.Conn, torcfg *bulb.NewOnionConfig) {
	if l, err := torctl.NewListener(torcfg, uint16(*listen)); err != nil {
		log.Fatalf("Failed to listen port: %s", err)
	} else {
		defer l.Close()
		addrVec := strings.SplitN(l.Addr().String(), ":", 2)
		os.Stderr.WriteString(addrVec[0])
		os.Stderr.WriteString("\n")
		if Verbose {
			os.Stderr.WriteString("[Waiting]")
			os.Stderr.WriteString("\n")
		}
		if conn, err := l.Accept(); err != nil {
			log.Fatalf("Failed to accept connection: %s", err)
		} else {
			defer conn.Close()
			if Verbose {
				os.Stderr.WriteString("[Connected]")
				os.Stderr.WriteString("\n")
			}
			if err := runIO(conn); err != nil {
				log.Fatalf("Failed conversation: %s", err)
			}
		}
	}
}

func Connect(torctl *bulb.Conn, dest string) {
	if dialer, err := torctl.Dialer(nil); err != nil {
		log.Fatalf("Failed to get Dialer: %s", err)
	} else if conn, err := dialer.Dial("tcp", dest); err != nil {
		log.Fatalf("Connection to %s failed", err)
	} else {
		defer conn.Close()
		if Verbose {
			os.Stderr.WriteString("[Connected]")
			os.Stderr.WriteString("\n")
		}
		if err := runIO(conn); err != nil {
			log.Fatalf("Failed conversation: %s", err)
		}
	}
}
