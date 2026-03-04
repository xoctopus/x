package urlx_test

import (
	"fmt"
	"net/url"

	"github.com/xoctopus/x/urlx"
)

func ExampleBuild() {
	u := urlx.Build(
		urlx.WithScheme("HTTPs"),
		urlx.WithUserinfo("user", "pass"),
		urlx.WithHost("example.com:9999"),
		urlx.WithPort(80),
		urlx.TrimPort(),
		urlx.WithDefaultPort(),
		urlx.WithPath("login"),
		urlx.WithQuery("k1=v1&k1=v2&k1=v3"),
		urlx.WithQueryParams(url.Values{
			"k2": []string{"v1", "v2", "v3"},
		}),
	)
	fmt.Println(u.String())

	u, _ = urlx.Parse(u.String(), urlx.TrimQuery(), urlx.WithFragment("section1"))
	fmt.Println(u.String())

	u = urlx.From(u.URL, urlx.TrimFragment(), urlx.TrimPort())
	fmt.Println(u.String())
	fmt.Println(u.Port())

	u = urlx.From(u.URL, urlx.WithScheme("test"), urlx.WithPort(30000))
	fmt.Println(u.Port())

	u = urlx.From(u.URL, urlx.WithScheme("test"), urlx.TrimPort())
	fmt.Println(u.Port())

	urlx.AddDefaultPort("test", 30000)
	fmt.Println(u.Port())

	u, _ = urlx.Parse("http://localhost\n:8080")
	fmt.Println(u.IsZero())

	// Output:
	// https://user:pass@example.com:443/login?k1=v1&k1=v2&k1=v3&k2=v1&k2=v2&k2=v3
	// https://user:pass@example.com:443/login#section1
	// https://user:pass@example.com/login
	// 443
	// 30000
	// 0
	// 30000
	// true
}
