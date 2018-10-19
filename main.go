package main

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "net/http"
    "strconv"
    "strings"
)

func page(i uint64) []string {
    return []string{}
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

}
