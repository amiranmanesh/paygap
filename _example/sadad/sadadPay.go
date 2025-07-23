package main

import (
	"context"
	"fmt"
	"github.com/amiranmanesh/paygap/client"
	"github.com/amiranmanesh/paygap/providers/sadad"
	"log"
	"math/rand"
)

func main() {

	c := client.New()
	s, err := sadad.New(c, "9001", "RecivedMerchenantkey", "1565879")
	if err != nil {
		log.Fatal(err)
	}
	//pay
	//اگر پرداخت به صورت تسهیمی است باید آبجکت مولتی پلکسینگ نیز مقدار دهی شود
	orderId := string(rand.Int())
	resp, err := s.PaymentRequest(context.Background(), 50000, orderId, "returnUrl", false, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)

	//verify
	v_res, v_err := s.VerifyRequest(context.Background(), *resp)
	if v_err != nil {
		log.Fatal(v_err)
	}

	fmt.Println(v_res)
}
