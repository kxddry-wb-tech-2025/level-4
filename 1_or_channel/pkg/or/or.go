package or

// Or returns a channel that closes when one of the channels in the argument list closes.
// It returns nil if no channels are provided.
func Or(channels ...<-chan struct{}) <-chan struct{} {
	// if no channels are provided, return nil
	switch len(channels) {
	case 0:
		return nil
	// if only one channel is provided, return it
	case 1:
		return channels[0]
	default:
	}

	orDone := make(chan struct{})

	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-Or(append(channels[1:], orDone)...):
			}
		}
	}()

	return orDone
}
