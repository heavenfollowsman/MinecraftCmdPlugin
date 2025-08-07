package MinecraftCmdPlugin

type Config struct {
	RconAddress  string `yaml:"rcon_address"`
	RconPassword string `yaml:"rcon_password"`
	// 基础的MC命令管理员
	BaseAdmin int64 `yaml:"base_admin"`
	// 授权文件路径
	PermissionFilePath string `yaml:"permission_file_path"`
}
