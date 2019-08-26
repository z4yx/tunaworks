package main

import prober "github.com/z4yx/tunaworks/prober"

func main() {
	cfg := prober.ProberConfig{
		Debug:    true,
		IPv4:     true,
		Server:   "192.168.1.50",
		Token:    "041c9bb6-5fbe-45a8-98da-1d2951e4f862",
		Interval: 60,
	}
	s := prober.MakeProber(&cfg)
	s.Run()
}
