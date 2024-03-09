package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type myClaims struct {
	UID int64
}

func Parse(tokenStr string, appSecret string) (*myClaims, error) {
	const op = "internal.lib.jwt.Parse"
	// Парсим и валидируем токен
	tokenParsed, err := jwt.Parse(tokenStr, func(tokenStr *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	// Если ключ окажется невалидным, мы получим соответствующую ошибку
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Преобразуем к типу jwt.MapClaims, в котором мы сохраняли данные
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	clms := myClaims{
		UID: int64(claims["uid"].(float64)),
	}

	return &clms, nil
}

/*
func Parse(tokenStr string, appSecret string) (*jwt.MapClaims, error) {
	const op = "internal.lib.jwt.Parse"
	// Парсим и валидируем токен
	tokenParsed, err := jwt.Parse(tokenStr, func(tokenStr *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	// Если ключ окажется невалидным, мы получим соответствующую ошибку
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Преобразуем к типу jwt.MapClaims, в котором мы сохраняли данные
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &claims, nil
}
*/
/*
func Parse(tokenStr string, appSecret string) (*jwt.Token, error) {
	const op = "internal.lib.jwt.Parse"
	// Парсим и валидируем токен
	tokenParsed, err := jwt.Parse(tokenStr, func(tokenStr *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	// Если ключ окажется невалидным, мы получим соответствующую ошибку
	if err != nil {

		return nil, err
	}
	return tokenParsed, nil
}
*/
