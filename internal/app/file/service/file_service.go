package service

import (
    `context`
    `crypto`
    `encoding/hex`
    `errors`
    `io`
    `io/fs`
    `path`
    `time`
    
    `cms/internal/app/file/internal/assembler`
    `cms/internal/app/file/internal/dto`
    `cms/internal/app/file/internal/errno`
    `cms/internal/app/file/internal/fileutil`
    `cms/internal/infra/logger`
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    
    `github.com/avast/retry-go`
)

type FileService struct {
    txManager contract.Transactor
    repo      contract.FileRepo
    logger    *logger.Logger
    driverSrv *FileDriverService
}

func NewFileService(txManager contract.Transactor, repo contract.FileRepo, driverSrv *FileDriverService, logger *logger.Logger) *FileService {
    return &FileService{
        txManager: txManager,
        repo:      repo,
        driverSrv: driverSrv,
        logger:    logger,
    }
}

func (s FileService) Save(ctx context.Context, user *model.User, params dto.FileSaveCommand) (*dto.FileInfo, error) {
    driverName, driver, err := s.driverSrv.GetCurrentDriver()
    if err != nil {
        return nil, err
    }
    
    var fileExt = path.Ext(params.Name)
    var fileWriter io.WriteCloser
    var savingKey string
    
    err = retry.Do(
        func() error {
            genPath, retryErr := fileutil.GeneratePath()
            if retryErr != nil {
                return retryErr
            }
            
            savingKey = genPath.FullPath() + fileExt
            writer, retryErr := driver.OpenWriter(ctx, savingKey)
            if retryErr != nil {
                if writer != nil {
                    _ = writer.Close()
                }
                if errors.Is(retryErr, fs.ErrExist) {
                    return retryErr
                }
                return retry.Unrecoverable(retryErr)
            }
            
            fileWriter = writer
            return nil
        },
        retry.Context(ctx),
        retry.Attempts(3),
    )
    
    if err != nil {
        return nil, errno.FilePreCreateFileFailed.Wrap(err).ToError()
    }
    
    hash := crypto.SHA256.New()
    multiWriter := io.MultiWriter(fileWriter, hash)
    fileSize, err := io.Copy(multiWriter, params.Data)
    
    delCtx, cancel := context.WithTimeout(ctx, time.Second*3)
    defer cancel()
    
    if err != nil {
        _ = fileWriter.Close()
        _ = driver.Delete(delCtx, savingKey)
        return nil, err
    }
    
    if err = fileWriter.Close(); err != nil {
        _ = driver.Delete(delCtx, savingKey)
        return nil, err
    }
    
    fileRecord := &model.File{
        UserID:       user.ID,
        OriginalName: params.Name,
        Ext:          fileExt,
        Path:         savingKey,
        Size:         datatype.SafeInt64(fileSize),
        SHA256:       hex.EncodeToString(hash.Sum(nil)),
        Driver:       driverName,
    }
    fileRecord.IsImage.FromBool(fileutil.IsImage(fileExt))
    
    err = s.repo.Create(ctx, fileRecord)
    if err != nil {
        if delErr := driver.Delete(delCtx, savingKey); delErr != nil {
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

func (s FileService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    fileRecord, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    driver, err := s.driverSrv.GetDriver(fileRecord.Driver)
    if err != nil {
        return err
    }
    if _, err = s.repo.Delete(ctx, id); err != nil {
        return err
    }
    return driver.Delete(ctx, fileRecord.Path)
}

func (s FileService) openFileReaderByRecord(ctx context.Context, fileRecord *model.File) (io.ReadCloser, *dto.FileMeta, error) {
    driver, err := s.driverSrv.GetDriver(fileRecord.Driver)
    if err != nil {
        return nil, nil, err
    }
    reader, err := driver.OpenReader(ctx, fileRecord.Path)
    if err != nil {
        return nil, nil, err
    }
    return reader, assembler.ToFileMetaDTO(fileRecord), nil
}

func (s FileService) OpenFileRecordByID(ctx context.Context, id datatype.SafeUint64) (io.ReadCloser, *dto.FileMeta, error) {
    fileRecord, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, nil, err
    }
    return s.openFileReaderByRecord(ctx, fileRecord)
}

func (s FileService) OpenFileRecordByPath(ctx context.Context, filePath string) (io.ReadCloser, *dto.FileMeta, error) {
    fileRecord, err := s.repo.FindByPath(ctx, filePath)
    if err != nil {
        return nil, nil, err
    }
    return s.openFileReaderByRecord(ctx, fileRecord)
}
