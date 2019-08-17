package tgclient

type Builder struct {
	config config
}

type config struct {
	apiId              int
	apiHash            string
	authPhone          string
	systemLanguageCode string
	systemVersion      string
	deviceModel        string
	applicationVersion string
	filesDirectory     string
	databaseDirectory  string
	useFileDatabase    bool
	checkCode          string
	password           string
	proxy              *socks5Proxy
}

type socks5Proxy struct {
	host     string
	port     int
	login    string
	password string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) ApiId(val int) *Builder {
	b.config.apiId = val
	return b
}

func (b *Builder) ApiHash(val string) *Builder {
	b.config.apiHash = val
	return b
}

func (b *Builder) AuthPhone(val string) *Builder {
	b.config.authPhone = val
	return b
}

func (b *Builder) SystemLanguageCode(val string) *Builder {
	b.config.systemLanguageCode = val
	return b
}

func (b *Builder) SystemVersion(val string) *Builder {
	b.config.systemVersion = val
	return b
}

func (b *Builder) DeviceModel(val string) *Builder {
	b.config.deviceModel = val
	return b
}

func (b *Builder) ApplicationVersion(val string) *Builder {
	b.config.applicationVersion = val
	return b
}

func (b *Builder) FilesDirectory(val string) *Builder {
	b.config.filesDirectory = val
	return b
}

func (b *Builder) DatabaseDirectory(val string) *Builder {
	b.config.databaseDirectory = val
	return b
}

func (b *Builder) UseFileDatabase(val bool) *Builder {
	b.config.useFileDatabase = val
	return b
}

func (b *Builder) Socks5Proxy(host string, port int, login, password string) *Builder {
	b.config.proxy = &socks5Proxy{
		host:     host,
		port:     port,
		login:    login,
		password: password,
	}
	return b
}

func (b *Builder) CheckCode(val string) *Builder {
	b.config.checkCode = val
	return b
}

func (b *Builder) Password(val string) *Builder {
	b.config.password = val
	return b
}

func (b *Builder) Build() *Client {
	return newClient(b.config)
	//if b.proxy != nil {
	//	_, err := cl.addProxy(*b.proxy)
	//	if err != nil {
	//		cl.Destroy()
	//return nil, err
	//}
	//}
}
