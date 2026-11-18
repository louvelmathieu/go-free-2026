package table

import "errors"

var ErrTransient = errors.New("erreur passagère")

// IsRetryable est la fonction sous test : une erreur passagère, wrappée ou non,
// se retente. Les autres non.
func IsRetryable(err error) bool { return errors.Is(err, ErrTransient) }
