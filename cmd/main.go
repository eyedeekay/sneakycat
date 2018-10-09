package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/yawning/bulb"
	"github.com/yawning/bulb/utils"

	"github.com/eyedeekay/torcat"
)

const (
	flagDescControl = "Control socket (unix/tcp)"
	flagDescListen  = "Listen port"
	flagDescVerbose = "Verbose output"
)

var (
	control = flag.String("control", "9051", flagDescControl)
	listen  = flag.Uint("listen", 0, flagDescListen)
	verbose = flag.Bool("verbose", false, flagDescVerbose)

	torctl *bulb.Conn
	torcfg = &bulb.NewOnionConfig{
		DiscardPK: true,
	}
	dest string
)

func init() {
	flag.StringVar(control, "c", "9051", flagDescControl)
	flag.UintVar(listen, "l", 0, flagDescListen)
	flag.BoolVar(verbose, "v", false, flagDescVerbose)

	flag.Parse()
	torcat.Verbose = true
}

func main() {
	if cproto, caddr, err := utils.ParseControlPortString(*control); err != nil {
		log.Fatalf("Failed to parse control socket: %s", err)
	} else if torctl, err = bulb.Dial(cproto, caddr); err != nil {
		log.Fatalf("Failed to connect to control socket: %s", err)
	} else {
		defer torctl.Close()
	}

	if err := torctl.Authenticate(os.Getenv("TORCAT_COOKIE")); err != nil {
		log.Fatalf("Authentication failed: %s", err)
	}

	if *listen > 65535 {
		log.Fatalf("Listen port %d is greater than 65535", *listen)
	} else if *listen != 0 {
		torcat.Listen(listen, torctl, torcfg)
	} else {
		if len(os.Args) != 3 {
			log.Fatalf("Invalid arguments. Must be `%s host port'", os.Args[0])
		} else if addr := os.Args[1]; len(addr) == 0 {
			log.Fatalf("Empty destination address")
		} else if port, err := strconv.Atoi(os.Args[2]); err != nil {
			log.Fatalf("Invalid port number: %s", err)
		} else {
			dest = fmt.Sprintf("%s:%d", addr, port)
		}
		torcat.Connect(torctl, dest)
	}
}
