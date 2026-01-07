package main

import (
    "fmt"
    "os"
    "path"
    
    "dpcms/internal/config"
    "dpcms/internal/httpserver"
    
    "github.com/bytedance/sonic"
    "github.com/spf13/pflag"
)

func registerCommand(appConfig *config.Config) error {
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
        configPath = path.Clean(configPath)
        
        err := os.MkdirAll(path.Dir(configPath), 0755)
        if err != nil {
            return err
        }
        
        _, err = os.Lstat(configPath)
        if err == nil {
            return fmt.Errorf("配置文件 %s 已存在", configPath)
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
        *appConfig = *cfg
    }
    
    return nil
}
