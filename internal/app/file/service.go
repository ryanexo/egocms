package file

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"
	
	"cms/internal/app/file/api"
	"cms/internal/app/file/errno"
	"cms/internal/app/file/model"
	"cms/internal/infra/store/modeltype"
	"cms/internal/public/apitype"
	"cms/internal/public/jsontype"
	
	"github.com/google/uuid"
)

type Repository interface {
    CreateFile(context.Context, *model.File) error
    DeleteFileIfUnused(context.Context, uint64) error
    RestoreFile(context.Context, uint64) error
    FindFileByID(context.Context, uint64) (*model.File, error)
    ListFiles(context.Context, *api.FileListParams) ([]*model.File, int64, error)
    CreateAttachment(context.Context, *model.Attachment) error
    UpdateAttachment(context.Context, *model.Attachment) error
    DeleteAttachment(context.Context, uint64) error
    FindAttachmentByID(context.Context, uint64) (*model.Attachment, error)
    ListAttachments(context.Context, *api.AttachmentListParams) ([]*model.Attachment, int64, error)
}

type Service struct {
    repo     Repository
    registry *manager.DriverRegistry
}

type countingReader struct {
    io.Reader
    n int64
}

func (r *countingReader) Read(p []byte) (int, error) {
    n, err := r.Reader.Read(p)
    r.n += int64(n)
    return n, err
}

func NewFileService(repo Repository, registry *manager.DriverRegistry) *Service {
    return &Service{repo: repo, registry: registry}
}

func (service *Service) Upload(ctx context.Context, originalName string, source io.Reader) (*api.File, error) {
    target, _ := service.registry.Default()
    return service.UploadTo(ctx, target, originalName, source)
}

func (service *Service) UploadTo(ctx context.Context, target, originalName string, source io.Reader) (*api.File, error) {
    originalName = path.Base(strings.ReplaceAll(originalName, "\\", "/"))
    if originalName == "." || originalName == "/" {
        originalName = ""
    }
    ext := strings.ToLower(path.Ext(originalName))
    record := &model.File{}
    if err := record.SetOriginalName(originalName); err != nil {
        return nil, err
    }
    if err := record.SetExt(ext); err != nil {
        return nil, err
    }
    driver, err := service.registry.Select(target)
    if err != nil {
        return nil, storageError(err)
    }
    storagePath := path.Join(time.Now().Format("2006/0102"), uuid.NewString()+ext)
    if err = record.SetStorage(target, storagePath); err != nil {
        return nil, err
    }
    
    hash := sha256.New()
    counting := &countingReader{Reader: io.TeeReader(source, hash)}
    err = driver.Create(ctx, storagePath, counting, -1)
    if err != nil {
        _ = driver.Delete(context.WithoutCancel(ctx), storagePath)
        return nil, storageError(err)
    }
    record.Size = uint64(counting.n)
    record.SHA256 = [32]byte(hash.Sum(nil))
    if err = service.repo.CreateFile(ctx, record); err != nil {
        cleanupErr := driver.Delete(context.WithoutCancel(ctx), storagePath)
        return nil, errors.Join(err, cleanupErr)
    }
    return service.fileDTO(ctx, record)
}

func (service *Service) FindFileByID(ctx context.Context, id uint64) (*api.File, error) {
    record, err := service.repo.FindFileByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return service.fileDTO(ctx, record)
}

func (service *Service) ListFiles(ctx context.Context, params api.FileListParams) (*apitype.PaginatedResult[*api.File], error) {
    records, total, err := service.repo.ListFiles(ctx, &params)
    if err != nil {
        return nil, err
    }
    list := make([]*api.File, 0, len(records))
    for _, record := range records {
        item, dtoErr := service.fileDTO(ctx, record)
        if dtoErr != nil {
            return nil, dtoErr
        }
        list = append(list, item)
    }
    return &apitype.PaginatedResult[*api.File]{Pagination: params.Pagination, Total: total, List: list}, nil
}

func (service *Service) Open(ctx context.Context, id uint64) (io.ReadCloser, *model.File, error) {
    record, err := service.repo.FindFileByID(ctx, id)
    if err != nil {
        return nil, nil, err
    }
    driver, err := service.registry.Select(record.Driver)
    if err != nil {
        return nil, nil, storageError(err)
    }
    reader, err := driver.OpenReader(ctx, record.Path)
    return reader, record, err
}

