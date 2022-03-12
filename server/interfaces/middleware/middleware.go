package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/csrf"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

// func AuthMid(next http.HandlerFunc, cookieApp application.AuthAppInterface) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, found := CheckCookies(r, cookieApp)
// 		if !found {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), entity.CookieInfoKey, cookie)
// 		r = r.Clone(ctx)

// 		next.ServeHTTP(w, r)
// 	})
// }

// func NoAuthMid(next http.HandlerFunc, cookieApp application.AuthAppInterface) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		_, found := CheckCookies(r, cookieApp)
// 		if found {
// 			w.WriteHeader(http.StatusForbidden)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// PanicMid logges error if handler errors
func PanicMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				logger.Info(err.(error).Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CSRFSettingMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r != nil {
			if r.Header.Get("X-CSRF-Token") == "" {
				token := csrf.Token(r)
				w.Header().Set("X-CSRF-Token", token)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// // CheckCookies returns *CookieInfo and true if cookie is present in sessions slice, nil and false othervise
// func CheckCookies(r *http.Request, cookieApp application.AuthAppInterface) (*entity.CookieInfo, bool) {
// 	cookie, err := r.Cookie(string(entity.CookieNameKey))
// 	if err == http.ErrNoCookie {
// 		return nil, false
// 	}

// 	return cookieApp.CheckCookie(cookie)
// }
