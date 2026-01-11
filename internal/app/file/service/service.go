package service

import (
    `bytes`
    `context`
    `crypto`
    `encoding/hex`
    `errors`
    `io`
    `io/fs`
    `path`
    `time`
    
    configkeys `dpcms/internal/app/config/constant`
    config `dpcms/internal/app/config/service`
    `dpcms/internal/app/file/errno`
    `dpcms/internal/app/file/internal/assembler`
    `dpcms/internal/app/file/internal/dto`
    `dpcms/internal/app/file/internal/fileutil`
    `dpcms/internal/erroz`
    `dpcms/internal/infra/file`
    `dpcms/internal/infra/logger`
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    
    `github.com/avast/retry-go`
)

type FileService struct {
    repo      FileRepo
    txManager contract.TxManager
    registry  *file.DriverRegistry
    configSrv *config.ConfigService
    logger    *logger.Logger
}

func NewFileService(txManager contract.TxManager, repo FileRepo, registry *file.DriverRegistry, configSrv *config.ConfigService, logger *logger.Logger) *FileService {
    return &FileService{
        txManager: txManager,
        repo:      repo,
        registry:  registry,
        configSrv: configSrv,
        logger:    logger,
    }
}

func (s FileService) getCurrentFileDriver() (file.Driver, error) {
    driverName, exists := s.configSrv.Get(configkeys.FileDriver)
    if !exists {
        return nil, erroz.Unknown.Wrap(
            errno.FileDriverConfigNotExists.ToError(),
        ).ToError()
    }
    
    driver, exists := s.registry.Get(driverName)
    if !exists {
        return nil, erroz.Unknown.Wrap(
            errno.FileDriverNotExists.ToError(),
        ).ToError()
    }
    
    return driver, nil
}

func (s FileService) Save(ctx context.Context, user *model.User, params dto.FileSaveCommand) (*dto.FileInfo, error) {
    driver, err := s.getCurrentFileDriver()
    if err != nil {
        return nil, err
    }
    
    genPath, err := fileutil.GeneratePath()
    if err != nil {
        return nil, err
    }
    
    fileData, err := io.ReadAll(params.Data)
    if err != nil {
        return nil, err
    }
    
    fileExt := path.Ext(params.Filename)
    savingKey := genPath.FullPath() + fileExt
    
    var fileSize int64
    var fileHash string
    
    err = retry.Do(func() error {
        writer, retryErr := driver.OpenWriter(ctx, savingKey)
        if retryErr != nil {
            if errors.Is(retryErr, fs.ErrExist) {
                return retryErr
            }
            return retry.Unrecoverable(retryErr)
        }
        
        hash := crypto.SHA256.New()
        writers := io.MultiWriter(writer, hash)
        written, retryErr := io.Copy(writers, bytes.NewReader(fileData))
        delCtx, cancel := context.WithTimeout(ctx, time.Second*5)
        defer cancel()
        
        if retryErr != nil {
            _ = writer.Close()
            _ = driver.Delete(delCtx, savingKey)
            return retryErr
        }
        
        if retryErr = writer.Close(); retryErr != nil {
            _ = driver.Delete(delCtx, savingKey)
            return retry.Unrecoverable(retryErr)
        }
        
        fileSize = written
        fileHash = hex.EncodeToString(hash.Sum(nil))
        
        return nil
    }, retry.Context(ctx), retry.Attempts(3))
    
    if err != nil {
        return nil, errno.FileGeneratePathFailed.Wrap(err).ToError()
    }
    
    fileRecord := &model.File{
        UserID:  user.ID,
        Name:    genPath.FullFilename(fileExt),
        Ext:     fileExt,
        Path:    savingKey,
        Size:    datatype.SafeInt64(fileSize),
        IsImage: 0,
        SHA256:  fileHash,
        Driver:  driver.Name(),
    }
    
    if fileutil.IsImage(fileExt) {
        fileRecord.IsImage = 1
    }
    
    err = s.repo.Create(ctx, fileRecord)
    if err != nil {
        delErr := driver.Delete(ctx, savingKey)
        if delErr != nil {
            s.logger.App.Error(
                errno.FileDeleteFailedInCreating.Wrap(delErr).ToError().Error(),
            )
        }
        return nil, err
    }
    
    url, err := driver.URL(ctx, savingKey)
    if err != nil {
        return nil, err
    }
    
    return assembler.ToFileInfoDTO(fileRecord, url), nil
}
