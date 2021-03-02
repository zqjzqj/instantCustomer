package sErr

type SErr struct {
	err string
}

func New(msg string) *SErr {
	return &SErr{err:msg}
}

func NewByError(err error) *SErr {
	return &SErr{err:err.Error()}
}

func (errS *SErr) Error() string {
	return errS.err
}
