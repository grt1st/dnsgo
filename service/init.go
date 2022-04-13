package service

var Config Settings
var R Counter

func Init(filepath string) {
	//_, err := toml.DecodeFile(filepath + filename, &Config)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	Config.FilePath = filepath
	R = Counter{
		mapCount: map[string]int{},
	}
}

type FileChecker struct {
	Filename      string
	FileTimestamp int64
}

type Settings struct {
	FilePath   string
	Nameserver string
	Hosts      string
	Resolv     string
}
