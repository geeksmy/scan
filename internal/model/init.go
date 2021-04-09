package model

const (
	prefix        = "scan_"
	TabNamePort   = prefix + "port"       // 端口扫描表
	TabNameMatch  = prefix + "match"      // 指纹正则表
	TabNameProbe  = prefix + "probe"      // 指纹探针表
	TabWeb        = prefix + "web"        // web 指纹识别表
	TabCyberspace = prefix + "cyberspace" // 网络空间扫描表
	TabBlasting   = prefix + "blasting"   // 密码爆破表
)
