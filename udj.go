package main

import (
    "flag"
    "os"
    "log"
    "net/http"
    "fmt"
    "strings"
    "io/ioutil"
    "encoding/json"
    "github.com/dustin/go-id3"
)

var addr = flag.String("addr", ":8010", "http service address")

func main() {
    flag.Parse()
    http.Handle("/list", http.HandlerFunc(list))
    http.Handle("/music", http.HandlerFunc(music))
    http.Handle("/stream", http.HandlerFunc(stream))
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

type FileInfo struct {
    Name string
    IsDir bool
}

func list(w http.ResponseWriter, req *http.Request) {
    path := req.FormValue("path")
    files, err := ioutil.ReadDir(path)
    if err != nil {
        log.Fatal("problem reading")
    } else {
        var output []FileInfo = make([]FileInfo, len(files))
        for i := range files {
            output[i] = FileInfo{Name: files[i].Name(),IsDir: files[i].IsDir()}
            fmt.Printf(files[i].Name())
        }
        jsonresp, _ := json.Marshal(output)
        fmt.Fprintf(w, "%s", jsonresp)
    }
}

type MusicInfo struct {
    Title string
    Artist string
    Album string
    Length string
    Path string
}

func getMusic(path string) []MusicInfo {
    var music []MusicInfo = make([]MusicInfo,0)
    files, _ := ioutil.ReadDir(path)
    for i := range files {
        newPath := fmt.Sprintf("%s/%s",path,files[i].Name())
        if files[i].IsDir() {
            var dirMusic []MusicInfo = getMusic(newPath)
            music = append(music,dirMusic...)
        } else {
            if strings.HasSuffix(files[i].Name(),".mp3") {
                fd, err := os.Open(newPath)
                if err == nil {
                    defer fd.Close()
                    info := id3.Read(fd)
                    if info != nil {
                        newMusicInfo := MusicInfo{
                            Title: info.Name, 
                            Artist: info.Artist,
                            Album: info.Album,
                            Length: info.Length,
                            Path: newPath}

                        music = append(music,newMusicInfo)
                    }
                }
            }
        }
    }
    return music
}

func music(w http.ResponseWriter, req *http.Request) {
    path := req.FormValue("path")
    output := getMusic(path)
    jsonresp, _ := json.Marshal(output)
    fmt.Fprintf(w, "%s", jsonresp)
}

func stream(w http.ResponseWriter, req *http.Request) {
    path := req.FormValue("path")
    http.ServeFile(w, req, path)
}