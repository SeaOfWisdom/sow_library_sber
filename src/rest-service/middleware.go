package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	srv "github.com/SeaOfWisdom/sow_library/src/service"
	"github.com/SeaOfWisdom/sow_library/src/service/storage"
	jwt "github.com/SeaOfWisdom/sow_proto/jwt-srv"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrNoToken      = errors.New("token is null")
)

func (rs *RestSrv) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var realRole storage.ParticipantRole
		path := strings.Split(r.URL.Path, "/")
		requestUri := r.URL.Path
		if len(path) > 1 {
			requestUri = path[1]
		} else {
			rs.logger.Error(fmt.Sprintf("WRONG URL PATH: %v", path))
		}
		for uri, requiredRole := range authURIs {
			if strings.EqualFold(requestUri, uri) {
				token, err := rs.getTokenFromHeader(r)
				if err != nil {
					rs.logger.Warn(fmt.Sprintf("while getting a jwt token from header, err: %v", err))
					responError(w, http.StatusBadGateway, err.Error())
					return
				}
				realRole, err = rs.VerifyJWT(token)
				if err != nil {
					responError(w, http.StatusBadGateway, err.Error())
					return
				}

				if realRole < requiredRole {
					responError(w, http.StatusNetworkAuthenticationRequired, "you haven't been granted access to this method")
					return
				}
				break
			}
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		ctx := context.WithValue(r.Context(), "role", realRole)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *RestSrv) VerifyJWT(token string) (storage.ParticipantRole, error) {
	if token == "" {
		return -1, fmt.Errorf("token is null")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := rs.jwtSrv.DecodeJWT(ctx, &jwt.Token{Token: token})
	if err != nil {
		return -1, err
	}

	if !resp.Valid || resp.GetBody().Iss != srv.ServiceName {
		return -1, fmt.Errorf("token is not valid")
	}

	return storage.ParticipantRole(resp.Body.Role), nil
}

func (rs *RestSrv) getTokenFromHeader(r *http.Request) (token string, err error) {
	// extract jwt token from the certain header
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	if tokenHeader == "" {                       //Token is missing, returns with error code 403 Unauthorized
		err = ErrNoToken

		return
	}

	// tokenHeader = "bearer+jwt_token"
	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		err = ErrInvalidToken

		return
	}

	return splitted[1], nil //Grab the token part, what we are truly interested in
}

func (rs *RestSrv) getWeb3Address(r *http.Request) (token string, err error) {
	token, err = rs.getTokenFromHeader(r)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	resp, dErr := rs.jwtSrv.DecodeJWT(ctx, &jwt.Token{Token: token})
	if dErr != nil {
		err = fmt.Errorf("%v, err: %v", ErrInvalidToken, dErr)

		return
	}

	if !common.IsHexAddress(strings.ToLower(resp.Body.Web3Address)) {
		err = fmt.Errorf("the wrong web3 address %s, decode err: %v", resp.Body.Web3Address, err)

		return
	}

	return resp.Body.Web3Address, nil
}
