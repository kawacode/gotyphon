package gotyphon

import (
	http "github.com/useflyent/fhttp"
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
