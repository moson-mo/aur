package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Jguer/aur"
)

const (
	boldCode  = "\x1b[1m"
	resetCode = "\x1b[0m"
)

// UseColor determines if package will emit colors.
var UseColor = true // nolint

const (
	searchMode = "search"
	infoMode   = "info"
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

	flag.Usage()

	fmt.Println("Example:", "aur-cli -verbose -by name search python3.7")
}

func versionRequestEditor(ctx context.Context, req *http.Request) error {
	req.Header.Add("User-Agent", "aur-cli/v1")

	return nil
}

func main() {
	var (
		by          string
		aurURL      string
		verbose     bool
		jsonDisplay bool
	)

	flag.StringVar(&by, "by", "name-desc", "Search for packages using a specified field"+
		"\n (name/name-desc/maintainer/depends/makedepends/optdepends/checkdepends)")
	flag.StringVar(&aurURL, "url", "https://aur.archlinux.org/", "AUR URL")
	flag.BoolVar(&verbose, "verbose", false, "display verbose information")
	flag.BoolVar(&jsonDisplay, "json", false, "display result as JSON")
	flag.Parse()

	if flag.NArg() < 2 {
		usage()

		os.Exit(1)
	}

	mode := flag.Arg(0)

	aurClient, err := aur.NewClient(aur.WithBaseURL(aurURL),
		aur.WithRequestEditorFn(versionRequestEditor))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	results, err := getResults(aurClient, by, mode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	if jsonDisplay {
		output, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		fmt.Println(string(output))
	} else {
		for i := range results {
			switch mode {
			case infoMode:
				printInfo(&results[i], os.Stdout, aurURL, verbose)
			case searchMode:
				printSearch(&results[i], os.Stdout)
			default:
				usage()
				os.Exit(1)
			}
		}
	}
}

func getResults(aurClient *aur.Client, by, mode string) (results []aur.Pkg, err error) {
	switch mode {
	case searchMode:
		results, err = aurClient.Search(context.Background(), strings.Join(flag.Args()[1:], " "), getSearchBy(by))
	case infoMode:
		results, err = aurClient.Info(context.Background(), flag.Args()[1:])
	default:
		usage()
		os.Exit(1)
	}

	if err != nil {
		err = fmt.Errorf("rpc request failed: %w", err)
	}

	return results, err
}

func stylize(startCode, in string) string {
	if UseColor {
		return startCode + in + resetCode
	}

	return in
}

func Bold(in string) string {
	return stylize(boldCode, in)
}

func printSearch(a *aur.Pkg, w io.Writer) {
	fmt.Fprintf(w, "- %s %s (%d %.2f)\n\t%s\n",
		Bold(a.Name), a.Version, a.NumVotes, a.Popularity, a.Description)
}

// PrintInfo prints package info like pacman -Si.
func printInfo(a *aur.Pkg, w io.Writer, aurURL string, verbose bool) {
	printInfoValue(w, "Name", a.Name)
	printInfoValue(w, "Version", a.Version)
	printInfoValue(w, "Description", a.Description)

	if verbose {
		printInfoValue(w, "Keywords", a.Keywords...)
		printInfoValue(w, "URL", a.URL)
		printInfoValue(w, "AUR URL", strings.TrimRight(aurURL, "/")+"/packages/"+a.Name)

		printInfoValue(w, "Groups", a.Groups...)
		printInfoValue(w, "Licenses", a.License...)
		printInfoValue(w, "Provides", a.Provides...)
		printInfoValue(w, "Depends On", a.Depends...)
		printInfoValue(w, "Make Deps", a.MakeDepends...)
		printInfoValue(w, "Check Deps", a.CheckDepends...)
		printInfoValue(w, "Optional Deps", a.OptDepends...)
		printInfoValue(w, "Conflicts With", a.Conflicts...)

		printInfoValue(w, "Maintainer", a.Maintainer)
		printInfoValue(w, "Votes", fmt.Sprintf("%d", a.NumVotes))
		printInfoValue(w, "Popularity", fmt.Sprintf("%f", a.Popularity))
		printInfoValue(w, "First Submitted", formatTimeQuery(a.FirstSubmitted))
		printInfoValue(w, "Last Modified", formatTimeQuery(a.LastModified))

		if a.OutOfDate != 0 {
			printInfoValue(w, "Out-of-date", formatTimeQuery(a.OutOfDate))
		} else {
			printInfoValue(w, "Out-of-date", "No")
		}

		printInfoValue(w, "ID", fmt.Sprintf("%d", a.ID))
		printInfoValue(w, "Package Base ID", fmt.Sprintf("%d", a.PackageBaseID))
		printInfoValue(w, "Package Base", a.PackageBase)
		printInfoValue(w, "Snapshot URL", aurURL+a.URLPath)
	}

	fmt.Fprintln(w)
}
