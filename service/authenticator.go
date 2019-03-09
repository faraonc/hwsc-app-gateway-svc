package service

import (
	"context"
	"encoding/base64"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	pbauth "github.com/hwsc-org/hwsc-lib/auth"
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
	//return tryBasicAuth(ctx)

	//TODO remove when we don't want to support insecure dialing
	newCtx, err = tryBasicAuth(ctx)
	if err == nil {
		return newCtx, nil
	}

	return context.TODO(), nil
}

func tryBasicAuth(ctx context.Context) (context.Context, error) {
	log.Info(consts.BasicAuthTag, consts.StrAuthAttempt)
	// auth := b.email + ":" + b.password
	// enc := base64.StdEncoding.EncodeToString([]byte(auth))
	// the format is "authorization": "Basic " + enc,
	auth, err := extractContextHeader(ctx, consts.StrMdBasicAuthHeader)
	if err != nil {
		log.Error(consts.BasicAuthTag, err.Error())
		return ctx, err
	}
	log.Info(consts.BasicAuthTag, auth)

	if !strings.HasPrefix(auth, consts.StrBasicAuthPrefix) {
		log.Error(consts.BasicAuthTag, consts.ErrMissingBasicAuthPrefix.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrMissingBasicAuthPrefix.Error())
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(consts.StrBasicAuthPrefix):])
	if err != nil {
		log.Error(consts.BasicAuthTag, consts.ErrInvalidBase64Header.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrInvalidBase64Header.Error())
	}

	emailPassword := string(c)
	s := strings.IndexByte(emailPassword, ':')
	// email:password
	if s < 0 || emailPassword[:s] == "" || emailPassword[s+1:] == "" {
		log.Error(consts.BasicAuthTag, consts.ErrInvalidBasicAuthFormat.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrInvalidBasicAuthFormat.Error())
	}

	// validate with user service here
	resp, err := userSvc.authenticateUser(emailPassword[:s], emailPassword[s+1:])
	if err != nil {
		log.Error(consts.BasicAuthTag, err.Error())
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}

	if err := pbauth.ValidateIdentification(resp.Identification); err != nil {
		log.Error(consts.BasicAuthTag, err.Error())
		return ctx, status.Error(codes.Internal, err.Error())
	}

	// Remove token from headers from here on
	cleanCtx, err := purgeContextHeader(ctx, consts.StrMdBasicAuthHeader)
	if err != nil {
		log.Error(consts.BasicAuthTag, err.Error())
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}
	// Add auth token string
	retMd := metadata.Pairs(consts.StrMdAuthToken, resp.Identification.Token)
	return metadata.NewIncomingContext(cleanCtx, retMd), nil
}

func extractContextHeader(ctx context.Context, header string) (string, error) {
	if strings.TrimSpace(header) == "" {
		return "", status.Error(codes.Unauthenticated, consts.ErrMissingHeader.Error())
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, consts.ErrMissingAuthHeadersFromCtx.Error())
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

func purgeContextHeader(ctx context.Context, header string) (context.Context, error) {
	if strings.TrimSpace(header) == "" {
		return nil, consts.ErrMissingHeader
	}
	md, _ := metadata.FromIncomingContext(ctx)
	mdCopy := md.Copy()
	mdCopy[header] = nil
	return metadata.NewIncomingContext(ctx, mdCopy), nil
}
