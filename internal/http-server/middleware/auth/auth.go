// Поскольку URL Shortener — это REST API сервис,
// авторизацию удобнее всего реализовать в виде middleware (аналог интерсепторов gRPC).
// В случае HTTP-запросов, JWT-токен обычно отправляют в заголовке вида:
// Authorization: "Bearer <jwt_token>"
// Поэтому, для получения токена напишем:

package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"url-shortener/internal/lib/jwt"
	"url-shortener/internal/lib/logger/sl"
)

//type PermissionProvider struct{
//
//}

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check is user is admin")
)

// extractBeclaims, err:=arerToken извлекает JWT-токен из заголовка http-запроса
func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

// сам middleware:
// added by Alexx:

type keyUID int64
type myBool bool
type myError error

var (
	uidKey     keyUID
	errorKey   myError
	isAdminKey myBool
)

// New creates new auth middleware.
func New(
	log *slog.Logger,
	appSecret string,
	//permProvider PermissionProvider,
) func(next http.Handler) http.Handler {
	const op = "middleware.auth.New"

	log = log.With(slog.String("op", op))

	// Возвращаем функцию-обработчик
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем JWT-токен из запроса
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				// It's ok, if user is not authorized
				next.ServeHTTP(w, r)
				return
			}

			// Парсим и валидируем токен, используя appSecret
			claims, err := jwt.Parse(tokenStr, appSecret)
			if err != nil {
				log.Warn("failed to parse token", sl.Err(err))

				// But if token is invalid, we shouldn't handle request
				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			// Отправляем запрос для проверки, является ли пользователь админом
			// пока закомментировал
			isAdmin := false
			/*
				isAdmin, err := permProvider.IsAdmin(r.Context(), claims.UID)
				if err != nil {
					log.Error("failed to check if user is admin", sl.Err(err))

					ctx := context.WithValue(r.Context(), errorKey, ErrFailedIsAdminCheck)
					next.ServeHTTP(w, r.WithContext(ctx))

					return
				}
			*/
			// Полученные данные сохраняем в контекст,
			// откуда его смогут получить следующие хэндлеры.

			ctx := context.WithValue(r.Context(), uidKey, claims.UID)
			ctx = context.WithValue(r.Context(), isAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(uidKey).(int64)
	return uid, ok
}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(errorKey).(error)
	return err, ok
}

/*
func New(
	log *slog.Logger,
	appSecret string,
	permProvider PermissionProvider,
) func(next http.Handler) http.Handler {
	const op = "middleware.auth.New"

	log = log.With(slog.String("op", op))

	//возвращаем функцию-обработчик
	return func(next http.Handler) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем JWT-токен из запроса
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				// it's ok, if user is not authorized
				next.ServeHTTP(w, r)
				return
			}

			// парсим и валидируем токен, используя appSecret
			claims, err := jwt.Parse(tokenStr, appSecret)
			if err != nil {
				log.Warn("failed to parse token", sl.Err(err))

				// But if token is invalid, we shouldn't handle request
				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			// отправялем запрос для проверки, является ли пользователь админом
			// TODO: - проверить
			isAdmin, err := permProvider.IsAdmin(r.Context(), claims.UserID)
			if err != nil {
				log.Error("failed to check if user is admin", sl.Err(err))

				ctx := context.WithValue(r.Context(), errorKey, ErrFailedIsAdminCheck)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			// полученные данные сохраняем в контекст,
			// откуда его смогут получить следующие хендлеры.
			ctx := context.WithValue(r.Context(), uidKey, claims.UserID)
			ctx = context.WithValue(r.Context(), isAdminKey, isAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
*/
