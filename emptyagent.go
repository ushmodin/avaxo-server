package avaxo

type EmptyAgent struct {
}

func (self EmptyAgent) Notify(opts ForwardOpts) error {
	return nil
}
