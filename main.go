package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	//  "os"
	"time"
	//"image"
	"io/ioutil"
	"math/rand"
	//"net/http"
	//"time"
	"strconv"
)

var numberCreatedImages int = 0
var numberGoroutines int = 0
var toLoad bool = false

type MathFunction func(int) int

func GenerateReq() string {
	x := rand.Intn(399) + 1
	y := rand.Intn(399) + 1
	w := rand.Intn(399) + 1
	h := rand.Intn(399) + 1
	x = x - 0
	y = y - 0
	w = w - 0
	h = h - 0
	//write random choose transformation
	var choose int = rand.Intn(4)
	t := "/process/"
	switch choose {
	case 0:
		t = t + "rotate-cw/"
	case 1:
		t = t + "rotate-ccw/"
	case 2:
		t = t + "flip-v/"
	case 3:
		t = t + "flip-h/"
	}
	t = t + strconv.Itoa(x)
	t = t + ","
	t = t + strconv.Itoa(y)
	t = t + ","
	t = t + strconv.Itoa(w)
	t = t + ","
	t = t + strconv.Itoa(h)
	return t
}

func ReadFile(path string) []byte {
	img, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Can not open file!\n")
	}
	return img
}

func Client(url string, img []byte, n int) {
	numberGoroutines++

	address2 := GenerateReq()
	address := url + address2
	reader := bytes.NewReader(img)
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	//fmt.Println(address)
	req, err := http.NewRequest("POST", address, reader)
	if err != nil {
		log.Println(err)
		log.Println("Oopes")
	}
	//fmt.Println("Req: ", req)
	nBytes := len(img)
	//fmt.Printf("Len is %d", nBytes)
	req.Header.Add("Content-Length", strconv.Itoa(nBytes))
	log.Printf("Routine %d opened\n", n)
	response, err := client.Do(req)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Routine number %d has Timeout", n)
			log.Printf("---------------------------------------------------------------------")
			log.Println(*req)
			panic("TIMEOUT!")
		}
	}()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Got it")
	}
	b := make([]byte, 10000)
	//strr := make([]string, 10000)
	//nh, errrs := response.Header.Get(strr)
	_, errs := response.Body.Read(b)
	if errs != nil {
		log.Println("Bleeeeeeeet  -------------------------------Can not read response!\n")
	} else {
		if toLoad {
			os.Create(".\\got_img" + strconv.Itoa(numberCreatedImages) + ".png")
			err := ioutil.WriteFile("D:\\golang\\ddos\\got_img"+strconv.Itoa(numberCreatedImages)+".png", b, 0644)
			numberCreatedImages++
		}
		if err != nil {
			panic("Can not write image!\n")
		}
		//a := <-num
		//a++
		//num <- a
	}
	//log.Println("Head is ", strr)
	//log.Println("Response is ", b)
	//log.Printf(". Bytes: %d\n\n", nb)
	defer response.Body.Close()
	defer func() {
		numberGoroutines--
	}()
	log.Printf("Routine %d closed\n", n)
}

func SendNum(n int, img []byte) {
	for i := 0; i < n; i++ {
		go Client("http://192.168.1.2:8080", img, i)
		log.Printf("Gorutine number %d is started\n", i)
		log.Printf("--------------------------Opened gourutines: %d\n", numberGoroutines)
	}
}

func SendLinear(n int, img []byte) {
	lim := 1
	nums := 0
	for i := 0; i < n; i++ {
		for j := 0; j < lim; j++ {
			go Client("http://192.168.1.2:8080", img, nums)
			log.Printf("Gorutine number %d is started\n", nums)
			nums++
		}
		log.Printf("--------------------------Opened gourutines: %d\nLim: %d\n", numberGoroutines, lim)
		time.Sleep(1 * time.Second)
		lim = lim + 1
	}
}

func Parabolla(n int) int {
	return (n - 2) * (n - 2)
}

func SendFunc(n int, f MathFunction, ms time.Duration, img []byte) {
	nums := 0
	for i := 0; i < n; i++ {
		for j := 0; j < f(n); j++ {
			go Client("http://192.168.1.2:8080", img, nums)
			log.Printf("Gorutine number %d is started\n", nums)
			nums++
		}
		log.Printf("--------------------------Opened gourutines: %d\n", numberGoroutines)
		time.Sleep(ms * time.Millisecond)
	}
}

func main() {
	//numberClosedGorutines := make(chan int)
	img := ReadFile("D:\\golang\\ddos\\zebra.png")
	//SendNum(60, img)
	SendFunc(100, Parabolla, 300, img)
	//log.Println("Got it!")
	var a string
	fmt.Scanln(a)
}
