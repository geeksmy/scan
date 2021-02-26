package model

const (
	prefix             = "scan_"
	TabNamePort        = prefix + "port"         // 端口扫描表
	TabNamePortService = prefix + "port_service" // 端口服务表
	TabNameMatch       = prefix + "match"        // 指纹正则表
	TabNameProbe       = prefix + "probe"        // 指纹探针表
)
