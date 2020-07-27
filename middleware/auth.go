package middleware

import (
	"my-filestore-server/redis"
	"net/http"
)

//HTTPInterceptor http请求拦截器
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		username := r.Form.Get("username")
		token := r.Form.Get("token")

		if len(username) < 3 || !isTokenValid(username, token) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

//IsTokenValid 验证token
func isTokenValid(username string, token string) bool {
	t, _ := redis.Get("token_" + username)
	return t == token
}
