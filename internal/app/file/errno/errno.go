package errno

import "cms/internal/public/erroz"

var (
	ErrOriginalNameRequired  = erroz.NewError("FILE_001", "文件名必填")
	ErrOriginalNameTooLong   = erroz.NewError("FILE_002", "文件名长度超出限制")
	ErrExtensionTooLong      = erroz.NewError("FILE_003", "文件扩展名长度超出限制")
	ErrInvalidDriver         = erroz.NewError("FILE_004", "文件存储驱动不合法")
	ErrInvalidPath           = erroz.NewError("FILE_005", "文件存储路径不合法")
	ErrInvalidSHA256         = erroz.NewError("FILE_006", "文件摘要不合法")
	ErrInvalidEntityType     = erroz.NewError("FILE_007", "附件实体类型不合法")
	ErrEntityIDRequired      = erroz.NewError("FILE_008", "附件实体 ID 必填")
	ErrAttachmentTypeTooLong = erroz.NewError("FILE_009", "附件类型长度超出限制")
	ErrFileInUse             = erroz.NewError("FILE_010", "文件已被附件引用，不能删除")
	ErrStorageUnavailable    = erroz.NewError("FILE_011", "文件存储服务不可用")
)
