package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"sync"
	"net"
)

func main() {
	app := cli.NewApp()
	app.Name = "snooper"
	app.Version = "1.0.0"
	app.Usage = "Snooper is ready to run"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "James Brown",
			Email: "jbrown@invoca.com",
		},
		cli.Author{
			Name:  "Christian Parkinson",
			Email: "cparkinson@invoca.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filename",
			Value: "iplist",
			Usage: "File containing IP's. 1 per line",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 80,
			Usage: "Port to send requests to",
		},
		cli.StringFlag{
			Name:  "urlPath",
			Value: "/is_alive",
			Usage: "Path for the URL",
		},
		cli.IntFlag{
			Name:  "concurrency",
			Value: 10,
			Usage: "Number of concurrent requests",
		},
	}

	app.Action = func(c *cli.Context) error {
		filename := c.String("filename")
		port := c.Int("port")
		urlPath := c.String("urlPath")
		concurrency := c.Int("concurrency")

		if c.NArg() > 0 {
		}

		ips := make(chan string)
		output := make(chan result)
		done := make(chan struct{})

		go publisher(filename, ips)
		go writer(output, done)

		var wg sync.WaitGroup

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				consumer(port, urlPath, ips, output)
			}()
		}

		wg.Wait()
		close(output)
		<- done

		return nil
	}

	app.Run(os.Args)
}

func publisher(filename string, ips chan string) {
	fileHandle, err := os.Open(filename)
	defer fileHandle.Close()
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(fileHandle)
	for fileScanner.Scan() {
		ip := fileScanner.Text()
		ips <- ip
	}
	close(ips)
}

func consumer(port int, urlPath string, ips chan string, output chan result) {
	timeout := time.Duration(1 * time.Minute)
	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 2 * time.Second)
			},
		},
	}
	for {
		ip, ok := <-ips
		if !ok {
			break
		}
		url := fmt.Sprintf("http://%s:%d%s", ip, port, urlPath)
		resp, err := client.Get(url)
		if err != nil {
			output <- result{ip: ip, response: fmt.Sprintf("%s\n",err.Error())}
			continue
		}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				output <- result{ip: ip, response: err.Error()}
				continue
			}
			bodyString := string(bodyBytes)
			output <- result{ip: ip, response: bodyString}
		} else {
			myoutput := fmt.Sprintf("Received HTTP response code: %d\n", resp.StatusCode)
			output <- result{ip: ip, response: myoutput}
		}
	}
}

func writer(output chan result, done chan struct{}) {
	for {
		r, ok := <-output
		if !ok {
			done <- struct{}{}
			break
		}
		filename := fmt.Sprintf("%s.log", r.ip)
		bodyBytes := []byte(r.response)
		err := ioutil.WriteFile(filename, bodyBytes, 0644)
		if err != nil {
			panic(err)
		}
	}
}

type result struct {
	ip       string
	response string
}
