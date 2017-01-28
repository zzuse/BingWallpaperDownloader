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
	"strings"
	"sync"
	"sync/atomic"
	_ "time"
)

const (
	bufferSize = 1280 * 1024 //one pic buffer
)

var (
	numPoller  = flag.Int("d", 1, "date num")
	savePath   = flag.String("s", "./downloads/", "save path")
	resolution = flag.String("r", "1920x1200", "resolution")
	Market     = []string{
		"en-US",
		"zh-CN",
		"ja-JP",
		"en-AU",
		"en-UK",
		"de-DE",
		"en-NZ",
	}
	wg sync.WaitGroup
)

type image struct {
	url      string
	filename string
}

type beautiContext struct {
	images     map[string]int
	imagesLock *sync.Mutex
	pageIndex  int32
	rootURL    string
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
	f, err := os.OpenFile("downloader.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a downloader log entry")
	ctx := &beautiContext{
		images:     make(map[string]int),
		imagesLock: &sync.Mutex{},
		pageIndex:  1,
		rootURL:    "http://www.bing.com/HPImageArchive.aspx?format=js&idx={}&n=",
	}
	os.MkdirAll(*savePath, 0777)
	ctx.start()
}

func (ctx *beautiContext) start() {
	for i := 0; i < len(Market); i++ {
		log.Println("download Market", Market[i])
		go ctx.downloadPage(Market[i])
		wg.Add(1)
	}

	wg.Wait()
	fmt.Printf("fetch done get img ok %d\n", ctx.okCounter)
}

func (ctx *beautiContext) downloadPage(market string) {
	select {
	default:
		url := fmt.Sprintf("%s%d&mkt=%s", ctx.rootURL, *numPoller, market)
		log.Println("download json page", url)
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

func (ctx *beautiContext) parsePage(body []byte) {
	defer wg.Done()
	body2 := string(body)
	//fmt.Printf("get json file: %s\n", body2)
	dec := json.NewDecoder(strings.NewReader(body2))
	var m Message
	if err := dec.Decode(&m); err == io.EOF {
		log.Fatal(err)
	} else if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("decode json file: %s \n", m.Images)
	oldfilename := ""

	idx := m.Images
	if idx == nil {
		log.Fatal("no image in json")
	} else {
		log.Println("json have ", len(idx), " images", m.Images)
		//fmt.Printf("idx len: %d\n", len(idx))
		for _, n := range idx {
			log.Println("urlbase ", n.Urlbase)
			url := fmt.Sprintf("http://www.bing.com%s", n.Urlbase)
			url += "_" + *resolution + ".jpg"
			str := strings.Split(url, "/")
			//fmt.Printf("split url by slash: %s \n", url)
			length := len(str)
			//fmt.Printf("slashed array: %s \n", str)
			imgeUrl := url
			//fmt.Printf("the URL: %s \n", imgeUrl)
			//get filename by "?"
			tmpfilename := strings.Split(str[length-1], "?")
			//fmt.Printf("the last field is filename: %s \n", tmpfilename)
			filename := tmpfilename[0]
			if filename == oldfilename {
				continue
			} else {
				oldfilename = filename
			}
			image := &image{url: imgeUrl, filename: filename}
			//fmt.Printf("start download %s\n", image.url)
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

					log.Println("start write file", image.filename)
					imgWriter := bufio.NewWriterSize(img, bufferSize)

					_, err = io.Copy(imgWriter, resp.Body)
					if err != nil {
						fmt.Print(err)
					}
					imgWriter.Flush()
					fmt.Printf("finish download %s\n", image.url)
					log.Println("finish download ", image.url)
					atomic.AddInt32(&ctx.okCounter, 1)
				}
			}
		}
	}
	log.Println("counter", ctx.okCounter)
}
