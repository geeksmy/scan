package zap

/*
# default yaml file
Logger:
    Level: info
*/
// Conf zap 的配置
type Conf struct {
	Level string
}

var C = Conf{Level: "info"}
