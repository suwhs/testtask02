package server

func makeHeaders() map[string]string {
	return map[string]string{
		"authority":          "www.rusprofile.ru",
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "ru-RU,ru;q=0.9,en-GB;q=0.8,en;q=0.7,en-US;q=0.6",
		"cache-control":      "no-cache",
		"pragma":             "no-cache",
		"referer":            "https://www.rusprofile.ru/",
		"sec-ch-ua":          `"Chromium";v="110", "Not A(Brand";v="24", "Google Chrome";v="110"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Linux\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	}
}
