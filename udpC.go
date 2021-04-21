package main

import (
	"bufio"
	"bytes"
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
		fmt.Println("Please provide a host:port string")
		return
	}
	CONNECT := arguments[1]
	CLIENTS := arguments[2]
	PORT, _ := strconv.Atoi(arguments[3])

	cli, err := strconv.Atoi(CLIENTS)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 1; i <= cli; i++ {
		s, err := net.ResolveUDPAddr("udp4", CONNECT+":"+strconv.Itoa(PORT+i))
		if err != nil {
			fmt.Println(err)
			return
		}
		c, err := net.DialUDP("udp4", nil, s)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()
		go manageConnection(c, strconv.Itoa(i), cli)
		time.Sleep(4 * time.Second)
	}
	for {

	}
}
func manageConnection(s *net.UDPConn, i string, clientes int) {

	boolean := true

	f, err := os.OpenFile("ArchivosRecibidos/Cliente"+i+"-Prueba-"+strconv.Itoa(clientes)+".txt", os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := []byte("Listo\n")
	_, err = s.Write(data)
	if err != nil {
		fmt.Println(err)
	}

	hash := make([]byte, 16)
	_, _, err = s.ReadFromUDP(hash)
	if err != nil {
		fmt.Println(err)
	}
	
	writer := bufio.NewWriter(f)

	paquetes := 0

	start := time.Now()
	kill := time.Now().Add(24 * time.Second)
	kill2 := time.Now().Add(240 * time.Second)
	s.SetDeadline(kill)
	for boolean && time.Now().Before(kill2){
		buffer := make([]byte, 1024)
		n, _, err := s.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}
		nn, err := writer.Write(buffer[:n])
		if err != nil {
			fmt.Println(err)
			break
		}
		if n != 1024 || nn != 1024{
			boolean = false
		}
		paquetes++
		time.Sleep(22)
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now()
	f.Close()
	f1, err := os.Open("ArchivosRecibidos/Cliente" + i + "-Prueba-" + strconv.Itoa(clientes) + ".txt")
	h := md5.New()
	if _, err := io.Copy(h, f1); err != nil {
		log.Fatal(err)
	}
	hashCalc := h.Sum(nil)

	elapsed := t.Sub(start)
	nombre := fmt.Sprintf("logs/Client-%d-%02d-%02d-%02d-%02d-%02d-log.txt", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	if err != nil {
		log.Fatal(err)
	}
	fil, err := os.OpenFile(nombre, os.O_CREATE, 0755)
	_, err2 := fil.WriteString("Se recibe el archivo: ArchivosRecibidos/Cliente" + i + "-Prueba-" + strconv.Itoa(clientes) + ".txt")
	if err2 != nil {
		log.Fatal(err2)
	}
	fi, err2 := f1.Stat()
	_, err2 = fil.WriteString("\nDe tamaño: " + strconv.FormatInt(fi.Size(), 10))
	str := ""
	if bytes.Compare(hash, hashCalc) == 0 {
		str = "si"
	} else {
		str = "no"
	}
	if err2 != nil {
		log.Fatal(err2)
	}
	_, err2 = fil.WriteString("\nArchivo entregado correctamente: " + str)
	if err2 != nil {
		log.Fatal(err2)
	}
	_, err2 = fil.WriteString("\nTiempo: " + elapsed.String())
	if err2 != nil {
		log.Fatal(err2)
	}
	_, err2 = fil.WriteString("\nCantidad de paquetes: " + strconv.Itoa(paquetes))
	if err2 != nil {
		log.Fatal(err2)
	}
	_, err2 = fil.WriteString("\nBytes enviados: " + strconv.FormatInt(fi.Size(), 10))
	if err2 != nil {
		log.Fatal(err2)
	}
}
