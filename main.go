package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var (
		accountID   string
		apiKey      string
		userAgent   string
		deleteUsers bool
	)
	flag.StringVar(&accountID, "id", "", "Your Drip Account ID. This is required.")
	flag.StringVar(&apiKey, "key", "", "Your Drip API Key. This is required.")
	flag.StringVar(&userAgent, "agent", "", "Eg: \"Your App Name (www.yourapp.com)\". This is optional.")
	flag.BoolVar(&deleteUsers, "delete", false, "Whether or not any subscribers found with temporary email addresses should be deleted. If left false, the app will print out the email addresses, if set to true it will also delete them from your Drip account.")
	flag.Parse()

	if accountID == "" || apiKey == "" {
		fmt.Println("id and key flags are required. use --help for more info.")
		os.Exit(1)
	}

	dc := DripClient{
		AccountID: accountID,
		APIKey:    apiKey,
		UserAgent: userAgent,
	}

	bl := blacklist()
	var toBoot []Subscriber

	page, numPages := 1, 1
	for page <= numPages {
		fmt.Printf("Getting page %d...\n", page)
		subs, meta, err := dc.Subscribers(page)
		if err != nil {
			panic(err)
		}
		numPages = meta.TotalPages
		if page == 1 {
			fmt.Println("Total pages:", numPages)
		}
		for _, sub := range subs {
			domain := strings.Split(sub.Email, "@")[1]
			if _, ok := bl[domain]; ok {
				toBoot = append(toBoot, sub)
			}
		}
		page++
	}

	fmt.Printf("Found %d total emails:\n", len(toBoot))
	for _, sub := range toBoot {
		fmt.Printf("  %s: %s\n", sub.ID, sub.Email)
	}
	// Quit here if we aren't deleting users.
	if !deleteUsers {
		return
	}

	// If you prefer to unsubscribe this is the code for it
	// for lo, hi := 0, 1000; lo < len(toBoot); lo, hi = hi, hi+1000 {
	// 	err := dc.Unsubscribe(toBoot[lo:min(hi, len(toBoot))])
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// Delete the bad emails
	for _, sub := range toBoot {
		err := dc.DeleteSub(sub)
		if err != nil {
			panic(err)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func blacklist() map[string]struct{} {
	f, err := os.Open("./disposable.txt")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(b), "\n")
	ret := make(map[string]struct{}, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ret[line] = struct{}{}
	}
	return ret
}
