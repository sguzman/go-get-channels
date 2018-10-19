package main

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "net/http"
    "strconv"
    "strings"
)

func page(i uint64) []string {
    var url string
    if i < 2 {
        url = "https://dbase.tube/chart/channels/subscribers/all"
    } else {
        url = fmt.Sprintf("https://dbase.tube/chart/channels/subscribers/all?page=%d", i)
    }

    res, err := http.Get(url)
    if err != nil {
        panic(err)
    }

    defer func() {
        err := res.Body.Close()
        if err != nil {
            panic(err)
        }
    }()
    if res.StatusCode != 200 {
        panic(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        panic(err)
    }

    channels := doc.Find("a[href^=\"/c/UC\"]")
    chans := make([]string, 0)
    for i := 0; i < len(channels.Nodes); i++ {
        split := strings.Split(channels.Get(i).Attr[0].Val, "/")[2]
        chans = append(chans, split)
        fmt.Println(split)
    }

    return chans
}

func max() uint64 {
    res, err := http.Get("https://dbase.tube/chart/channels/subscribers/all")
    if err != nil {
        panic(err)
    }
    defer res.Body.Close()
    if res.StatusCode != 200 {
        panic(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        panic(err)
    }

    links := doc.Find("a[href^=\"/chart/channels/subscribers/all?page=\"]")
    pageCountStr := links.Last().Get(0).Attr[0].Val
    countIdx := strings.Index(pageCountStr, "=")
    countStr := pageCountStr[countIdx+1:len(pageCountStr)]

    if count, err := strconv.ParseUint(countStr, 10, 64); err == nil {
        return count
    } else {
        panic(err)
    }
}

func main() {
    pages := max()
    fmt.Println("Found", pages, "pages")

    var i uint64
    for i = 1; i <= pages; i++ {
        channels := page(i)
        fmt.Println("Found", len(channels), "channels", "on page", i)
    }

}
