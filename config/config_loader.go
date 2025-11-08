package config

type Loader interface {
	Init(configStructPointer any) error        // 初始化方法,相关参数走环境变量和启动命令或者文件都可以
	SetCallback(callback ReloadConfigCallback) // 进行动态调用callback的
}

type ReloadConfigCallback func(configStruct interface{}) error
