package main

import (
    "database/sql"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "github.com/deckarep/golang-set"
    _ "github.com/lib/pq"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func page(i uint64) mapset.Set {
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
    chanSet := mapset.NewSet()
    for i := 0; i < len(channels.Nodes); i++ {
        split := strings.Split(channels.Get(i).Attr[0].Val, "/")[2]
        chanSet.Add(split)
        fmt.Println(split)
    }

    return chanSet
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

func insert(db *sql.DB, channel string) {
    sqlInsert := "INSERT INTO youtube.entities.channels (serial) VALUES ($1) ON CONFLICT (serial) DO NOTHING"

    _, err := db.Exec(sqlInsert, channel)
    if err != nil {
        panic(err)
    }
}

func main() {
    for {
        pages := max()
        fmt.Println("Found", pages, "pages")
        dur, err := time.ParseDuration("3s")
        if err != nil {
            panic(err)
        }

        connStr := "user=root dbname=youtube host=localhost port=5432 sslmode=disable"
        db, err := sql.Open("postgres", connStr)
        if err != nil {
            panic(err)
        }

        var i uint64
        for i = 1; i <= pages; i++ {
            time.Sleep(dur)
            channels := page(i)
            fmt.Println("Found", channels.Cardinality(), "channels", "on page", i)

            sliceSet := channels.ToSlice()
            for j := 0; j < channels.Cardinality(); j++ {
                var str string
                str = sliceSet[j].(string)
                insert(db, str)
            }
        }
    }

}
