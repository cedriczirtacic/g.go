/**
 https://github.com/cedriczirtacic/g.go
**/
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// import html parsing capabilities
import "golang.org/x/net/html"

// default parameters
const (
	default_port = 8080
	default_file = "g.html"
	default_buff = 1024
	default_comm = false
)

// global vars
var (
	port         int
	ifile        string
	bfile        *os.File
	comm_disable bool
)

// how our bookmarks are going to be
type Link struct {
	link string
	name string
}

// i will store all your bookmarks
var bookmarks []Link

func init() {
	flag.IntVar(&port, "port", default_port, "Listening port.")
	flag.BoolVar(&comm_disable, "disable_cmds", default_comm, "Disable commands.")
	flag.StringVar(&ifile, "if", default_file, "Input file to parse.")

	// parse all params
	flag.Parse()

	// check if -h exists
	if f := flag.Lookup("h"); f != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	var err error

	// print runing info
	log.Printf("Reading bookmarks from: %s\n", ifile)
	log.Printf("Listening on port: %d\n", port)

	// bookmarks file exists?
	if _, err := os.Stat(ifile); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	// lets start reading the bookmarks
	bfile, err = os.Open(ifile)
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	//create list with name->links
	bookmarks = make([]Link, default_buff)
	err = load_bookmarks(&bookmarks, &bfile)
	if err != nil {
		log.Fatal(err)
		goto quit
	}

	// get ready to serve the bookmarks!
	http.ListenAndServe(
		fmt.Sprintf(":%d", port), http.HandlerFunc(handle_bookmarks))

quit:
	err = bfile.Close() // lets assume nothing fails here and we are happy
	os.Exit(0)
}

// this little fella is going to get and process the requests
func handle_bookmarks(w http.ResponseWriter, r *http.Request) {
	var bookmark string = html.UnescapeString(r.URL.Path[1:len(r.URL.Path)]) // unescape and erase the '/'

	// use querys for special commands
	if len(r.URL.RawQuery) > 0 {
		if !comm_disable {
			log.Printf("Got a command: %s", r.URL.RawQuery)
			// Everything is OK
			w.WriteHeader(http.StatusOK)

			switch r.URL.RawQuery {
			case "reload":
				var err error
				err = load_bookmarks(&bookmarks, &bfile)
				if err != nil {
					fmt.Fprintf(w, "%s", err)
				}
				break
			case "print":
				fmt.Fprintf(w, "Bookmarks:\n")
				for i, _ := range bookmarks {
					if bookmarks[i].name == "" {
						break
					}
					fmt.Fprintf(w, "\t(*) %s\n", bookmarks[i].name)
				}
				break
			}
		} else {
			fmt.Fprintf(w, "Commands disabled!\n")
		}

	} else {
		// there's no special commands?
		// lets look for the bookmark then
		if bookmark != "favicon.ico" {
			for i, _ := range bookmarks {
				if bookmark == bookmarks[i].name {
					log.Printf("Delivered bookmark for: %s", r.RemoteAddr)
					http.Redirect(w, r, bookmarks[i].link, 302)
					return
				}
			}
		}

		// respond with 503 if nothing found
		http.Error(w, "Bookmark doesn't exists!", 503)
		return
	}
}

// read and re-read bookmarks file
func load_bookmarks(b *[]Link, f **os.File) error {
	var linklen int = 0
	var read_text bool = false
	var tokenizer *html.Tokenizer

	tokenizer = html.NewTokenizer(*f)

	for {
		var token html.TokenType = tokenizer.Next()
		if token == html.ErrorToken {
			break
		}

		switch token {
		case html.TextToken:
			if read_text {
				var text []byte
				text = tokenizer.Raw()
				// we need to have a name for the bookmark
				if string(text) != "" {
					(*b)[linklen].name = string(text)
					linklen++
				}
			}

			read_text = false
			break
		case html.StartTagToken:
			var tname []byte
			var hasattr bool
			tname, hasattr = tokenizer.TagName()
			// we only need the bookmarks links
			if string(tname) == "a" && hasattr == true {
				for {
					attr, val, more := tokenizer.TagAttr()
					if string(attr) == "href" {
						(*b)[linklen] = Link{
							link: string(val),
						}
						read_text = true
					}
					if !more {
						break
					}
				}
			}
			break
		}
	}

	if len(bookmarks[0].name) < 1 {
		return errors.New("load_bookmarks() error: we got an empty bookmarks array")
	}

	// let's return to the begining of the file
	// in case it's a reload
	_, _ = (*f).Seek(0, 0)

	return nil
}
