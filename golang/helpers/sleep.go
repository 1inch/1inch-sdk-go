package helpers

import "time"

// Sleep is a helper function that enables free accounts to stay within their limit of 1 request per second
func Sleep() {
	time.Sleep(1 * time.Second)
}
