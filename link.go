package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var files stringSlice
	flag.Var(&files, "file", "Link mapping files to use.")
	flag.Parse()

	if len(files) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	mappings, err := loadMappings(files)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := server{mappings: mappings}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func(s *server) {
		for range c {
			mappings, err := loadMappings(files)
			if err != nil {
				fmt.Println(err)
				return
			}
			s.setMappings(mappings)
		}
	}(&s)

	http.Handle("/", &s)
	err = http.ListenAndServe(":9876", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
