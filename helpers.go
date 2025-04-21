package echoext

// M is a helper type for map[string]any
type M map[string]any

func ErrM(err error) M {
	return M{"error": err.Error()}
}
