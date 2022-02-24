/*
	2022 (c) Yariya
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const threadDelay = time.Microsecond * 200
const socketTimeout = time.Millisecond * 500

type Protocol struct {
	Port    int
	Payload []byte
}

var mutex sync.Mutex
var protocols = map[string]*Protocol{
	"dns":  {Port: 53, Payload: []byte("\\x9c\\x88\\x01\\x20\\x00\\x01\\x00\\x00\\x00\\x00\\x00\\x01\\x07\\x75\\x63\\x64\\x61\\x76\\x69\\x73\\x03\\x65\\x64\\x75\\x00\\x00\\xff\\x00\\x01\\x00\\x00\\x29\\x10\\x00\\x00\\x00\\x80\\x00\\x00\\x00")},
	"dvr":  {Port: 37810, Payload: []byte("\\x44\\x48\\x49\\x50")},
	"ntp":  {Port: 123, Payload: []byte("\x17\x00\x03\x2a\x00\x00\x00\x00")},
	"snmp": {Port: 161, Payload: []byte("\\x30\\x20\\x02\\x01\\x01\\x04\\x06\\x70\\x75\\x62\\x6c\\x69\\x63\\xa5\\x13\\x02\\x02\\x00\\x01\\x02\\x01\\x00\\x02\\x01\\x46\\x30\\x07\\x30\\x05\\x06\\x01\\x28\\x05\\x00")},
	"wsd":  {Port: 3702, Payload: []byte("\\x3c\\x3a\\x2f\\x3e")},
}

var working = make(map[string]interface{})

func main() {

	if len(os.Args) < 5 {
		fmt.Printf("[#] AmpFilter\n[!] Usage: %s <input> <output> <protocol> <bytes>\n", os.Args[0])
		return
	}

	protocol := protocols[os.Args[3]]
	bytes, _ := strconv.Atoi(os.Args[4])

	if protocol == nil {
		log.Fatalln("[-] Protocol does not exist!")
		return
	}

	inputIO, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln("[-] This file does not exist!")
		return
	}
	input, err := ioutil.ReadAll(inputIO)
	if err != nil {
		log.Fatalln("[-] util Error ", err)
		return
	}

	in := strings.Split(string(input), "\n")

	var wg sync.WaitGroup
	for _, ip := range in {
		go func() {
			wg.Add(1)
			defer wg.Done()

			c, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: protocol.Port, IP: net.ParseIP(ip)})
			if err != nil {
				return
			}
			c.SetReadDeadline(time.Now().Add(socketTimeout))

			defer c.Close()
			_, err = c.Write(protocol.Payload)
			if err != nil {
				return
			}
			b := make([]byte, 2<<15)
			n, addr, err := c.ReadFrom(b)
			if err != nil {
				return
			}

			if n >= bytes {
				mutex.Lock()
				if working[ip] == nil {
					fmt.Printf("[+] Received working server %s bytes: %d\n", addr.String(), n)
					working[ip] = ip
				}
				mutex.Unlock()
			}
		}()
		time.Sleep(threadDelay)
	}
	time.Sleep(time.Second)
	wg.Wait()
	f, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalln("[-] Couldn't create file!")
		return
	}
	var allIps string
	for ip, _ := range working {
		allIps += ip + "\n"
	}
	_, err = f.Write([]byte(allIps))
	if err != nil {
		log.Fatalln("[-] Couldn't write to file!")
		return
	}
	fmt.Printf("[+] Filter saved to output file with %d working servers!\n", len(working))

}
