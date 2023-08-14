package adyen_encrypt

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"

	jose "github.com/go-jose/go-jose/v3"
)

type Encryptor struct {
	Key       string
	OriginKey string
	Domain    string
	RsaPubKey *rsa.PublicKey
}

var (
	errNoKey = errors.New("key field is empty")
)

func NewEncryptor(key string) *Encryptor {
	return &Encryptor{}
}

func (enc *Encryptor) SetKey(key string) {
	enc.Key = key
}

func (enc *Encryptor) SetOriginKey(originkey string) {
	enc.OriginKey = originkey
}

func (enc *Encryptor) SetDomain(domain string) {
	enc.Domain = domain
}

func (enc *Encryptor) ParseKey() (err error) {
	if enc.Key == "" {
		return errNoKey
	}
	jwk := DefaultJWK()
	if err = jwk.ParseAdyenKey(enc.Key); err != nil {
		return
	}
	enc.RsaPubKey = jwk.JWKToPem()
	return
}

func PrepareEncryptor(key, originkey, domain string) (enc *Encryptor, err error) {
	if originkey == "" {
		originkey = "live_YCN5QJ4BXJHSTL24DUQMIHO4JQP2XDLK"
	}
	if domain == "" {
		domain = "https://www.bstn.com"
	}
	enc = &Encryptor{
		Key:       key,
		OriginKey: originkey,
		Domain:    domain,
	}
	err = enc.ParseKey()
	return
}

func (enc *Encryptor) EncryptSingle(data []byte) (ret string, err error) {
	rcpt := jose.Recipient{
		Algorithm: jose.KeyAlgorithm(jose.RSA_OAEP),
		Key:       enc.RsaPubKey,
	}

	opts := &jose.EncrypterOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"version": "1",
		},
	}

	jwe, err := jose.NewEncrypter(jose.ContentEncryption(jose.A256CBC_HS512), rcpt, opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	cipherText, err := jwe.Encrypt(data)
	if err != nil {
		return
	}
	return cipherText.CompactSerialize()
}

type AdyenData struct {
	EncryptedCardNumber   string `json:"encryptedCardNumber"`
	EncryptedExpiryMonth  string `json:"encryptedExpiryMonth"`
	EncryptedExpiryYear   string `json:"encryptedExpiryYear"`
	EncryptedSecurityCode string `json:"encryptedSecurityCode"`
}

func (enc *Encryptor) EncryptData(number, expMonth, expYear, cvc string) (ret *AdyenData, err error) {
	ret = &AdyenData{}

	ts := NowTimeISO()
	ref := "https://checkoutshopper-live.adyen.com/checkoutshopper/securedfields/" + enc.OriginKey + "/4.5.0/securedFields.html?type=card&d=" + base64.StdEncoding.EncodeToString([]byte(enc.Domain))
	number = FormatCardNumber(number)

	numberPayload := []byte(`{"number":"` + number + `","generationtime":"` + ts + `","numberBind":"1","activate":"3","referrer":"` + ref + `","numberFieldFocusCount":"3","numberFieldLog":"fo@44070,cl@44071,KN@44082,fo@44324,cl@44325,cl@44333,KN@44346,KN@44347,KN@44348,KN@44350,KN@44351,KN@44353,KN@44354,KN@44355,KN@44356,KN@44358,fo@44431,cl@44432,KN@44434,KN@44436,KN@44438,KN@44440,KN@44440","numberFieldClickCount":"4","numberFieldKeyCount":"16"}`)
	expMonthPayload := []byte(`{"expiryMonth":"` + expMonth + `","generationtime":"` + ts + `"}`)
	expYearPayload := []byte(`{"expiryYear":"` + expYear + `","generationtime":"` + ts + `"}`)
	cvcPayload := []byte(`{"cvc":"` + cvc + `","generationtime":"` + ts + `","cvcBind":"1","activate":"4","referrer":"` + ref + `","cvcFieldFocusCount":"4","cvcFieldLog":"fo@122,cl@123,KN@136,KN@138,KN@140,fo@11204,cl@11205,ch@11221,bl@11221,fo@33384,bl@33384,fo@50318,cl@50319,cl@50321,KN@50334,KN@50336,KN@50336","cvcFieldClickCount":"4","cvcFieldKeyCount":"6","cvcFieldChangeCount":"1","cvcFieldBlurCount":"2","deactivate":"2"}`)

	ret.EncryptedCardNumber, err = enc.EncryptSingle(numberPayload)
	if err != nil {
		return
	}
	ret.EncryptedExpiryMonth, err = enc.EncryptSingle(expMonthPayload)
	if err != nil {
		return
	}
	ret.EncryptedExpiryYear, err = enc.EncryptSingle(expYearPayload)
	if err != nil {
		return
	}
	ret.EncryptedSecurityCode, err = enc.EncryptSingle(cvcPayload)
	if err != nil {
		return
	}
	return
}
