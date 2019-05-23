package service

import (
	"context"
	"encoding/base64"
	"github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
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
// Return a context used for the authentication.
func Authenticate(ctx context.Context) (newCtx context.Context, err error) {
	log.RequestService(consts.AuthTag)
	if ctx == nil {
		log.Error(consts.AuthTag, consts.ErrNilContext.Error())
		return context.TODO(), consts.ErrNilContext
	}
	newCtx, errEmailToken := tryEmailTokenVerification(ctx)
	if errEmailToken == nil {
		return newCtx, nil
	}
	// TODO
	//newCtx, errTLSAuth := a.tryTLSAuth(ctx)
	//if errTLSAuth == nil {
	//	return newCtx, nil
	//}
	newCtx, errTokenAuth := tryTokenAuth(ctx)
	if errTokenAuth == nil {
		return newCtx, nil
	}
	newCtx, errBasicAuth := tryBasicAuth(ctx)
	if errBasicAuth == nil {
		return newCtx, nil
	}
	// TODO return values will change once authentication is required
	return context.TODO(), consts.StatusUnauthenticated
}

// tryEmailTokenVerification checks if browser intends to verify an email.
// The header has a format of "authorization": "Email Token " + token
// Returns a context with no token or an error.
func tryEmailTokenVerification(ctx context.Context) (context.Context, error) {
	log.Info(consts.EmailVerificationTag, consts.StrEmailTokenVerificationAttempt)
	if ctx == nil {
		return context.TODO(), consts.ErrNilContext
	}
	auth, err := extractContextHeader(ctx, consts.StrMdBasicAuthHeader)
	if err != nil {
		log.Error(consts.EmailVerificationTag, err.Error())
		return ctx, err
	}
	if !strings.HasPrefix(auth, consts.StrEmailTokenVerificationPrefix) {
		log.Error(consts.EmailVerificationTag, consts.ErrMissingEmailTokenPrefix.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrMissingEmailTokenPrefix.Error())
	}
	log.Info(consts.EmailVerificationTag, auth)

	err = userSvc.verifyEmailToken(auth[len(consts.StrEmailTokenVerificationPrefix):])
	if err != nil {
		log.Error(consts.EmailVerificationTag, err.Error())
		st, _ := status.FromError(err)
		return ctx, status.Error(st.Code(), st.Message())
	}

	return context.TODO(), nil
}

// tryTokenAuth checks if browser intends to authenticate using an existing auth token.
// The header has a format of "authorization": "Auth Token " + token.
// Returns a context with token or an error.
func tryTokenAuth(ctx context.Context) (context.Context, error) {
	log.Info(consts.TokenAuthTag, consts.StrAuthTokenAttempt)
	if ctx == nil {
		return context.TODO(), consts.ErrNilContext
	}
	auth, err := extractContextHeader(ctx, consts.StrMdBasicAuthHeader)
	if err != nil {
		log.Error(consts.TokenAuthTag, err.Error())
		return ctx, err
	}
	if !strings.HasPrefix(auth, consts.StrTokenAuthPrefix) {
		log.Error(consts.TokenAuthTag, consts.ErrMissingAuthTokenPrefix.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrMissingAuthTokenPrefix.Error())
	}
	log.Info(consts.TokenAuthTag, auth)

	resp, err := userSvc.verifyAuthToken(auth[len(consts.StrTokenAuthPrefix):])
	if err != nil {
		log.Error(consts.TokenAuthTag, err.Error())
		st, _ := status.FromError(err)
		return ctx, status.Error(codes.Unauthenticated, st.Message())
	}

	return finalizeAuth(ctx, consts.TokenAuthTag, resp.Identification)
}

// tryBasicAuth checks if browser intends to authenticate using a base64 encoded "email:password"
// The header has a format of "authorization": "Basic " + base64 encoded "email:password".
// Returns a context with token or an error.
func tryBasicAuth(ctx context.Context) (context.Context, error) {
	log.Info(consts.BasicAuthTag, consts.StrBasicAuthAttempt)
	if ctx == nil {
		return context.TODO(), consts.ErrNilContext
	}
	auth, err := extractContextHeader(ctx, consts.StrMdBasicAuthHeader)
	if err != nil {
		log.Error(consts.BasicAuthTag, err.Error())
		return ctx, err
	}
	if !strings.HasPrefix(auth, consts.StrBasicAuthPrefix) {
		log.Error(consts.BasicAuthTag, consts.ErrMissingBasicAuthPrefix.Error())
		return ctx, status.Error(codes.Unauthenticated, consts.ErrMissingBasicAuthPrefix.Error())
	}
	log.Info(consts.BasicAuthTag, auth)

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
		st, _ := status.FromError(err)
		return ctx, status.Error(codes.Unauthenticated, st.Message())
	}
	return finalizeAuth(ctx, consts.BasicAuthTag, resp.Identification)
}

// finalizeAuth validates the Identification, and sanitizes the context with the token.
// Returns a context with token or an error.
func finalizeAuth(ctx context.Context, tag string, id *lib.Identification) (context.Context, error) {
	if ctx == nil {
		return context.TODO(), consts.ErrNilContext
	}
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
	if ctx == nil {
		return "", consts.ErrNilContext
	}
	if strings.TrimSpace(header) == "" {
		return "", status.Error(codes.Unauthenticated, consts.ErrMissingHeader.Error())
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, consts.ErrMissingHeadersFromCtx.Error())
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
	if ctx == nil {
		return nil, consts.ErrNilContext
	}
	if strings.TrimSpace(header) == "" {
		return nil, consts.ErrMissingHeader
	}
	md, _ := metadata.FromIncomingContext(ctx)
	mdCopy := md.Copy()
	mdCopy[header] = nil
	return metadata.NewIncomingContext(ctx, mdCopy), nil
}
