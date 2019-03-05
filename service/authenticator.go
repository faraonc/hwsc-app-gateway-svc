package service

import (
	"context"
	"encoding/base64"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// Authenticate attempts to use TLS, Token, and Basic for authentication.
// Return the context used for the authentication.
func Authenticate(ctx context.Context) (newCtx context.Context, err error) {
	// TODO
	//err = a.tryTLSAuth(ctx)
	//if err == nil {
	//	return ctx, nil
	//}

	//newCtx, err = a.tryTokenAuth(ctx)
	//if err == nil {
	//	return newCtx, nil
	//}
	newCtx, err = tryBasicAuth(ctx)
	if err == nil {
		return newCtx, nil
	}

	// TODO remove when we don't want to support insecure dialing
	return context.TODO(), nil
}

func tryBasicAuth(ctx context.Context) (context.Context, error) {
	log.Info(consts.BasicAuthTag)
	// auth := b.email + ":" + b.password
	// enc := base64.StdEncoding.EncodeToString([]byte(auth))
	// the format is "authorization": "Basic " + enc,
	auth, err := extractHeader(ctx, consts.BasicAuthHeader)
	if err != nil {
		log.Error(consts.BasicAuthTag, err.Error())
		return ctx, err
	}
	log.Info(consts.BasicAuthTag, auth)

	if !strings.HasPrefix(auth, consts.BasicAuthPrefix) {
		log.Error(consts.BasicAuthTag, consts.ErrMissingBasicAuthPrefix.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrMissingBasicAuthPrefix.Error())
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(consts.BasicAuthPrefix):])
	if err != nil {
		log.Error(consts.BasicAuthTag, consts.ErrInvalidBase64Header.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrInvalidBase64Header.Error())
	}

	emailPassword := string(c)
	s := strings.IndexByte(emailPassword, ':')
	if s < 0 {
		log.Error(consts.BasicAuthTag, consts.ErrInvalidBasicAuthFormat.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrInvalidBasicAuthFormat.Error())
	}

	// validate with user service here
	// TODO figure out what to do with resp containing a User object
	// TODO maybe add as metadata in context
	_, err = userSvc.authenticateUser(emailPassword[:s], emailPassword[s+1:])
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}

	// Remove token from headers from here on
	return purgeHeader(ctx, consts.BasicAuthHeader), nil
}

func extractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, consts.ErrMissingAuthHeaders.Error())
	}
	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, consts.ErrMissingAuthHeader.Error())
	}
	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, consts.ErrMultipleAuthHeaders.Error())
	}
	return authHeaders[0], nil
}

func purgeHeader(ctx context.Context, header string) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	mdCopy := md.Copy()
	mdCopy[header] = nil
	return metadata.NewIncomingContext(ctx, mdCopy)
}
