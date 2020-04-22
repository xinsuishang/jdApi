package trace

import (
	"bufio"
	"fmt"
	"io"
	"jdApi/api"
	"jdApi/conf"
	"log"
	"os"
	"time"
)

func Trace() {
	ts := time.Now()
	expressSlice := loadFile()
	elapsed1 := time.Since(ts)
	tm := time.Now()
	checkLen := len(expressSlice)
	nums := 128
	fmt.Println(checkLen)
	if checkLen > 0 {
		jobChan := make(chan string)
		resChan := make(chan string)

		conf.WG.Add(1)
		go getExpress(jobChan, expressSlice)
		conf.WG.Add(nums)
		go createPool(nums, jobChan, resChan)
		conf.WG.Add(1)
		go resDeal(resChan, checkLen)

		conf.WG.Wait()
	}
	elapsed2 := time.Since(tm)
	fmt.Println("elapsed1=", elapsed1, "elapsed2=", elapsed2)
}

func createPool(num int, jobChan chan string, resChan chan string) {

	// 根据开协程个数，去跑运行
	for i := 0; i < num; i++ {
		go trace(jobChan, resChan)
	}
}

func trace(jobChan, resChan chan string) {
	for express := range jobChan {
		traceRes, err := api.Api(express)
		if err != nil {
			log.Fatal(traceRes, err)
		}
		if traceRes.Response.Code == "0" {
			res := traceRes.Response.TraceApiDtos
			if len(res) > 1 {
				resChan <- fmt.Sprintln(express, res)
			} else {
				resChan <- fmt.Sprintln(express, []conf.TraceRoute{})
			}
		} else {
			conf.ErrResp.Store(express, traceRes.Response)
			resChan <- fmt.Sprintln(express, []conf.TraceRoute{})
		}
	}
	conf.WG.Done()
}

func getExpress(jobChan chan string, expressSlice []string) {
	defer conf.WG.Done()
	for i := 0; i < len(expressSlice); i++ {
		jobChan <- expressSlice[i]
	}
	close(jobChan)
}

func loadFile() (res []string) {

	file := "express.txt"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		res = append(res, string(a))
	}

	return
}

func resDeal(resChan chan string, checkLen int) {
	defer conf.WG.Done()
	i := 0
	for {

		r, ok := <-resChan
		if !ok {
			break
		} else {
			i++
			fmt.Println(i, r)
		}
		if checkLen == i {
			close(resChan)
		}
	}
}
