package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/victorlpgazolli/lightdns"

	"github.com/miekg/dns"
)

var records = map[string]string{
	"mail.amazon.com":  "192.162.1.2",
	"paste.amazon.com": "191.165.0.3",
}
var blackListDomains = []string{}

func lookupFunc(domain string) (string, error) {
	for _, invalidDomain := range blackListDomains {
		if invalidDomain == domain {
			log.Println("blocked:", invalidDomain)
			return "0.0.0.0", nil
		}
	}
	var msg dns.Msg
	var ip string
	fqdn := dns.Fqdn(domain)
	msg.SetQuestion(fqdn, dns.TypeA)
	in, err := dns.Exchange(&msg, "1.1.1.1:53")
	if err != nil {
		panic(err)
	}
	if len(in.Answer) < 1 {
		log.Println("No records")
		return "0.0.0.0", errors.New("no records")
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ip = a.A.String()
		}
	}
	log.Printf("%s => %s", domain, ip)
	return ip, nil
}
func writeDomainsToFile() {
	resp, err := http.Get("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	out, err := os.Create("hosts.txt")
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}
func getUpdatedAdsDomains() {

	if _, err := os.Open("./hosts.txt"); err != nil {
		writeDomainsToFile()
	}
	invalidDomain := "0.0.0.0 0.0.0.0"
	notADomain := "# "
	domainIndicative := "0.0.0.0 "
	file, err := os.Open("./hosts.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		domainLine := scanner.Text()
		lineIsNotADomain := strings.HasPrefix(domainLine, notADomain)
		domainIsInvalid := strings.HasPrefix(domainLine, invalidDomain)
		startsWithCorrectIp := strings.HasPrefix(domainLine, domainIndicative)
		hasToIncludeDomain := !lineIsNotADomain && !domainIsInvalid && startsWithCorrectIp
		if !hasToIncludeDomain {
			continue
		}
		blackListDomains = append(blackListDomains, strings.Trim(domainLine, domainIndicative))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}



func main() {

	getUpdatedAdsDomains()

	portStr := os.Getenv("PORT")
	port, _ := strconv.Atoi(portStr)

	if hasPort := port != 0; !hasPort {
		port = 53
		portStr = "53"
	}
	dnsServer := lightdns.NewDNSServer(port)


	dnsServer.AddZoneData(".", nil, lookupFunc, lightdns.DNSForwardLookupZone)

	log.Printf("dns server is starting on port: %v", port)

	dnsServer.StartAndServe()
}
