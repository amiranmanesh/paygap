# PayGap [![Go Reference](https://pkg.go.dev/badge/github.com/amiranmanesh/paygap.svg)](https://pkg.go.dev/github.com/amiranmanesh/paygap)
Payment gateway providers SDK Go for zarinpal, idpay, pay.ir, and other gateway providers.

![paygap](assets/banner.jpg)

## Install For ![Go Version](https://img.shields.io/badge/go%20version-%3E=1.19-61CFDD.svg?style=flat-square)

```shell
go get -u github.com/amiranmanesh/paygap
```

## Usage

example for zarinpal provider ([Other examples](https://github.com/amiranmanesh/paygap/tree/main/_example))

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/amiranmanesh/paygap/client"
	"github.com/amiranmanesh/paygap/providers/zarinpal"
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
```

## Contributing

1. fork project in your GitHub account.
2. create new branch for new changes.
3. after change code, send Pull Request.

## TODO

- [x] zarinpal
- [x] idpay
- [x] pay.ir
- [x] payping
- [x] vandar.io
- [ ] rayanpay
- [ ] nextpay
- [x] mellat
- [ ] parsian
- [ ] pasargad
- [x] sadad
- [x] [zibal](https://zibal.ir)
