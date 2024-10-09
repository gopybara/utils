package redis

type Params struct {
	DSN      string
	Host     string
	Port     string
	Password string
	DB       int
}

func (c *Params) GetDSN() string {
	if c.DSN == "" {
		c.DSN = c.Host + ":" + c.Port
	}

	return c.DSN
}

func (c *Params) GetPassword() string {
	return c.Password
}

func (c *Params) GetDB() int {
	return c.DB
}

type Config interface {
	GetDSN() string
	GetPassword() string
	GetDB() int
}

type SentinelParams struct {
	MasterName       string
	Nodes            []string
	Password         string
	SentinelPassword string
	DB               int
}

type SentinelConfig interface {
	GetMasterName() string
	GetNodes() []string
	GetPassword() string
	GetSentinelPassword() string
	GetDB() int
}

func (c *SentinelParams) GetMasterName() string {
	return c.MasterName
}

func (c *SentinelParams) GetNodes() []string {
	return c.Nodes
}

func (c *SentinelParams) GetPassword() string {
	return c.Password
}

func (c *SentinelParams) GetSentinelPassword() string {
	return c.SentinelPassword
}

func (c *SentinelParams) GetDB() int {
	return c.DB
}
