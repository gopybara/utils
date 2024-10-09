package logger

type Params struct {
	SentryDsn   *string
	Environment string
	ServiceName string
	CommitTag   string
}

func (p *Params) GetSentryDsn() *string {
	return p.SentryDsn
}

func (p *Params) GetEnvironment() string {
	return p.Environment
}

func (p *Params) GetServiceName() string {
	return p.ServiceName
}

func (p *Params) GetCommitTag() string {
	return p.CommitTag
}

type Config interface {
	GetSentryDsn() *string
	GetEnvironment() string
	GetServiceName() string
	GetCommitTag() string
}
