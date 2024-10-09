package configs

type LoggerParams struct {
	SentryDsn   *string
	Environment string
	ServiceName string
	CommitTag   string
}

func (p *LoggerParams) GetSentryDsn() *string {
	return p.SentryDsn
}

func (p *LoggerParams) GetEnvironment() string {
	return p.Environment
}

func (p *LoggerParams) GetServiceName() string {
	return p.ServiceName
}

func (p *LoggerParams) GetCommitTag() string {
	return p.CommitTag
}

type LoggerConfig interface {
	GetSentryDsn() *string
	GetEnvironment() string
	GetServiceName() string
	GetCommitTag() string
}
