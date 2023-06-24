package main

import (
	//snmp "github.com/tiebingzhang/WapSNMP"
	"log"
	"math/rand"
	"net"
	"time"

	snmp "github.com/ales999/go-snmp"
)

func myUDPServer(listenIPAddr string, port int) *net.UDPConn {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(listenIPAddr),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Printf("udp Listen error.")
		panic(err)
	}
	return conn
}

func main() {
	rand.Seed(0)
	target := ""
	community := ""
	version := snmp.SNMPv2c

	udpsock := myUDPServer("0.0.0.0", 162)
	defer udpsock.Close()

	wsnmp := snmp.NewWapSNMPOnConn(target, community, version, 2*time.Second, 5, udpsock)
	defer wsnmp.Close()

	wsnmp.Trapusers = append(wsnmp.Trapusers,
		snmp.V3user{User: "pcb.snmpv3", AuthAlg: "SHA1", AuthPwd: "this_is_my_pcb", PrivAlg: "AES", PrivPwd: "my_pcb_is_4_me"})

	packet := make([]byte, 3000)
	for {
		_, addr, err := udpsock.ReadFromUDP(packet)
		if err != nil {
			log.Fatal("udp read error\n")
		}

		log.Printf("Received trap from %s:\n", addr.IP)

		err = wsnmp.ParseTrap(packet)
		if err != nil {
			log.Printf("Error processing trap: %v.", err)
		}
	}
	//udpsock.Close()

}
