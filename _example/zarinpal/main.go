package main

import (
	"context"
	"fmt"
	"github.com/amiranmanesh/paygap/client"
	"github.com/amiranmanesh/paygap/providers/zarinpal"
	"log"
)

func main() {
	c := client.New()
	z, err := zarinpal.New(c, "YOUR_MERCHANT_ID", false)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := z.RequestPayment(context.Background(), 1000, "YOUR_CALL_BACK", "YOUR_CURRENCY", "description", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
