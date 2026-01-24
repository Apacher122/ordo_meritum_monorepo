package contexts

import "context"

// FromContext retrieves a UserContext from the given context. If no UserContext is found,
// it returns nil and false. Otherwise, it returns the UserContext and true.
func FromContext(ctx context.Context) (*UserContext, bool) {
	userCtx, ok := ctx.Value(UserContextKey).(*UserContext)
	return userCtx, ok
}
