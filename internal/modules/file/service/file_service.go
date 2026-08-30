package service

import (
	"context"
	"crypto"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"path"
	"time"

	"cms/internal/public/jsontype"
	"cms/internal/public/model"

	"cms/internal/infra/logger"
	contract2 "cms/internal/modules/file/contract"
	"cms/internal/modules/file/internal/assembler"
	"cms/internal/modules/file/internal/dto"
	"cms/internal/modules/file/internal/errno"
	"cms/internal/modules/file/internal/fileutil"

	"github.com/avast/retry-go"
)

type FileService struct {
	txManager persistence.Transactor
	repo      contract2.FileRepo
	logger    *logger.Logger
	driverSrv *FileDriverService
}

func NewFileService(txManager persistence.Transactor, repo contract2.FileRepo, driverSrv *FileDriverService, logger *logger.Logger) *FileService {
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
	var savingKey string
	hash := crypto.SHA256.New()
	counted := &countingReader{Reader: io.TeeReader(params.Data, hash)}

	err = retry.Do(
		func() error {
			genPath, retryErr := fileutil.GeneratePath()
			if retryErr != nil {
				return retryErr
			}

			savingKey = genPath.FullPath() + fileExt
			_, statErr := driver.Stat(ctx, savingKey)
			if statErr == nil {
				return fs.ErrExist
			}
			if !errors.Is(statErr, fs.ErrNotExist) {
				return retry.Unrecoverable(statErr)
			}
			retryErr := driver.Create(ctx, savingKey, counted, -1)
			if retryErr != nil {
				if errors.Is(retryErr, fs.ErrExist) {
					return retryErr
				}
				return retry.Unrecoverable(retryErr)
			}

			return nil
		},
		retry.Context(ctx),
		retry.Attempts(3),
	)

	if err != nil {
		return nil, errno.FilePreCreateFileFailed.Wrap(err).ToError()
	}

	fileSize := counted.n

	delCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	fileRecord := &model.File{
		UserID:       user.ID,
		OriginalName: params.Name,
		Ext:          fileExt,
		Path:         savingKey,
		Size:         jsontype.SafeInt64(fileSize),
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

	return assembler.ToFileInfoDTO(fileRecord, savingKey), nil
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

func (s FileService) Delete(ctx context.Context, id jsontype.SafeUint64) error {
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

func (s FileService) OpenFileRecordByID(ctx context.Context, id jsontype.SafeUint64) (io.ReadCloser, *dto.FileMeta, error) {
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
