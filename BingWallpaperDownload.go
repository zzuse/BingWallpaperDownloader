package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"regexp"
	"strings"
	"sync"
	"sync/atomic"
	_ "time"
)

const (
	bufferSize = 1280 * 1024 //one pic buffer
)

var (
	numPoller = flag.Int("p", 1, "page loader num")
	savePath  = flag.String("s", "./downloads/", "save path")
)

type image struct {
	url      string
	filename string
}

//TODO:
//filename add date
//json to interface{}?

type beautiContext struct {
	pollerDone chan struct{}
	images     map[string]int
	imagesLock *sync.Mutex
	imageChan  chan *image
	pageIndex  int32
	rootURL    string
	done       bool
	okCounter  int32
}

type Message struct {
	Images   []ImagesType
	Tooltips TooltipsType
}
type ImagesType struct {
	Startdate     string
	Fullstartdate string
	Enddate       string
	Url           string
	Urlbase       string
	Copyright     string
	Copyrightlink string
	Quiz          string
	Wp            bool
	Hsh           string
	Drk           int
	Top           int
	Bot           int
	Hs            []TablesType
}
type TooltipsType struct {
	Loading  string
	Previous string
	Next     string
	Walle    string
	Walls    string
}
type TablesType struct {
	Id string
}

func main() {
	flag.Parse()
	ctx := &beautiContext{
		pollerDone: make(chan struct{}),
		images:     make(map[string]int),
		imagesLock: &sync.Mutex{},
		imageChan:  make(chan *image, 100),
		pageIndex:  1,
		//rootURL: "http://cn.bing.com/HPImageArchive.aspx?format=js&idx={}&n=1&nc=1421741858945&pid=hp",
		rootURL: "http://cn.bing.com/HPImageArchive.aspx?format=js&idx={}&n=7",
	}
	os.MkdirAll(*savePath, 0777)
	ctx.start()
}

func (ctx *beautiContext) start() {
	fmt.Printf("Poller%d\n", *numPoller)
	for i := 0; i < *numPoller; i++ {
		fmt.Printf("download%d\n", *numPoller)
		go ctx.downloadPage()
	}
	waits := sync.WaitGroup{}

	<-ctx.pollerDone
	ctx.done = true
	//close(ctx.pollerDone)
	waits.Wait()
	fmt.Printf("fetch done get img ok %d\n", ctx.okCounter)
}

func (ctx *beautiContext) downloadPage() {
	isDone := false
	for !isDone {
		select {
		case <-ctx.pollerDone:
			isDone = true
		default:
			url := fmt.Sprintf("%s", ctx.rootURL)
			fmt.Printf("download page %s\n", url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("failed to load url %s with error %v", url, err)
			} else {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("failed to load url %s with error %v", url, err)
				} else {
					ctx.parsePage(body)
				}
			}
		}
	}
}

func (ctx *beautiContext) parsePage(body []byte) {
	body2 := string(body)
	fmt.Printf("get json file: %s\n", body2)
	dec := json.NewDecoder(strings.NewReader(body2))
	for {
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("decode json file: %s \n", m.Images)

		idx := m.Images
		if idx == nil {
			ctx.pollerDone <- struct{}{}
		} else {
			fmt.Printf("idx len: %d\n", len(idx))
			for _, n := range idx {
				url := fmt.Sprintf("http://cn.bing.com%s", n.Url)
				str := strings.Split(url, "/")
				fmt.Printf("split url by slash: %s \n", url)
				length := len(str)
				fmt.Printf("slashed array: %s \n", str)
				imgeUrl := url
				fmt.Printf("the URL: %s \n", imgeUrl)
				//get filename by "?"
				tmpfilename := strings.Split(str[length-1], "?")
				fmt.Printf("the last field is filename: %s \n", tmpfilename)
				filename := tmpfilename[0]
				image := &image{url: imgeUrl, filename: filename}
				fmt.Printf("start download %s\n", image.url)
				resp, err := http.Get(image.url)
				if err != nil {
					fmt.Printf("failed to load url %s with error %v\n", image.url, err)
				} else {
					defer resp.Body.Close()
					saveFile := *savePath + image.filename //path.Base(imgUrl)

					img, err := os.Create(saveFile)
					if err != nil {
						fmt.Print(err)

					} else {
						defer img.Close()

						imgWriter := bufio.NewWriterSize(img, bufferSize)

						_, err = io.Copy(imgWriter, resp.Body)
						if err != nil {
							fmt.Print(err)
						}
						imgWriter.Flush()
						fmt.Printf("finish download %s\n", image.url)
					}
				}
				atomic.AddInt32(&ctx.okCounter, 1)
				if ctx.okCounter > 6 {
					ctx.pollerDone <- struct{}{}
					fmt.Printf("counter greater than 7 abort download %s\n", image.url)
					break
				}
			}
		}
	}
}
