package filedriller

// RedisConf holds the config of a redis server
type RedisConf struct {
	Server *string
	Port   *string
}

// Config maps all flags to a struct. This is used in griller, the filedriller GUI
type Config struct {
	RootDir string
	HashAlg string
	RedisServer string
	RedisPort string
	SFile bool
	OFile string
	IFile string
	Entro bool
}