package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/url"
	"sort"
	"strings"
)

type Credential struct {
	PublicKey  string
	PrivateKey string
}

func (c *Credential) CreateSign(query string) string {
	// replace "=" "&"
	str := strings.Replace(query, "&", "", -1)
	str = strings.Replace(str, "=", "", -1)

	// crypto by SHA1
	strUnescaped, _ := url.QueryUnescape(str)
	h := sha1.New()
	s := strUnescaped + c.PrivateKey
	io.WriteString(h, s)
	bs := h.Sum(nil)
	result := hex.EncodeToString(bs)

	return result
}

func (c *Credential) BuildCredentialedQuery(query map[string]string) string {
	var queryList []string
	for k, v := range query {
		queryList = append(queryList, k+"="+url.QueryEscape(v))
	}
	queryList = append(queryList, "PublicKey="+url.QueryEscape(c.PublicKey))
	sort.Strings(queryList)
	queryString := strings.Join(queryList, "&")

	sign := c.CreateSign(queryString)
	queryString = queryString + "&Signature=" + sign
	return queryString
}
