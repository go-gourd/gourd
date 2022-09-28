package sessions

type SessionConfig struct {
	Type   string
	Path   string //文件路径，file方式可用
	Expire int    //过期时间，秒
	Domain string
	Secure bool
}
