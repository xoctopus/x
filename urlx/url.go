// Package urlx provides a more user-friendly URL parser, modifier and builder.
// TODO Opaque support
package urlx

import (
	"net/url"
	"strconv"
	"strings"
)

func Parse(raw string, modifiers ...Modifier) (*URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return &URL{}, err
	}
	return From(*u, modifiers...), nil
}

func From(from url.URL, modifiers ...Modifier) *URL {
	from.ForceQuery = false
	from.OmitHost = false

	u := &URL{from}
	u.Scheme = strings.ToLower(u.Scheme)

	for _, modifier := range modifiers {
		modifier(&u.URL)
	}
	return u
}

func Build(modifiers ...Modifier) *URL {
	return From(url.URL{}, modifiers...)
}

type Modifier func(u *url.URL)

func WithScheme(scheme string) Modifier {
	return func(u *url.URL) {
		u.Scheme = strings.ToLower(scheme)
	}
}

// TODO
// func WithOpaque(opaque string) Modifier {
// 	return func(u *url.URL) {
// 		u.Opaque = opaque
// 	}
// }

func WithUserinfo(username, password string) Modifier {
	return func(u *url.URL) {
		u.User = url.UserPassword(username, password)
	}
}

func WithHost(host string) Modifier {
	return func(u *url.URL) {
		u.Host = host
	}
}

func WithPort(port uint16) Modifier {
	return func(u *url.URL) {
		u.Host = strings.TrimSuffix(u.Host, ":"+u.Port())
		u.Host += ":" + strconv.Itoa(int(port))
	}
}

func TrimPort() Modifier {
	return func(u *url.URL) {
		u.Host = strings.TrimSuffix(u.Host, ":"+u.Port())
	}
}

func WithDefaultPort() Modifier {
	return func(u *url.URL) {
		if port, ok := defaultPorts[u.Scheme]; ok {
			u.Host = strings.TrimSuffix(u.Host, ":"+u.Port())
			u.Host += ":" + strconv.Itoa(int(port))
		}
	}
}

func WithPath(path string) Modifier {
	return func(u *url.URL) {
		if len(path) > 0 {
			path = "/" + strings.TrimPrefix(path, "/")
			u.Path = path
		}
	}
}

func WithQuery(query string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		if q1, err := url.ParseQuery(query); err == nil {
			for key, values := range q1 {
				for _, v := range values {
					q.Add(key, v)
				}
			}
		}
		u.RawQuery = q.Encode()
	}
}

func WithQueryParams(q url.Values) Modifier {
	return func(u *url.URL) {
		q0 := u.Query()
		for key, values := range q {
			for _, v := range values {
				q0.Add(key, v)
			}
		}
		u.RawQuery = q0.Encode()
	}
}

func TrimQuery() Modifier {
	return func(u *url.URL) {
		u.RawQuery = ""
	}
}

func WithFragment(fragment string) Modifier {
	return func(u *url.URL) {
		fragment = strings.TrimPrefix(fragment, "#")
		u.Fragment = fragment
	}
}

func TrimFragment() Modifier {
	return func(u *url.URL) {
		u.Fragment = ""
	}
}

var defaultPorts = map[string]uint16{
	// web + base
	"http":  80,
	"https": 443,
	"ssh":   22,
	"dns":   53,

	// rdb
	"mysql":    3306,
	"postgres": 5432,
	"pg":       5432,
	"oracle":   1521,
	"mssql":    1433,
	"tidb":     4000,

	// no-sql + cache
	"redis":     6379,
	"memcached": 11211,
	"mongodb":   27017,

	// mq
	"amqp":     5672,
	"kafka":    9092,
	"pulsar":   6650,
	"mqtt":     1883,
	"rocketmq": 9876,
}

func DefaultPort(scheme string) (port uint16, ok bool) {
	port, ok = defaultPorts[strings.ToLower(scheme)]
	return
}

func AddDefaultPort(scheme string, port uint16) (conflicted bool) {
	scheme = strings.ToLower(scheme)
	_, conflicted = defaultPorts[scheme]
	defaultPorts[scheme] = port
	return
}

type URL struct {
	url.URL
}

func (u *URL) Port() uint16 {
	if port := u.URL.Port(); len(port) > 0 {
		if i, err := strconv.ParseUint(port, 10, 16); err == nil {
			return uint16(i)
		}
	}
	if port, ok := DefaultPort(u.Scheme); ok {
		return port
	}
	return 0
}

func (u *URL) IsZero() bool {
	return len(u.Host) == 0
}
