package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
    "io/ioutil"
	"github.com/tomnomnom/gahttp"
	"golang.org/x/net/html"
	"github.com/fatih/color"
)

func status(stat int) string {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	if stat == 200 {
		return green(stat)
	} else if (stat >= 300) && (stat <= 399) {
		return blue(stat)
	} else if (stat >= 400) && (stat <= 499) {
		return yellow(stat)
	} else if (stat >= 500) && (stat <= 599) {
		return red(stat)
	} else {
		return white(stat)
	}
}

func extractTitle(req *http.Request, resp *http.Response, err error) {
	if err != nil {
		return
	}

	z := html.NewTokenizer(resp.Body)
	for {
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		tt := z.Next()
		if tt == html.ErrorToken {
			bbb, _ := ioutil.ReadAll(resp.Body)
			bss := string(bbb)
			fmt.Printf("[ %s ] No Title (%s) [ %s ]\n", status(resp.StatusCode), yellow(req.URL), red(len(bss)))
			break
		}

		t := z.Token()

		if t.Type == html.StartTagToken {
			if t.Data == "title" {
				if z.Next() == html.TextToken {
					//fmt.Printf(z.Token().Data)
					title := strings.TrimSpace(z.Token().Data)

					bb, _ := ioutil.ReadAll(resp.Body)
					bs := string(bb)
					if len([]rune(title)) != 0 && title != "" {
						fmt.Printf("[ %s ] %s (%s) [ %s ]\n", status(resp.StatusCode), cyan(title), yellow(req.URL), red(len(bs)))
					} else {
						fmt.Printf("[ %s ] Blank Title (%s) [ %s ]\n", status(resp.StatusCode), yellow(req.URL), red(len(bs)))
					}				
					break
				}
			} 
		}
	}
}

func main() {

	var concurrency = 20
	flag.IntVar(&concurrency, "c", 20, "Concurrency")
	fr := flag.Bool("follow-redir", false, "Follow Redirect")
	flag.Parse()

	p := gahttp.NewPipelineWithClient(gahttp.NewClient(gahttp.SkipVerify))
	if !*fr {
		p.SetClient(gahttp.NewClient(gahttp.NoRedirects))
	}
	p.SetConcurrency(concurrency)
	extractFn := gahttp.Wrap(extractTitle, gahttp.CloseBody)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		p.Get(sc.Text(), extractFn)
	}
	p.Done()

	p.Wait()

}
