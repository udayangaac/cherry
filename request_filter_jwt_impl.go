// Copyright 2021 by Chamith Udayanga. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cherry

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"time"
)

type claims struct {
	ID    int         `json:"id"`
	Roles []string    `json:"roles"`
	Data  interface{} `json:"data"`
	Iat   int         `json:"iat"`
	Exp   int         `json:"exp"`
}

type jwtRequestFilter struct {
	Secret        string
	ValidDuration time.Duration
}

type JwtConfig struct {
	Secret        string
	ValidDuration time.Duration
}

func NewJwtRequestFilter(config JwtConfig) RequestFilter {
	return &jwtRequestFilter{
		Secret:        config.Secret,
		ValidDuration: config.ValidDuration,
	}
}

func (j jwtRequestFilter) DoFilter(ctx context.Context, token string, roles []string) (authentication Authentication, status EntryPointStatus) {
	c := claims{}
	if token != "" {
		t, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error in parsing token")
			}
			return []byte(j.Secret), nil
		})

		if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			if err := mapstructure.Decode(claims, &c); err != nil {
				status = EntryPointStatus{
					Code: InvalidToken,
					Desc: "invalid token",
				}
				return
			}

			if !isValidRole(roles, c.Roles) {
				status = EntryPointStatus{
					Code: Unauthorized,
					Desc: "unauthorized",
				}
			}

			if isExpired(c.Iat, c.Exp, j.ValidDuration) {
				status = EntryPointStatus{
					Code: TokenExpired,
					Desc: "expired",
				}
			}

			authentication = Authentication{
				ID:    c.ID,
				Roles: c.Roles,
				Data:  c.Data,
			}

			return
		} else {
			status = EntryPointStatus{
				Code: InvalidToken,
				Desc: "invalid token",
			}
			return
		}
	} else {
		status = EntryPointStatus{
			Code: EmptyToken,
			Desc: "empty token",
		}
		return
	}
}

func isValidRole(roles []string, rolesOfToken []string) (isValidRole bool) {
	rolesMap := make(map[string]bool)
	for _, v := range roles {
		rolesMap[v] = true
	}
	for _, v := range rolesOfToken {
		if rolesMap[v] {
			isValidRole = true
		}
	}
	return
}

func isExpired(iat, exp int, validDuration time.Duration) (isExpired bool) {
	isExpired = exp-iat < int(validDuration.Seconds())
	return
}

// GenerateJwtToke generates Jwt token with the values of fields inside the
// parameter auth  Valid duration and Secret key must be parse with JwtConfig.
// NewJwtRequestFilter(./request_filter_jwt_impl.go > NewJwtRequestFilter ) config
// structure must be same.
func GenerateJwtToke(auth Authentication, config JwtConfig) (token string, err error) {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    auth.ID,
		"roles": auth.Roles,
		"data":  auth.Data,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Unix() + int64(config.ValidDuration.Seconds()),
	})
	token, err = tokenObj.SignedString([]byte(config.Secret))
	return
}
