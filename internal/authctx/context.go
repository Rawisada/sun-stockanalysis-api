package authctx

import "context"

type contextKey struct{}

var userIDContextKey contextKey

func UserIDContextKey() any {
	return userIDContextKey
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}
