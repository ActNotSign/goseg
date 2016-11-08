package main

import (
    "log"
    "github.com/huichen/sego"
    "github.com/codegangsta/cli"
    "sync"
    "net/http"
    "os"
    "path/filepath"
    "path"
    "fmt"
    "io"
    "encoding/json"
    "strings"
)

var (
    VERSION = "1.0.1"
    NAME    = "goseg"
    httpPort = ":6800"
)

var(
    segmenter Segmenter
    dictPath string
)

type Segmenter struct {
    sego    *sego.Segmenter
    mutex   sync.Mutex
}

type JsonResult struct {
    Ok     bool     `json:"ok"`
    Words  []string `json:"words"`
    Msg    string   `json:"error"`
}

func main () {
    app := cli.NewApp()
    app.Name = NAME
    app.Usage = "http service"
    app.Version = VERSION
    app.Action = func(c *cli.Context) error {
        log.Println("run service")
        // get dict path
        dict := "dict"
        if c.Args().Get(0) != "" {
            dict = c.Args().Get(0)
        }
        filepath.Walk(dict, walkHandler)

        // load dict
        segmenter.Load()

        // run server
        if c.Args().Get(1) != "" {
            httpPort = c.Args().Get(1)
        }
        log.Println("listen ", httpPort)
        httpService(httpPort)
        return nil
    }
    app.Run(os.Args)
}

// get dict filename
func walkHandler(filename string, info os.FileInfo, err error) error {
    if path.Ext(filename) == ".txt" {
        if dictPath != "" {
            dictPath = fmt.Sprintf("%s,%s", dictPath, filename)
        } else {
            dictPath = filename
        }
    }
    return nil
}

func httpService(port string) {
    http.HandleFunc("/segment", segmentHandler)
    err := http.ListenAndServe(port, nil)
    if err != nil {
        panic(err)
    }
}

func segmentHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    r.ParseMultipartForm(32 << 20)
    content := r.FormValue("content")
    tagging := r.FormValue("postagging")
    log.Println(r.RemoteAddr, content, tagging)
    if content == "" {
        out, err := json.Marshal(&JsonResult{Ok: false, Words: []string{}, Msg: "content is empty"})
        checkError(err)
        io.WriteString(w, string(out))
    } else {
        var words = make([]string, 0)
        if tagging == "true" {
            result := segmenter.SegmentsToString([]byte(content), false)
            for _, word := range strings.Split(result, " ") {
                words = append(words, word)
            }
        } else {
            words = segmenter.SegmentsToSlice([]byte(content), false)
        }
        out, err := json.Marshal(&JsonResult{Ok: true, Words: words})
        checkError(err)
        io.WriteString(w, string(out))
    }
}

func checkError(err error) error {
    if err != nil {
        log.Println(err)
    }
    return err
}

// segment
func (s *Segmenter) Load() {
    s.sego= new(sego.Segmenter)
    s.sego.LoadDictionary(dictPath)
}

func (s *Segmenter) GetCurrent() (*sego.Segmenter){
    return s.sego
}

func (s *Segmenter) SegmentsToString(text []byte, model bool) (output string) {
    return sego.SegmentsToString(s.sego.Segment(text), model)
}

func (s *Segmenter) SegmentsToSlice(text []byte, model bool) (output []string) {
    return sego.SegmentsToSlice(s.sego.Segment(text), model)
}