func (service *Service) DeleteFile(ctx context.Context, id uint64) error {
    record, err := service.repo.FindFileByID(ctx, id)
    if err != nil {
        return err
    }
    driver, err := service.registry.Select(record.Driver)
    if err != nil {
        return storageError(err)
    }
    if err = service.repo.DeleteFileIfUnused(ctx, id); err != nil {
        return err
    }
    if err = driver.Delete(context.WithoutCancel(ctx), record.Path); err != nil {
        return errors.Join(storageError(err), service.repo.RestoreFile(context.WithoutCancel(ctx), id))
    }
    return nil
}

func (service *Service) CreateAttachment(ctx context.Context, params api.AttachmentCreateParams) (jsontype.SafeUint64, error) {
    if _, err := service.repo.FindFileByID(ctx, params.FileID.Uint64()); err != nil {
        return 0, err
    }
    attachment := &model.Attachment{FileID: params.FileID.Uint64(), Sort: params.Sort}
    if err := attachment.SetEntity(params.EntityType, params.EntityID.Uint64()); err != nil {
        return 0, err
    }
    if err := attachment.SetType(params.Type); err != nil {
        return 0, err
    }
    if err := service.repo.CreateAttachment(ctx, attachment); err != nil {
        return 0, err
    }
    return jsontype.SafeUint64(attachment.ID), nil
}

func (service *Service) UpdateAttachment(ctx context.Context, params api.AttachmentUpdateParams) error {
    if _, err := service.repo.FindAttachmentByID(ctx, params.ID.Uint64()); err != nil {
        return err
    }
    attachment := &model.Attachment{Base: modeltype.Base{ID: params.ID.Uint64()}, Sort: params.Sort}
    if err := attachment.SetType(params.Type); err != nil {
        return err
    }
    return service.repo.UpdateAttachment(ctx, attachment)
}

func (service *Service) DeleteAttachment(ctx context.Context, id uint64) error {
    if _, err := service.repo.FindAttachmentByID(ctx, id); err != nil {
        return err
    }
    return service.repo.DeleteAttachment(ctx, id)
}

func (service *Service) FindAttachmentByID(ctx context.Context, id uint64) (*api.Attachment, error) {
    record, err := service.repo.FindAttachmentByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return service.attachmentDTO(ctx, record)
}

func (service *Service) ListAttachments(ctx context.Context, params api.AttachmentListParams) (*apitype.PaginatedResult[*api.Attachment], error) {
    records, total, err := service.repo.ListAttachments(ctx, &params)
    if err != nil {
        return nil, err
    }
    list := make([]*api.Attachment, 0, len(records))
    for _, record := range records {
        item, dtoErr := service.attachmentDTO(ctx, record)
        if dtoErr != nil {
            return nil, dtoErr
        }
        list = append(list, item)
    }
    return &apitype.PaginatedResult[*api.Attachment]{Pagination: params.Pagination, Total: total, List: list}, nil
}

func (service *Service) fileDTO(ctx context.Context, record *model.File) (*api.File, error) {
    return &api.File{
        ID: jsontype.SafeUint64(record.ID), CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt,
        OriginalName: record.OriginalName, Ext: record.Ext, Size: jsontype.SafeUint64(record.Size),
        Driver: record.Driver, URL: record.Path,
    }, nil
}

func storageError(err error) error {
    return fmt.Errorf("%w: %v", errno.ErrStorageUnavailable, err)
}

func (service *Service) attachmentDTO(ctx context.Context, record *model.Attachment) (*api.Attachment, error) {
    result := &api.Attachment{
        ID: jsontype.SafeUint64(record.ID), CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt,
        FileID: jsontype.SafeUint64(record.FileID), EntityType: record.EntityType,
        EntityID: jsontype.SafeUint64(record.EntityID), Type: record.Type, Sort: record.Sort,
    }
    if record.File.ID != 0 {
        fileDTO, err := service.fileDTO(ctx, &record.File)
        if err != nil {
            return nil, err
        }
        result.File = fileDTO
    }
    return result, nil
}
