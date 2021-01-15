package utils

// Result : a monad help solve error chain
type Result struct {
	context interface{}
	err     error
}

// NewResult : create a result object
func NewResult(context interface{}, err error) *Result {
	result := &Result{
		context: context,
		err:     err,
	}
	return result
}

// AndThen : if earlier functions have failed, then do nothing, else run next function
func (result *Result) AndThen(next func(interface{}) (interface{}, error)) *Result {
	if result.err != nil {
		return result
	}
	result.context, result.err = next(result.context)
	return result
}

// Unwrap : read contents of result
func (result *Result) Unwrap() (interface{}, error) {
	return result.context, result.err
}

// Error : read error of result
func (result *Result) Error() error {
	return result.err
}

// Must : read context of result
func (result *Result) Must() interface{} {
	if result.err != nil {
		return nil
	}
	return result.context
}

// Try : do a chain of actions which may produce error and stop when error occured
func Try(actions ...func() error) error {
	for _, action := range actions {
		if err := action(); err != nil {
			return err
		}
	}
	return nil
}
