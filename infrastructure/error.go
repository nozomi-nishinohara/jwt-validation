package infrastructure

const (
	ErrNotFound = JwtValudationError("Not found")
)

type JwtValudationError string

func (e JwtValudationError) Error() string     { return string(e) }
func (JwtValudationError) JwtValudationError() {}
