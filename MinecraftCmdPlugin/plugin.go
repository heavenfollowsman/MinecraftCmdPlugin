package MinecraftCmdPlugin

import (
	"encoding/json"
	"fmt"
	"github.com/gorcon/rcon"
	"github.com/kohmebot/pkg/command"
	"github.com/kohmebot/pkg/version"
	"github.com/kohmebot/plugin"
	zero "github.com/wdvxdr1123/ZeroBot"
	"os"
	"strconv"
	"strings"
)

type PluginMinecraftCmdPlugin struct {
	env plugin.Env
	// 配置文件
	conf Config
	// 可执行指令的人员map
	permissionUsers map[int64]bool
}

func NewPlugin() plugin.Plugin {
	return new(PluginMinecraftCmdPlugin)
}

// 插件初始化方法
func (p *PluginMinecraftCmdPlugin) Init(engine *zero.Engine, env plugin.Env) error {
	p.env = env
	// 读取插件的配置文件内容
	err := env.GetConf(&p.conf)
	if err != nil {
		return err
	}

	// 加载授权数据
	p.loadAuthorized()

	engine.OnCommand("MC").Handle(func(ctx *zero.Ctx) {

		//ctx.Send("读取到的默认用户配置信息为：" + strconv.FormatInt(p.conf.BaseAdmin, 10))
		qq := ctx.Event.UserID
		//ctx.Send("当前发送消息的人qq为" + strconv.FormatInt(qq, 10))
		qqName := ctx.Event.Sender.NickName
		msg := ctx.Event.Message.String()

		/*
			func TrimPrefix(A, B string) string
			检查A是否以B开头，且如果是就移除A中的B，如果不是就返回原本的值

			func TrimSpace(s string) string
			移除字符串 s 两端的空白字符（包括空格、制表符 \t、换行符 \n 等），返回处理后的字符串。
		*/
		cmd := strings.TrimSpace(strings.TrimPrefix(msg, "/MC"))
		// 权限校验
		if !p.permissionUsers[qq] {
			ctx.Send("住口！" + qqName + "！你这皓首匹夫，苍髯老贼！你没资格命令我")
			return
		}

		/*
			func HasPrefix(s, prefix string) bool
			检查字符串 cmd 是否以 "QQadd " 开头，返回 true（如果以该前缀开头）或 false（如果不以该前缀开头）
		*/
		// 添加授权
		if strings.HasPrefix(cmd, "QQadd ") && qq == p.conf.BaseAdmin {
			newUser := strings.TrimSpace(strings.TrimPrefix(cmd, "QQadd "))
			// 转换为整型
			num, err := strconv.ParseInt(newUser, 10, 64)

			if err != nil {
				ctx.Send("添加失败，QQ号格式不正确")
				return
			}
			p.permissionUsers[num] = true
			// 持久化改用户权限
			p.saveAuthorized()
			ctx.Send(fmt.Sprintf("已授权用户 %d 可使用 /MC 指令", num))
			return
		}

		// 删除授权
		if strings.HasPrefix(cmd, "QQdel ") && qq == p.conf.BaseAdmin {
			delUser := strings.TrimSpace(strings.TrimPrefix(cmd, "QQdel "))
			num, err := strconv.ParseInt(delUser, 10, 64)
			if err != nil {
				ctx.Send("删除失败，QQ号格式不正确")
				return
			}
			// 从nap中删除该人员
			delete(p.permissionUsers, num)
			p.saveAuthorized()
			ctx.Send(fmt.Sprintf("已移除用户 %d 的使用权限", num))
			return
		}

		// 执行 MC 指令
		if cmd == "" {
			ctx.Send("请附带要发送的 Minecraft 指令")
			return
		}

		for _, keyword := range dangerousPrefixes {
			if strings.Contains(strings.ToLower(cmd), strings.ToLower(keyword)) {
				ctx.Send("该指令包含高危关键字，禁止执行")
				return
			}
		}

		conn, err := rcon.Dial(p.conf.RconAddress, p.conf.RconPassword)
		if err != nil {
			ctx.Send(fmt.Sprintf("连接 RCON 失败: %v", err))
			return
		}
		defer conn.Close()

		resp, err := conn.Execute("/" + cmd)
		if err != nil {
			ctx.Send(fmt.Sprintf("命令执行失败: %v", err))
			return
		}

		ctx.Send(fmt.Sprintf("Minecraft 返回结果: %s", resp))
	})

	return nil
}

func (p *PluginMinecraftCmdPlugin) Name() string {
	return "MinecraftCmdPlugin"
}

func (p *PluginMinecraftCmdPlugin) Description() string {
	return "本插件是通过QQ机器人来远程执行服务器命令，实现在群内进行一定程度上的运维。"
}

func (p *PluginMinecraftCmdPlugin) Commands() fmt.Stringer {
	return command.NewCommands(
		command.NewCommand("发送/MC即可输入MC的具体指令，如/MC list 查看当前在线人员。", "MC"))
}

func (p *PluginMinecraftCmdPlugin) Version() uint64 {
	return uint64(version.NewVersion(0, 0, 30))
}

func (p *PluginMinecraftCmdPlugin) OnBoot() {

}

func (p *PluginMinecraftCmdPlugin) loadAuthorized() {
	// 实例化map
	p.permissionUsers = make(map[int64]bool)

	// 打开文件
	file, err := os.Open(p.conf.PermissionFilePath)
	if err != nil {
		// 文件不存在就默认给 baseAdmin 授权
		fmt.Println("权限文件不存在，初始化默认权限")
		fmt.Println(p.conf.BaseAdmin)
		p.permissionUsers[p.conf.BaseAdmin] = true
		p.saveAuthorized()
		return
	}
	defer file.Close()

	// 将上面的map编码为json格式并写入指定的json文件
	json.NewDecoder(file).Decode(&p.permissionUsers)
}

func (p *PluginMinecraftCmdPlugin) saveAuthorized() {

	// 保存新授权人员
	file, err := os.Create(p.conf.PermissionFilePath)
	if err != nil {
		fmt.Println("无法保存权限文件:", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(p.permissionUsers)
}

// 危险指令
var dangerousPrefixes = []string{
	"kill @e", "kill @a", "stop", "op", "deop",
	"ban", "ban-ip", "whitelist off", "clear @a",
	"data", "structure", "clone", "fill", "tp @a",
	"kill",
}
