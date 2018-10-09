Sneakycat: A netcat-like tool for Tor
=====================================

Originally based on torcat, this was split to make it easier to embed in
existing Go applications to allow themselves to automatically self-forward to
Tor Hidden Services, and allow client applications written in Go to more easily
use Tor Control Port. Right now it doesn't do very much, but it's growing.
