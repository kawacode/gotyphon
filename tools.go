package gotyphon

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

// Delete a slice of uint16s from a slice of uint16s
func DeleteFromIntSlice(slice []uint16, start int, end int) {
	slices.Delete(slice, start, end)
}

// It takes a map of string slices and returns a pointer to an http.Header
func MapStringSlicesToHttpHeaders(headers map[string][]string) *http.Header {
	var result = http.Header{}
	for key, value := range headers {
		for _, v := range value {
			http.Header.Add(result, key, v)
		}
	}
	return &result
}

// > Converts a map of strings to a map of string slices
func MapStringToMapStringSlice(MapString map[string]string, bot *BotData) map[string][]string {
	var result = make(map[string][]string)
	for key, value := range MapString {
		result[key] = []string{value}
	}
	if len(bot.HttpRequest.Request.HeaderOrderKey) > 0 {
		result["Header-Order:"] = bot.HttpRequest.Request.HeaderOrderKey
	}
	if len(bot.HttpRequest.Request.PHeaderOrderKey) > 0 {
		result["PHeader-Order:"] = bot.HttpRequest.Request.PHeaderOrderKey
	}
	return result
}
// Delay pauses the program for a given number of seconds.
func Delay(seconds time.Duration) {
	time.Sleep(seconds * time.Second)
}

// It takes a string, and two other strings, and returns the string between the two other strings.
func Parse(str string, left string, right string) string {
	return strings.Split(strings.Split(str, left)[1], right)[0]
}

// It takes a command name and a list of arguments and returns the output of the command as a string
// and an error
func Executecommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

// Convert a map of string slices to a map of strings.
func MapStringSliceToMapString(headers map[string][]string) map[string]string {
	var result = make(map[string]string)
	for key, value := range headers {
		for _, value := range value {
			result[key] = value
		}
	}
	return result
}

// It takes a slice of uint16s, creates a new random number generator, and then iterates over the
// slice, swapping the current element with a random element in the slice
func ShuffleSlice(slice []uint16) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(slice); n > 0; n-- {
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
	}
}

// RandomizeHeaders takes a map of headers and a slice of headers to skip, and returns a map of
// randomized headers
func RandomizeHeaders(headers map[string]string, skip []string) (map[string]string, error) {
	result := map[string]string{}
	for key, value := range headers {
		for _, valueskip := range skip {
			if key != valueskip {
				if RandomBool().Bool() {
					result[key] = value
				}
			} else {
				result[key] = value
			}
		}
	}
	return result, nil
}

// It takes two integers, min and max, and returns a random integer between min and max
func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return (min + rand.Intn(max-min))
}

