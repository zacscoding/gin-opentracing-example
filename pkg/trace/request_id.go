package trace

import "context"

type contextKey = string

const requestIdKey = contextKey("requestId")

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}

func ExtractRequestId(ctx context.Context) string {
	requestId, ok := ctx.Value(requestIdKey).(string)
	if ok {
		return requestId
	}
	return ""
}
