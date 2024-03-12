package helpers

// Sleep is a helper function that enables free accounts to stay within their limit of 1 request per second
func Sleep() {
	// Uncomment this to enable testing on low-TPS API keys
	//time.Sleep(1 * time.Second)
}