// > It takes a slice of uint16 and returns a slice of uint16 with duplicates removed
func RemoveDuplicatesFromInt(src []uint16) []uint16 {
	keys := make(map[uint16]bool)
	list := []uint16{}
	for _, entry := range src {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// ReadFile takes a file location as a string, opens the file, reads the file, and returns the file as
// a string and a slice of strings
func ReadFileByContain(filelocation string, contain string) (string, error) {
	var result []string
	file, err := os.Open(filelocation)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	defer file.Close()
	var resultstring string
	for _, v := range result {
		if strings.Contains(v, contain) {
			resultstring = v
		}
	}
	return resultstring, nil
}

// ReadFile takes a file location as a string, opens the file, reads the file, and returns the file as
// a string and a slice of strings
func ReadFile(filelocation string) ([]string, string, error) {
	var result []string
	file, err := os.Open(filelocation)
	if err != nil {
		return nil, "", err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	defer file.Close()
	var resultstring string
	for _, v := range result {
		resultstring += v + "\r\n"
	}
	return result, resultstring, nil
}

// ReadDir returns a list of filenames in a directory.
func ReadDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	filelist := []string{}
	for _, file := range files {
		if !file.IsDir() {
			filelist = append(filelist, file.Name())
		}
	}
	return filelist, nil
}

// It generates a random UUID
func RandomUUID() string {
	return uuid.New().String()
}

// It takes a list of strings and an integer, and returns a list of strings
func RandomSelectFromStringList(list []string, amount_of_items int) []string {
	items := []string{}
	for range make([]string, amount_of_items) {
		items = append(items, list[rand.Intn(len(list))])
	}
	return items
}

// A boolgen is a random number generator that generates random booleans.
//
// The first field, src, is a random number generator. The second field, cache, is a cache of random
// bits. The third field, remaining, is the number of random bits left in the cache.
//
// The next step is to define the methods of the type.
// @property src - The source of randomness.
// @property {int64} cache - The cache is the last random number generated by the source.
// @property {int} remaining - The number of bits left in the cache.
type boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

// A function that returns a random boolean value.
func (b *boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

// `RandomBool()` returns a pointer to a `boolgen` struct that has a `src` field that is a pointer to a
// `rand.Source` struct that has been initialized with a random seed
func RandomBool() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}

// It generates a random JA3 fingerprint
//Siteciper is the site default cipher suite
func RandomJA3() string {
	version := "771"
	ciphersuite := RandomCiphers()
	var ciphersuit []string
	for _, ciphersuite_values := range ciphersuite {
		ciphersuit = append(ciphersuit, fmt.Sprint(ciphersuite_values))
	}
	ciphersuites := strings.Join(ciphersuit, "-")
	extentioning := RandomExtension()
	var extentionn []string
	for _, extention_values := range extentioning {
		extentionn = append(extentionn, fmt.Sprint(extention_values))
	}
	extentions := strings.Join(extentionn, "-")
	ellipticcurv := RandomEllipticcurves()
	var ellipticcurvess []string
	for _, ellipticcurve_values := range ellipticcurv {
		ellipticcurvess = append(ellipticcurvess, fmt.Sprint(ellipticcurve_values))
	}
	ellipticcurves := strings.Join(ellipticcurvess, "-")
	ellipticcurvepointformat := "0"
	randomja3 := version + "," + ciphersuites + "," + extentions + "," + ellipticcurves + "," + ellipticcurvepointformat

	return randomja3
}
func RandomEllipticcurves() []uint16 {
	ellipticcurve := []uint16{29, 23, 24, 25, 30, 256, 257, 258, 259, 260}
	ShuffleSlice(ellipticcurve)
	ellipticcurv := RemoveDuplicatesFromInt(ellipticcurve)
	return ellipticcurv
}
func RandomExtension() []uint16 {
	extention := []uint16{0, 5, 10, 11, 13, 16, 18, 21, 22, 23, 27, 28, 34, 35, 43, 45, 49, 50, 51, 13172, 17513, 30032, 65281}
	ShuffleSlice(extention)
	extentions := RemoveDuplicatesFromInt(extention)
	return extentions
}
func RandomCiphers() []uint16 {
	ciphers := []uint16{4866, 4865, 4867, 49196, 49195, 52393, 49200, 52392, 49199, 159, 52394, 163, 158, 162, 49188, 49192, 49187, 49191, 107, 106, 103, 64, 49198, 49202, 49197, 49201, 49190, 49194, 49189, 49193, 49162, 49172, 49161, 49171, 57, 56, 51, 50, 49157, 49167, 49156, 49166, 157, 156, 61, 60, 53, 47, 49160, 49170, 22, 19, 49155, 49165, 10, 255}
	ShuffleSlice(ciphers)
	if RandomBool().Bool() {
		DeleteFromIntSlice(ciphers, RandomInt(0, 15), RandomInt(16, 30))
	}
	ciphers = append(ciphers, 4866)
	ciphersuite := RemoveDuplicatesFromInt(ciphers)
	return ciphersuite
}

/*
0 = Chrome; 1 = Firefox; 2 = Android; 3 = Computer; 4 = IOS; 5 = IPad 6 = IPhone; 7 = Android;
8 = InternetExplorer; 9 = Linux; 10 = MacOSX; 11 = Mobile; 12 = Safari; 13  = Random; Default = Random;
*/
func RandomUserAgent(browser_type int) string {
	switch browser_type {
	case 0:
		return browser.Chrome()
	case 1:
		return browser.Firefox()
	case 2:
		return browser.Android()
	case 3:
		return browser.Computer()
	case 4:
		return browser.IOS()
	case 5:
		return browser.IPad()
	case 6:
		return browser.IPhone()
	case 7:
		return browser.Android()
	case 8:
		return browser.InternetExplorer()
	case 9:
		return browser.Linux()
	case 10:
		return browser.MacOSX()
	case 11:
		return browser.Mobile()
	case 12:
		return browser.Safari()
	case 13:
		return browser.Random()
	default:
		return browser.Random()
	}
}
