package pub

type PublisherOption func(*Publisher)

func SetErrorHook(errorHook func(err error, requestID string)) PublisherOption {
	return func(p *Publisher) {
		p.errorHook = errorHook
	}
}
