# Adyen-Encrypt-4.5.0

Simple golang package to encrypt sensitive card data using Adyen's v4.5.0 encryption method.

## Installation

```
go get github.com/AlienRecall/adyen-encrypt-4.5
```

## Usage

Explained Usage:
```go
enc, _ := adyen_encrypt.PrepareEncryptor(
	"adyen key", //required
	"adyen origin key", //optional
	"site domain", //optional
)
encryptedData, err := enc.EncryptData("card number", "card expire month", "card expire year", "card cvc")
```

Single field encryption:
```go
package main

import (
	"fmt"
	"time"

	adyen_encrypt "github.com/AlienRecall/adyen-encrypt-4.5"
)

func main() {
	enc, _ := adyen_encrypt.PrepareEncryptor(
		"10001|C621C7E8267CF5A0758EC2E0530AF2B59625EFA2A26174690B401476BA5FF1AD079D881838CD625384D546DAB4E82CF1E414F1F2C7EB5420AFD9F8FF516479FD2F7EDA66572BB9C08672961C8BF528FFD0B1951B29C2332FBF301A96BA1D41DA28F39718095222C4CCFF0C0BCAECDEF944D2994D45FB81FE210090B46E5BE22CCCBAC4F413C08F90229D0E9096046BDB6745E5C549A7FEDC907646661C79A0A14ECE4EA351A07832D7228AA8D3398874D173076E475196E1DFBF35E0FDA83C047DED0156D6839D67DF1DC0D00509E8876DF209169832607B3FAE834F0DD8E78123A991E50EFD485740622FBE3EAAE6FA33BEE2DDA42465DA36D468500AF7BD01",
		"live_YCN5QJ4BXJHSTL24DUQMIHO4JQP2XDLK",
		"https://www.bstn.com",
	)
	ts := adyen_encrypt.NowTimeISO()
	data := []byte(`{"expiryMonth":"12","generationtime":"` + ts + `"}`)
	encrypted, err := enc.EncryptSingle(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(encrypted)
}
```

Card data encryption:
```go
package main

import (
	"fmt"
	"time"

	adyen_encrypt "github.com/AlienRecall/adyen-encrypt-4.5"
)

func main() {
	enc, _ := adyen_encrypt.PrepareEncryptor(
		"10001|C621C7E8267CF5A0758EC2E0530AF2B59625EFA2A26174690B401476BA5FF1AD079D881838CD625384D546DAB4E82CF1E414F1F2C7EB5420AFD9F8FF516479FD2F7EDA66572BB9C08672961C8BF528FFD0B1951B29C2332FBF301A96BA1D41DA28F39718095222C4CCFF0C0BCAECDEF944D2994D45FB81FE210090B46E5BE22CCCBAC4F413C08F90229D0E9096046BDB6745E5C549A7FEDC907646661C79A0A14ECE4EA351A07832D7228AA8D3398874D173076E475196E1DFBF35E0FDA83C047DED0156D6839D67DF1DC0D00509E8876DF209169832607B3FAE834F0DD8E78123A991E50EFD485740622FBE3EAAE6FA33BEE2DDA42465DA36D468500AF7BD01",
		"live_YCN5QJ4BXJHSTL24DUQMIHO4JQP2XDLK",
		"https://www.bstn.com",
	)
	encryptedData, err := enc.EncryptData("4242424242424242", "12", "2023", "123")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("encryptedCardNumber: %s\nencryptedExpiryMonth: %s\nencryptedExpiryYear: %s\nencryptedSecurityCode: %s\n", encryptedData.EncryptedCardNumber, encryptedData.EncryptedExpiryMonth, encryptedData.EncryptedExpiryYear, encryptedData.EncryptedSecurityCode)
}
```

Please note that you need to replace `adyen key` with the actual website's Adyen key.

## Contributing

We appreciate any contributions you might make. Please feel free to submit a pull request, issue, or suggestion.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more information.
