package main

import (
	"github.com/heavenfollowsman/MinecraftCmdPlugin/MinecraftCmdPlugin"
	"github.com/kohmebot/kohme"
	"github.com/kohmebot/kohme/pkg/conf"
)

func main() {

	// 读取配置文件
	var zeroConf = new(conf.ZeroConf)
	err := zeroConf.ParseJsonFile("L:\\All_Project\\GO\\XJH\\KohmeBot\\kohme\\conf\\config.json")
	if err != nil {
		panic(err)
	}

	var pluginConf = new(conf.PluginConf)
	pluginConf.ParseYamlFile("L:\\All_Project\\GO\\XJH\\KohmeBot\\kohme\\conf\\plugins.yaml")

	kohme.Register(MinecraftCmdPlugin.NewPlugin())

	kohme.RunKohme(*zeroConf, *pluginConf)
}
