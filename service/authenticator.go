package service

import (
	"context"
	"encoding/base64"
	"github.com/hwsc-org/hwsc-api-blocks/int/lib"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	pbauth "github.com/hwsc-org/hwsc-lib/auth"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// Authenticate attempts to use TLS, Token, and Basic for authentication.
// Email Token is also authenticated to complete email verification.
// Return the context used for the authentication.
func Authenticate(ctx context.Context) (newCtx context.Context, err error) {
	newCtx, err = tryEmailTokenVerification(ctx)
	if err == nil {
		return newCtx, nil
	}
	// TODO
	//err = a.tryTLSAuth(ctx)
	//if err == nil {
	//	return ctx, nil
	//}

	newCtx, err = tryTokenAuth(ctx)
	if err == nil {
		return newCtx, nil
	}

	newCtx, err = tryBasicAuth(ctx)
	if err == nil {
		return newCtx, nil
	}

	//TODO remove when we don't want to support insecure dialing
	return context.TODO(), nil
}

// tryEmailTokenVerification checks if browser intends to verify an email.
// The header has a format of "authorization": "Email Token " + token
// Returns context with no token or an error.
func tryEmailTokenVerification(ctx context.Context) (context.Context, error) {
	log.Info(consts.EmailVerificationTag, consts.StrTokenAuthAttempt)
	return nil, nil
}

// tryTokenAuth checks if browser intends to authenticate using an existing auth token.
// The header has a format of "authorization": "Auth Token " + token.
// Returns context with token or an error.
func tryTokenAuth(ctx context.Context) (context.Context, error) {
	log.Info(consts.TokenAuthTag, consts.StrTokenAuthAttempt)
	auth, err := extractContextHeader(ctx, consts.StrMdBasicAuthHeader)
	if err != nil {
		return ctx, err
	}
	log.Info(consts.TokenAuthTag, auth)
	if !strings.HasPrefix(auth, consts.StrTokenAuthPrefix) {
		return ctx, status.Error(codes.Unauthenticated, consts.ErrMissingTokenPrefix.Error())
	}

	resp, err := userSvc.verifyAuthToken(auth[len(consts.StrBasicAuthPrefix):])
	st, ok := status.FromError(err)
	if !ok {
		log.Error(consts.TokenAuthTag, st.Message())
		return ctx, status.Error(codes.Unauthenticated, st.Message())
	}

	return finalizeAuth(consts.TokenAuthTag, resp.Identification, ctx)
}

// tryBasicAuth checks if browser intends to authenticate using a base64 encoded "email:password"
// The header has a format of "authorization": "Basic " + base64 encoded "email:password".
// Returns context with token or an error.
func tryBasicAuth(ctx context.Context) (context.Context, error) {
	log.Info(consts.BasicAuthTag, consts.StrBasicAuthAttempt)
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
	st, ok := status.FromError(err)
	if !ok {
		log.Error(consts.BasicAuthTag, st.Message())
		return ctx, status.Error(codes.Unauthenticated, st.Message())
	}

	return finalizeAuth(consts.BasicAuthTag, resp.Identification, ctx)
}

// finalizeAuth validates the Identification, and sanitizes the context with the token.
// Returns the context with token or an error.
func finalizeAuth(tag string, id *lib.Identification, ctx context.Context) (context.Context, error) {
	if strings.TrimSpace(tag) == "" {
		return ctx, status.Error(codes.Unauthenticated, consts.ErrMissingTag.Error())
	}
	if err := pbauth.ValidateIdentification(id); err != nil {
		log.Error(tag, err.Error())
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}

	// Sanitize the header used for dialing
	cleanCtx, err := purgeContextHeader(ctx, consts.StrMdBasicAuthHeader)
	if err != nil {
		log.Error(tag, err.Error())
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}

	// Add auth token string
	retMd := metadata.Pairs(consts.StrMdAuthToken, id.Token)
	return metadata.NewIncomingContext(cleanCtx, retMd), nil
}

// extractContextHeader extracts the header from the context.
// Returns the header from the context.
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

// purgeContextHeader removes a specific header from the context.
// Returns the sanitized context.
func purgeContextHeader(ctx context.Context, header string) (context.Context, error) {
	if strings.TrimSpace(header) == "" {
		return nil, consts.ErrMissingHeader
	}
	md, _ := metadata.FromIncomingContext(ctx)
	mdCopy := md.Copy()
	mdCopy[header] = nil
	return metadata.NewIncomingContext(ctx, mdCopy), nil
}
