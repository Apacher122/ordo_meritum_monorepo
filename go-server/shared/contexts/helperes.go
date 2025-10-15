package contexts

import "context"

func FromContext(ctx context.Context) (*UserContext, bool) {
	userCtx, ok := ctx.Value(UserContextKey).(*UserContext)
	return userCtx, ok
}
