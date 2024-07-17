package sender

type DingSender struct {
	AccessToken string
	Sec         string
}

func (s *DingSender) Init(msg string, to []string) error {
	return nil
}
