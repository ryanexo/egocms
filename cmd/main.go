package main

import (
    _ `embed`
    `fmt`
    `os`

    `GoBlog/internal/server`
    `GoBlog/internal/server/config`
    `github.com/spf13/pflag`
)

func main() {
    initCommand()
    cfg, err := config.New(configPath)
    if err != nil {
        panic(err)
    }
    s, err := initServer(cfg)
    if err != nil {
        panic(err)
    }
    if err = s.Listen(); err != nil {
        panic(err)
    }
}

//go:embed example.json
var exampleConfig string
var (
    doPrintVersion bool
    doCreateExFile bool
    configPath     string
)

func initCommand() {
    pflag.BoolVarP(&doPrintVersion, "version", "v", false, "版本号")
    pflag.BoolVarP(&doCreateExFile, "example", "e", false, "创建配置文件")
    pflag.StringVarP(&configPath, "conf", "c", "./config/config.json", "配置文件路径")
    pflag.Parse()

    if doPrintVersion {
        fmt.Println(server.Version)
        os.Exit(0)
    }

    if doCreateExFile {
        _, oErr := os.Lstat("./config/config.json")
        if oErr == nil {
            fmt.Println("配置文件 config/config.json 已存在")
            return
        }

        f, err := os.OpenFile("./config/config.json", os.O_WRONLY|os.O_CREATE, 0644)
        if err != nil {
            fmt.Println(err)

            return
        }

        _, err = f.WriteString(exampleConfig)
        if err != nil {
            fmt.Println(err)

            return
        }

        _ = f.Close()
        fmt.Println("示例文件config.json已创建")

        os.Exit(0)
    }
}
