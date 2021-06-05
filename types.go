package filedriller

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