package extensions

import "net/http"

func IsSuccessStatusCode(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode <= 299
}