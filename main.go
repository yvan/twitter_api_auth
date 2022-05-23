package main

import (
	"os"
	"fmt"
	"time"
	"sort"
	"strings"
	"strconv"
	"net/url"
	"net/http"
	"math/rand"
	"io/ioutil"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/gorilla/mux"
)

var (
	TOKEN = os.Getenv("TW_TOKEN")
	TOKENSECRET = os.Getenv("TW_TOKENSECRET")
	CONSUMERKEY = os.Getenv("TW_CONSUMERKEY")
	CONSUMERSECRET = os.Getenv("TW_CONSUMERSECRET")
)

// take the number of bytes
// and product a random string
// of length nChars
func makeNonce(nBytes int) string{
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, nBytes)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// this function makes the key
// so that requests can be signed
func makeKey(aSecret, bSecret string) string{
	return strings.Join([]string{aSecret, bSecret,}, "&")
}

// this function takes a string
// and signs it using hmac/sha1
func makeSignature(key, baseStr string) string{
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(baseStr))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// this function helps us build
// a properly formatted authentication
// header for our request
func makeAuthHeader(signingKey, method, queryUrl string, oauthParams map[string]string) string{
	// get the header fields
	auth := url.Values{}
	for k,v := range oauthParams {
		auth.Add(k,v)
	}
	authParams := strings.Replace(auth.Encode(), "+", "%20", -1)
	baseStr := strings.ToUpper(method)+"&"+url.QueryEscape(queryUrl)+"&"+url.QueryEscape(authParams)
	auth.Add("oauth_signature", makeSignature(signingKey, baseStr))

	// sort keys of the map
	keys := make([]string, 0, len(oauthParams)+1)
	for k := range auth {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	header := "OAuth "
	for i, k := range keys {
		header += url.QueryEscape(k)+"="+"\""+url.QueryEscape(auth.Get(k))+"\""
		if i < len(keys)-1 {
			header += ","
		}
	}
	return header
}

func makeAuthRequest(method, queryUrl, header string) (int, string){
	client := &http.Client{}
	req, _ := http.NewRequest(strings.ToUpper(method), queryUrl, nil)
	req.Header = http.Header{
		"Authorization": []string{header},
	}
	res, _ := client.Do(req)

	if res.StatusCode != 200 {
		fmt.Println("err")
	}
	body, _ := ioutil.ReadAll(res.Body)
	return res.StatusCode, string(body)
}

func makeTokenRequest(method, queryUrl, header, data string) (int, string){
	client := &http.Client{}
	payload := url.Values{}
	payload.Set("oauth_verifier", data)
	req, _ := http.NewRequest(strings.ToUpper(method), queryUrl, strings.NewReader(payload.Encode()))
	req.Header = http.Header{
		"Authorization": []string{header},
		"Content-Legnth": []string{"57"},
		"Content-Type": []string{"application/x-www-form-urlencoded"},
	}

	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	return res.StatusCode, string(body)
}

func parseResponse(res string) map[string]string{
	ret := make(map[string]string)
	for _, s := range strings.Split(res, "&") {
		split := strings.Split(s, "=")
		
		ret[split[0]] = split[1]
	}
	return ret
}

func makeRedirect(token string) string{
	return fmt.Sprintf("https://api.twitter.com/oauth/authenticate?oauth_token=%s", token)
}


// step 1 of three legged oauth flow
func auth(w http.ResponseWriter, r *http.Request){
	var method string = "POST"
	var queryUrl string = "https://api.twitter.com/oauth/request_token"
	var authParams = make(map[string]string)
	var key = makeKey(CONSUMERSECRET,TOKENSECRET)
	authParams["oauth_callback"] = "http://localhost:8080/callback"
	authParams["oauth_token"] = TOKEN
	authParams["oauth_consumer_key"] = CONSUMERKEY
	authParams["oauth_nonce"] = makeNonce(32)
	authParams["oauth_signature_method"] = "HMAC-SHA1"
	authParams["oauth_timestamp"] = strconv.Itoa(int(time.Now().Unix()))
	authParams["oauth_version"] = "1.0"
	header := makeAuthHeader(key, method, queryUrl, authParams)
	
	fmt.Println("step 1 oauthflow: ")
	fmt.Printf("request header: %s", header)
	fmt.Println()
	code, body := makeAuthRequest("POST", "https://api.twitter.com/oauth/request_token", header)
	fmt.Printf("code: %d body: %s", code, body)
	fmt.Println()

	oauthInfo := parseResponse(body)

	// step 2 of three leg flow
	redirect := fmt.Sprintf("https://api.twitter.com/oauth/authenticate?oauth_token=%s", oauthInfo["oauth_token"])
	fmt.Printf("redirecting to %s", redirect)
	http.Redirect(w, r, redirect, 302)
	return
}

// step 3 of three leg flow
func callback(w http.ResponseWriter, r *http.Request){
	fmt.Println()
	fmt.Println("starting step 3")
	method := "POST"
	qParams := r.URL.Query()
	queryUrl := fmt.Sprintf("https://api.twitter.com/oauth/access_token")
	token := qParams.Get("oauth_token")

	var authParams = make(map[string]string)
	authParams["oauth_consumer_key"] = CONSUMERKEY
	authParams["oauth_token"] = token
	authParams["oauth_nonce"] = makeNonce(32)
	authParams["oauth_signature_method"] = "HMAC-SHA1"
	authParams["oauth_timestamp"] = strconv.Itoa(int(time.Now().Unix()))
	authParams["oauth_version"] = "1.0"
	header := makeAuthHeader(makeKey(CONSUMERSECRET,TOKENSECRET), "POST", queryUrl, authParams)
	data := qParams.Get("oauth_verifier")
	code, body := makeTokenRequest(method, queryUrl, header, data)
	fmt.Println("final result: ")
	fmt.Println(code, body)
	return 
}

func testEm() {
	fmt.Println("test functions: ")
	
	fmt.Println(makeNonce(32))
	fmt.Println(makeKey("abcd1234","efgh5678"))
	fmt.Println(makeKey("abcd1234",""))
	fmt.Println(makeSignature(makeKey("abcd1234","efgh5678"), "signthisstring"))
	
	var method string = "POST"
	var queryUrl string = "https://api.twitter.com/oauth/request_token"
	var authParams = make(map[string]string)
	var key = makeKey(CONSUMERSECRET,TOKENSECRET)
	
	authParams["oauth_callback"] = "http://localhost:8080/callback"
	authParams["oauth_token"] = TOKEN
	authParams["oauth_consumer_key"] = CONSUMERKEY
	authParams["oauth_nonce"] = makeNonce(32)
	authParams["oauth_signature_method"] = "HMAC-SHA1"
	authParams["oauth_timestamp"] = strconv.Itoa(int(time.Now().Unix()))
	authParams["oauth_version"] = "1.0"
	header := makeAuthHeader(key, method, queryUrl, authParams)
	fmt.Println(header)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth", auth)
	r.HandleFunc("/callback", callback)
	http.ListenAndServe(":8080", r)	
}
