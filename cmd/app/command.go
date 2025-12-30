package app

import (
    "fmt"
    "os"
    
    app2 `dpcms/cmd/app`
    "dpcms/internal/config"
    "dpcms/internal/httpserver"
    
    "github.com/bytedance/sonic"
    "github.com/spf13/pflag"
)

func registerCommand() error {
    var (
        doPrintVersion bool
        doCreateExFile bool
        configPath     string
    )
    
    pflag.BoolVarP(&doPrintVersion, "version", "v", false, "版本号")
    pflag.BoolVarP(&doCreateExFile, "example", "e", false, "创建配置文件")
    pflag.StringVarP(&configPath, "conf", "c", "./runtime/config.json", "配置文件路径")
    pflag.Parse()
    
    if doPrintVersion {
        fmt.Println(httpserver.Version)
        os.Exit(0)
    }
    
    if doCreateExFile {
        _, oErr := os.Lstat("./runtime/config.json")
        if oErr == nil {
            return fmt.Errorf("配置文件 runtime/config.json 已存在")
        }
        
        f, err := os.OpenFile("./runtime/config.json", os.O_WRONLY|os.O_CREATE, 0644)
        defer func(f *os.File) {
            _ = f.Close()
        }(f)
        
        if err != nil {
            return err
        }
        
        json, err := sonic.MarshalIndent(config.NewWithBasicConfig(), "", "  ")
        if err != nil {
            return err
        }
        _, err = f.Write(json)
        if err != nil {
            return err
        }
        
        fmt.Println("示例文件config.json已创建")
        
        os.Exit(0)
    }
    
    if configPath != "" {
        cfg, err := config.New(configPath)
        if err != nil {
            panic(err)
        }
        app2.appConfig = cfg
    }
    
    return nil
}
