package service

import (
	"context"
	"encoding/base64"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"strings"
)

func  Authenticate(ctx context.Context) (newCtx context.Context, err error) {

	// TODO
	//err = a.tryTLSAuth(ctx)
	//if err == nil {
	//	return ctx, nil
	//}
	//
	//newCtx, err = a.tryTokenAuth(ctx)
	//if err == nil {
	//	return newCtx, nil
	//}

	return tryBasicAuth(ctx)
}


// TODO do we need a lock here?
func tryBasicAuth(ctx context.Context) (context.Context, error) {
	// The format:
	// auth := b.email + ":" + b.password
	// enc := base64.StdEncoding.EncodeToString([]byte(auth))
	// "authorization": "Basic " + enc,
	auth, err := extractHeader(ctx, "authorization")
	if err != nil {
		return ctx, err
	}

	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return ctx, status.Error(codes.Unauthenticated, `missing "Basic " prefix in "Authorization" header`)
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, `invalid base64 in header`)
	}

	emailPassword := string(c)
	s := strings.IndexByte(emailPassword, ':')
	if s < 0 {
		return ctx, status.Error(codes.Unauthenticated, `invalid basic auth format`)
	}

	// talk to user service here
	// TODO figure out what to do with resp containing a User object
	// TODO maybe add as metadata in context
	_, err = userSvc.authenticateUser(emailPassword[:s], emailPassword[s+1:])
	if err != nil{
		return ctx, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Remove token from headers from here on
	return purgeHeader(ctx, "authorization"), nil
}

func extractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no header in request")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}

func purgeHeader(ctx context.Context, header string) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	mdCopy := md.Copy()
	mdCopy[header] = nil
	return metadata.NewIncomingContext(ctx, mdCopy)
}
