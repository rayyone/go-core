package stgoption

// Option Function to change request options
type OptionFunc func(*Option)

type Option struct {
	ContentDisposition string
}

func ContentDisposition(value string) OptionFunc {
	return func(opt *Option) {
		opt.ContentDisposition = value
	}
}

func GetDefaultOptions() Option {
	return Option{
		ContentDisposition: "inline",
	}
}
