package file

import (
	"context"

	"cms/internal/app/file/api"
	"cms/internal/app/file/errno"
	"cms/internal/app/file/model"
	"cms/internal/infra/store/gorm/dbscope"
	"cms/internal/infra/store/gorm/gquery"

	"gorm.io/gorm/clause"
)

type Repo struct {
	q *gquery.Query
}

func NewFileRepo(q *gquery.Query) *Repo { return &Repo{q: q} }

func (repo *Repo) CreateFile(ctx context.Context, data *model.File) error {
	return repo.q.File.WithContext(ctx).Create(data)
}

func (repo *Repo) DeleteFileIfUnused(ctx context.Context, id uint64) error {
	return repo.q.Transaction(func(tx *gquery.Query) error {
		fileQuery := tx.File
		if _, err := fileQuery.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(fileQuery.ID.Eq(id)).
			First(); err != nil {
			return err
		}
		attachmentQuery := tx.Attachment
		count, err := attachmentQuery.WithContext(ctx).Where(attachmentQuery.FileID.Eq(id)).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return errno.ErrFileInUse
		}
		_, err = fileQuery.WithContext(ctx).Where(fileQuery.ID.Eq(id)).Delete()
		return err
	})
}

func (repo *Repo) RestoreFile(ctx context.Context, id uint64) error {
	query := repo.q.File
	_, err := query.WithContext(ctx).Unscoped().Where(query.ID.Eq(id)).Update(query.DeletedAt, nil)
	return err
}

func (repo *Repo) FindFileByID(ctx context.Context, id uint64) (*model.File, error) {
	query := repo.q.File
	return query.WithContext(ctx).Where(query.ID.Eq(id)).First()
}

func (repo *Repo) ListFiles(ctx context.Context, params *api.FileListParams) ([]*model.File, int64, error) {
	fileQuery := repo.q.File
	query := fileQuery.WithContext(ctx)
	if params.OriginalName != nil {
		query = query.Where(fileQuery.OriginalName.Like("%" + *params.OriginalName + "%"))
	}
	if params.Ext != nil {
		query = query.Where(fileQuery.Ext.Eq(*params.Ext))
	}
	if params.Driver != nil {
		query = query.Where(fileQuery.Driver.Eq(*params.Driver))
	}
	total, err := query.Count()
	if err != nil {
		return nil, 0, err
	}
	list, err := query.Scopes(dbscope.Paginate(params.PageNo, params.PageSize)).Order(fileQuery.ID.Desc()).Find()
	return list, total, err
}

func (repo *Repo) CreateAttachment(ctx context.Context, data *model.Attachment) error {
	return repo.q.Transaction(func(tx *gquery.Query) error {
		fileQuery := tx.File
		if _, err := fileQuery.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(fileQuery.ID.Eq(data.FileID)).
			First(); err != nil {
			return err
		}
		return tx.Attachment.WithContext(ctx).Create(data)
	})
}

func (repo *Repo) UpdateAttachment(ctx context.Context, data *model.Attachment) error {
	query := repo.q.Attachment
	_, err := query.WithContext(ctx).
		Where(query.ID.Eq(data.ID)).
		Select(query.Type, query.Sort).
		Updates(data)
	return err
}

func (repo *Repo) DeleteAttachment(ctx context.Context, id uint64) error {
	query := repo.q.Attachment
	_, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Delete()
	return err
}

func (repo *Repo) FindAttachmentByID(ctx context.Context, id uint64) (*model.Attachment, error) {
	query := repo.q.Attachment
	return query.WithContext(ctx).Preload(query.File).Where(query.ID.Eq(id)).First()
}

func (repo *Repo) ListAttachments(ctx context.Context, params *api.AttachmentListParams) ([]*model.Attachment, int64, error) {
	attachmentQuery := repo.q.Attachment
	query := attachmentQuery.WithContext(ctx)
	if params.FileID != nil {
		query = query.Where(attachmentQuery.FileID.Eq(params.FileID.Uint64()))
	}
	if params.EntityType != nil {
		query = query.Where(attachmentQuery.EntityType.Eq(*params.EntityType))
	}
	if params.EntityID != nil {
		query = query.Where(attachmentQuery.EntityID.Eq(params.EntityID.Uint64()))
	}
	if params.Type != nil {
		query = query.Where(attachmentQuery.Type.Eq(*params.Type))
	}
	total, err := query.Count()
	if err != nil {
		return nil, 0, err
	}
	list, err := query.
		Scopes(dbscope.Paginate(params.PageNo, params.PageSize)).
		Preload(attachmentQuery.File).
		Order(attachmentQuery.Sort, attachmentQuery.ID).
		Find()
	return list, total, err
}
