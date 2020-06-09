package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/joho/godotenv"
	"github.com/tendermint/tmlibs/bech32"
	"github.com/tomasen/realip"
)

var chain string
var recaptchaSecretKey string
var amountFaucet string
var amountSteak string
var key string
var pass string
var node string
var publicURL string
var faucetHome string
var fees string

type claimStruct struct {
	Address  string
	Response string
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(key, "=", value)
		return value
	} else {
		log.Fatal("Error loading environment variable: ", key)
		return ""
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rand.Seed(time.Now().UnixNano())

	chain = getEnv("FAUCET_CHAIN")
	recaptchaSecretKey = getEnv("FAUCET_RECAPTCHA_SECRET_KEY")
	amountFaucet = getEnv("FAUCET_AMOUNT_FAUCET")
	amountSteak = getEnv("FAUCET_AMOUNT_STEAK")
	key = getEnv("FAUCET_KEY")
	pass = getEnv("FAUCET_PASS")
	node = getEnv("FAUCET_NODE")
	publicURL = getEnv("FAUCET_PUBLIC_URL")
	faucetHome = getEnv("FAUCET_HOME")
	fees = getEnv("FAUCET_FEES")

	r := mux.NewRouter()
	recaptcha.Init(recaptchaSecretKey)

	r.HandleFunc("/claim", getCoinsHandler)
	r.HandleFunc("/claim/wallet", getWalletCoinsHandler)

	log.Fatal(http.ListenAndServe(publicURL, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Token"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), handlers.AllowedOrigins([]string{"*"}))(r)))

}

func executeCmd(command string, writes ...string) {
	cmd := getCmd(command)

	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		for _, write := range writes {
			stdin.Write([]byte(write + "\n"))
		}
	}()
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
}

func getCmd(command string) *exec.Cmd {
	// split command into command and args
	split := strings.Split(command, " ")

	var cmd *exec.Cmd
	if len(split) == 1 {
		cmd = exec.Command(split[0])
	} else {
		cmd = exec.Command(split[0], split[1:]...)
	}

	return cmd
}

func getCoinsHandler(w http.ResponseWriter, request *http.Request) {
	var claim claimStruct

	// decode JSON response from front end
	decoder := json.NewDecoder(request.Body)
	decoderErr := decoder.Decode(&claim)
	if decoderErr != nil {
		panic(decoderErr)
	}

	// make sure address is bech32
	readableAddress, decodedAddress, decodeErr := bech32.DecodeAndConvert(claim.Address)
	if decodeErr != nil {
		panic(decodeErr)
	}
	// re-encode the address in bech32
	encodedAddress, encodeErr := bech32.ConvertAndEncode(readableAddress, decodedAddress)
	if encodeErr != nil {
		panic(encodeErr)
	}

	// make sure captcha is valid
	clientIP := realip.FromRequest(request)
	captchaResponse := claim.Response
	captchaPassed, captchaErr := recaptcha.Confirm(clientIP, captchaResponse)
	if captchaErr != nil {
		panic(captchaErr)
	}

	nonce := rand.Uint64()

	// send the coins!
	if captchaPassed {

		fmt.Println(encodedAddress)
		txCount += 1
		sendFaucet := fmt.Sprintf("colorcli tx send %s %s --memo %d --from %s --chain-id %s --fees %s --home %s --node %s --nonce %d",
			encodedAddress, amountFaucet, txCount, key, chain, fees, faucetHome, node, nonce)
		fmt.Println("Command: ", sendFaucet)
		fmt.Println(time.Now().UTC().Format(time.RFC3339), encodedAddress, "[1]")
		go executeCmd(sendFaucet, "y", pass)
	}
	return
}

var (
	txCount = 0
)

func getWalletCoinsHandler(w http.ResponseWriter, request *http.Request) {
	var claim claimStruct

	// decode JSON response from front end
	decoder := json.NewDecoder(request.Body)
	decoderErr := decoder.Decode(&claim)
	if decoderErr != nil {
		panic(decoderErr)
	}

	// make sure address is bech32
	readableAddress, decodedAddress, decodeErr := bech32.DecodeAndConvert(claim.Address)
	if decodeErr != nil {
		panic(decodeErr)
	}
	// re-encode the address in bech32
	encodedAddress, encodeErr := bech32.ConvertAndEncode(readableAddress, decodedAddress)
	if encodeErr != nil {
		panic(encodeErr)
	}

	nonce := rand.Uint64()

	txCount += 1
	sendFaucet := fmt.Sprintf("colorcli tx send %s %s --memo %d --from %s --chain-id %s --fees %s --home %s --node %s --nonce %d",
		encodedAddress, amountFaucet, txCount, key, chain, fees, faucetHome, node, nonce)
	fmt.Println("Command: ", sendFaucet)
	fmt.Println(time.Now().UTC().Format(time.RFC3339), encodedAddress, "[1]")
	go executeCmd(sendFaucet, "y", pass)

	return
}
