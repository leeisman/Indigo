package sub

type SubscriberOption func(*Subscriber)

func SyncMode() SubscriberOption {
	return func(sb *Subscriber) {
		sb.pullMode = syncMode
	}
}

func AsyncMode() SubscriberOption {
	return func(sb *Subscriber) {
		sb.pullMode = asyncMode
	}
}

func SetMaxOutstandingMessages(i int) SubscriberOption {
	return func(sb *Subscriber) {
		sb.subscription.ReceiveSettings.MaxOutstandingMessages = i
	}
}

func SetNumGoroutines(i int) SubscriberOption {
	return func(sb *Subscriber) {
		sb.subscription.ReceiveSettings.NumGoroutines = i
	}
}

func SetErrorHook(errorHook func(err error)) SubscriberOption {
	return func(sb *Subscriber) {
		sb.errorHook = errorHook
	}
}
