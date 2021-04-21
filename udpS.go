package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	PORT, _ := strconv.Atoi(arguments[1])
	CLIENTS := arguments[2]
	FILE := arguments[3]

	cli, err := strconv.Atoi(CLIENTS)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("file-" + FILE + ".txt")
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	hash := h.Sum(nil)

	defer f.Close()

	for i := 1; i <= cli; i++ {
		s, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(PORT+i))
		if err != nil {
			fmt.Println(err)
		}
		connection, err := net.ListenUDP("udp4", s)
		if err != nil {
			fmt.Println(err)
		}
		defer connection.Close()
		go manageConnectionSend(connection, FILE, hash)
		time.Sleep(4 * time.Second)
	}
	for {

	}
}
func manageConnectionSend(s *net.UDPConn, i string, hash []byte) {
	buffer := make([]byte, 1024)
	boolean := true

	_, addr, err := s.ReadFromUDP(buffer)

	s.WriteToUDP(hash, addr)
	if err != nil {
		fmt.Println(err)
	}
	f1, err := os.Open("file-" + i + ".txt")
	reader := bufio.NewReader(f1)

	paquetes := 0

	start := time.Now()
	for boolean {
		n, err := reader.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}
		_, err = s.WriteToUDP(buffer[:n], addr)
		if err != nil {
			fmt.Println(err)
		}
		if n < 1024 {
			boolean = false
		}
		paquetes++
		time.Sleep(25)
		//fmt.Println(n)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	nombre := fmt.Sprintf("logs/Server-%d-%02d-%02d-%02d-%02d-%02d-log.txt", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fil, err := os.OpenFile(nombre, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err)
	}
	_, err2 := fil.WriteString("Se envia el archivo: file-" + i + ".txt")
	if err2 != nil {
		fmt.Println(err2)
	}
	fi, err2 := f1.Stat()
	if err2 != nil {
		fmt.Println(err2)
	}
	_, err2 = fil.WriteString("\nDe tamaño: " + strconv.FormatInt(fi.Size(), 10))
	if err2 != nil {
		fmt.Println(err2)
	}
	_, err2 = fil.WriteString("\nArchivo entregado: si")
	if err2 != nil {
		fmt.Println(err2)
	}
	_, err2 = fil.WriteString("\nTiempo: " + elapsed.String())
	if err2 != nil {
		fmt.Println(err2)
	}
	_, err2 = fil.WriteString("\nCantidad de paquetes: " + strconv.Itoa(paquetes))
	if err2 != nil {
		fmt.Println(err2)
	}
	_, err2 = fil.WriteString("\nBytes enviados: " + strconv.FormatInt(fi.Size(), 10))
	if err2 != nil {
		fmt.Println(err2)
	}
}
