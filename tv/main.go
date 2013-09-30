package main

import (
    "flag"
    "fmt"
    "github.com/goopi/shows"
    "log"
    "os"
)

var (
    searchTerm = flag.String("s", "", "show name")
    showId = flag.String("e", "", "show id")
)

func usage() {
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
}

func main() {
    flag.Parse()
    if *searchTerm == "" && *showId == "" {
        usage()
    }

    if *searchTerm != "" {
        showList := shows.Search(*searchTerm)

        for _, show := range showList {
            fmt.Println(show)
        }
    }

    if *showId != "" {
        episodeList, err := shows.GetEpisodes(*showId)
        if err != nil { log.Fatal(err) }

        for _, epi := range episodeList {
            fmt.Println(epi)
        }
    }
}
