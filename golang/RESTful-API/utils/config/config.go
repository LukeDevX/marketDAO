// 配置文件方法封装
package config

import (
	"fmt"             // 用于格式化 I/O，比如输出错误信息
	"gopkg.in/ini.v1" //  一个处理 .ini 配置文件的库，能够方便地读取、修改 .ini 文件中的键值。
	"io/ioutil"       // 用于读取目录中的文件。
	"os"              // 提供与操作系统交互的功能，例如文件系统操作、获取工作目录等。
	"path/filepath"   // 用于处理文件路径的库，提供跨平台的路径操作功能。
	"strings"         // 提供对字符串的处理功能，如分割字符串。
)

var (
	// 是一个全局变量，类型为 *ini.File，用于存储 .ini 文件的配置内容。
	appConfig *ini.File
)

// init() 是 Go 中的特殊函数，它会在程序启动时自动运行，用于初始化一些全局设置。在这个例子中，它用于读取和加载配置文件。
func init() {
	var (
		err            error        // 用于捕获错误信息
		configsDirname = "/configs" // 定义配置文件所在的目录
		configPath     string       // 存储实际的配置文件路径。
	)

	workPath, err := os.Getwd() // 使用 os.Getwd() 函数获取存储当前工作目录。如果获取当前工作目录失败，程序会直接崩溃（使用 panic）。
	//fmt.Println("workPath", workPath)
	if err != nil {
		panic(err)
	}
	configPath = workPath + configsDirname // 将当前工作目录与 configs 目录拼接，构成配置文件路径。

	if fileExists(configPath) == false { // 如果当前工作目录下的 configs 目录不存在，尝试使用可执行文件所在的目录作为基准路径。
		execPath, err := os.Executable() // os.Executable() 获取当前可执行文件的路径
		if err != nil {
			panic(err)
		}
		configPath = filepath.Dir(execPath) + configsDirname // 将可执行文件所在的目录与 configs 目录拼接，构成配置文件路径。filepath.Dir(execPath) 返回其所在的目录路径。
	}
	files, err := ioutil.ReadDir(configPath) // ioutil.ReadDir(configPath) 读取 configPath 下的所有文件列表，如果读取失败，程序会崩溃。
	if err != nil {
		panic(err)
	}
	appConfig = ini.Empty() // 初始化 appConfig 为一个空的 ini 文件，用于存储配置文件内容。
	// 遍历 files，对于每个文件，使用 strings.Split(file.Name(), ".") 分割文件名和扩展名。
	// 检查文件是否有 .ini 扩展名，如果是，调用 appConfig.Append 将 .ini 文件加载到 appConfig 中。
	// 如果读取 .ini 文件失败，会打印错误信息并退出程序。
	for _, file := range files {
		fileInfo := strings.Split(file.Name(), ".")
		if len(fileInfo) == 2 && fileInfo[1] == "ini" {
			err = appConfig.Append(configPath + "/" + file.Name())
			if err != nil {
				fmt.Printf("Fail to read file:%s,err: %v", file.Name(), err)
				os.Exit(1)
			}
		}
	}
}

// 该函数用于检查某个路径是否存在。
// os.Stat(path) 获取文件或目录的状态信息。如果没有报错，说明文件/目录存在，返回 true。
// 如果报错并且错误类型不是 os.IsExist，则返回 false。
// 判断所给路径文件/文件夹是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// GetConfig 用于获取配置项的值，支持按 RunMode 分不同的环境（如生产环境、开发环境等）来获取配置。
// 读取 RunMode 配置
// appConfig.Section("").Key("RunMode").String() 获取全局 RunMode 配置，如果没有设置，默认值为 "prod"（生产模式）。
// 获取对应的配置项
// 如果当前 RunMode 下的配置中有指定的 keyName，则返回该配置项。
// 否则，返回全局配置中 keyName 对应的值。
func GetConfig(keyName string) *ini.Key {
	runMode := appConfig.Section("").Key("RunMode").String()
	if runMode == "" {
		runMode = "prod"
	}
	//如果runMode下不存在，从非runMode取
	if appConfig.Section(runMode).HasKey(keyName) == false {
		return appConfig.Section("").Key(keyName)
	}
	return appConfig.Section(runMode).Key(keyName)
}
