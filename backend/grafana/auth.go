package grafana

import "net/http"

type TokenSource interface {
	Get() (string, error)
}

func BearerToken(token TokenSource) Auth {
	return func(req *http.Request) error {
		v, err := token.Get()
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+v)
		return nil
	}
}

func BasicAuth(username string, password TokenSource) Auth {
	return func(req *http.Request) error {
		v, err := password.Get()
		if err != nil {
			return err
		}
		req.SetBasicAuth(username, v)
		return nil
	}
}

func APIKeyHeader(header string, key TokenSource) Auth {
	return func(req *http.Request) error {
		v, err := key.Get()
		if err != nil {
			return err
		}
		req.Header.Set(header, v)
		return nil
	}
}
