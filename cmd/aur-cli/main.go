package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Jguer/aur"
)

func getSearchBy(value string) aur.By {
	switch value {
	case "name":
		return aur.Name
	case "maintainer":
		return aur.Maintainer
	case "depends":
		return aur.Depends
	case "makedepends":
		return aur.MakeDepends
	case "optdepends":
		return aur.OptDepends
	case "checkdepends":
		return aur.CheckDepends
	default:
		return aur.NameDesc
	}
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "<opts>", "<command>", "<pkg(s)>")
	fmt.Println("Available commands:", "info, search")
	fmt.Println("Available opts:", "-by <Search for packages using a specified field>")
	fmt.Println("name/name-desc/maintainer/depends/makedepends/optdepends/checkdepends")
}

func main() {
	if len(os.Args) < 3 {
		usage()

		os.Exit(1)
	}

	var by string
	flag.StringVar(&by, "by", "name-desc", "Search for packages using a specified field")
	flag.Parse()

	aurClient, err := aur.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	var results []aur.Pkg
	if os.Args[1] == "search" {
		results, err = aurClient.Search(context.Background(), strings.Join(os.Args[2:], " "), getSearchBy(by))
	} else if os.Args[1] == "info" {
		results, err = aurClient.Info(context.Background(), os.Args[2:])
	} else {
		usage()

		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	fmt.Println(results)
}
