package file

import (
	"errors"
	"mime"
	"net/http"
	"strconv"

	"cms/internal/app/file/api"
	"cms/internal/httpx"
	"cms/internal/middleware/authz"
	"cms/internal/public/apitype"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
	auth    *authz.Factory
}

func NewFileController(service *Service, auth *authz.Factory) *Controller {
	return &Controller{service: service, auth: auth}
}

func (controller *Controller) Setup(router httpx.Router) {
	auth := controller.auth.Resource("file")
	group := router.Group("/file")
	group.POST("/upload", auth.Permission("upload").Wrap(httpx.HandlerFunc(controller.Upload)))
	group.POST("/list", auth.Permission("read").Wrap(httpx.HandlerFunc(controller.List)))
	group.POST("/detail", auth.Permission("read").Wrap(httpx.HandlerFunc(controller.Detail)))
	group.POST("/delete", auth.Permission("delete").Wrap(httpx.HandlerFunc(controller.Delete)))
	group.GET("/download", auth.Permission("read").Wrap(httpx.HandlerFunc(controller.Download)))
	group.POST("/attachment/create", auth.Permission("update").Wrap(httpx.HandlerFunc(controller.CreateAttachment)))
	group.POST("/attachment/update", auth.Permission("update").Wrap(httpx.HandlerFunc(controller.UpdateAttachment)))
	group.POST("/attachment/delete", auth.Permission("update").Wrap(httpx.HandlerFunc(controller.DeleteAttachment)))
	group.POST("/attachment/detail", auth.Permission("read").Wrap(httpx.HandlerFunc(controller.AttachmentDetail)))
	group.POST("/attachment/list", auth.Permission("read").Wrap(httpx.HandlerFunc(controller.AttachmentList)))
}

// Upload 上传文件
// @x-apifox-folder "文件"
// @Security ApiKeyAuth
// @Summary 上传文件
// @Tags 文件
// @Accept mpfd
// @Produce json
// @Param file formData file true "文件"
// @Param target formData string false "存储目标；为空时使用 default"
// @Success 200 {object} api.ApiFile
// @Router /file/upload [post]
func (controller *Controller) Upload(ctx *gin.Context) error {
	header, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	source, err := header.Open()
	if err != nil {
		return err
	}
	defer source.Close()
	target := ctx.PostForm("target")
	var result *api.File
	if target == "" {
		result, err = controller.service.Upload(ctx, header.Filename, source)
	} else {
		result, err = controller.service.UploadTo(ctx, target, header.Filename, source)
	}
	if err != nil {
		return err
	}
	httpx.JSON(ctx, httpx.OK.WithData(result))
	return nil
}

// List 文件列表
// @x-apifox-folder "文件"
// @Security ApiKeyAuth
// @Summary 文件列表
// @Tags 文件
// @Accept json
// @Produce json
// @Param body body api.FileListParams true "请求参数"
// @Success 200 {object} api.ApiFileList
// @Router /file/list [post]
func (controller *Controller) List(ctx *gin.Context) error {
	return httpx.BindJSON[api.FileListParams](ctx, func(params api.FileListParams) (any, error) {
		return controller.service.ListFiles(ctx, params)
	})
}

// Detail 文件详情
// @x-apifox-folder "文件"
// @Security ApiKeyAuth
// @Summary 文件详情
// @Tags 文件
// @Accept json
// @Produce json
// @Param body body apitype.ResourceID true "请求参数"
// @Success 200 {object} api.ApiFile
// @Router /file/detail [post]
func (controller *Controller) Detail(ctx *gin.Context) error {
	return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
		return controller.service.FindFileByID(ctx, params.ID.Uint64())
	})
}

// Delete 删除文件
// @x-apifox-folder "文件"
// @Security ApiKeyAuth
// @Summary 删除文件
// @Tags 文件
// @Accept json
// @Produce json
// @Param body body apitype.ResourceID true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /file/delete [post]
func (controller *Controller) Delete(ctx *gin.Context) error {
	return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
		return nil, controller.service.DeleteFile(ctx, params.ID.Uint64())
	})
}

// Download 下载文件
// @x-apifox-folder "文件"
// @Security ApiKeyAuth
// @Summary 下载文件
// @Tags 文件
// @Produce application/octet-stream
// @Param id query string true "文件 ID"
// @Success 200 {file} binary
// @Router /file/download [get]
func (controller *Controller) Download(ctx *gin.Context) error {
	id, err := strconv.ParseUint(ctx.Query("id"), 10, 64)
	if err != nil || id == 0 {
		return errors.New("invalid file id")
	}
	reader, record, err := controller.service.Open(ctx, id)
	if err != nil {
		return err
	}
	defer reader.Close()
	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": record.OriginalName})
	ctx.DataFromReader(http.StatusOK, int64(record.Size), "application/octet-stream", reader, map[string]string{
		"Content-Disposition":    disposition,
		"X-Content-Type-Options": "nosniff",
	})
	return nil
}

// CreateAttachment 创建附件绑定
// @x-apifox-folder "文件/附件"
// @Security ApiKeyAuth
// @Summary 创建附件绑定
// @Tags 附件
// @Accept json
// @Produce json
// @Param body body api.AttachmentCreateParams true "请求参数"
// @Success 200 {object} apitype.ApiCreateResult
// @Router /file/attachment/create [post]
func (controller *Controller) CreateAttachment(ctx *gin.Context) error {
	return httpx.BindJSON[api.AttachmentCreateParams](ctx, func(params api.AttachmentCreateParams) (any, error) {
		return controller.service.CreateAttachment(ctx, params)
	})
}

// UpdateAttachment 更新附件
// @x-apifox-folder "文件/附件"
// @Security ApiKeyAuth
// @Summary 更新附件
// @Tags 附件
// @Accept json
// @Produce json
// @Param body body api.AttachmentUpdateParams true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /file/attachment/update [post]
func (controller *Controller) UpdateAttachment(ctx *gin.Context) error {
	return httpx.BindJSON[api.AttachmentUpdateParams](ctx, func(params api.AttachmentUpdateParams) (any, error) {
		return nil, controller.service.UpdateAttachment(ctx, params)
	})
}

// DeleteAttachment 删除附件绑定
// @x-apifox-folder "文件/附件"
// @Security ApiKeyAuth
// @Summary 删除附件绑定
// @Tags 附件
// @Accept json
// @Produce json
// @Param body body apitype.ResourceID true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /file/attachment/delete [post]
func (controller *Controller) DeleteAttachment(ctx *gin.Context) error {
	return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
		return nil, controller.service.DeleteAttachment(ctx, params.ID.Uint64())
	})
}

// AttachmentDetail 附件详情
// @x-apifox-folder "文件/附件"
// @Security ApiKeyAuth
// @Summary 附件详情
// @Tags 附件
// @Accept json
// @Produce json
// @Param body body apitype.ResourceID true "请求参数"
// @Success 200 {object} api.ApiAttachment
// @Router /file/attachment/detail [post]
func (controller *Controller) AttachmentDetail(ctx *gin.Context) error {
	return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
		return controller.service.FindAttachmentByID(ctx, params.ID.Uint64())
	})
}

// AttachmentList 附件列表
// @x-apifox-folder "文件/附件"
// @Security ApiKeyAuth
// @Summary 附件列表
// @Tags 附件
// @Accept json
// @Produce json
// @Param body body api.AttachmentListParams true "请求参数"
// @Success 200 {object} api.ApiAttachmentList
// @Router /file/attachment/list [post]
func (controller *Controller) AttachmentList(ctx *gin.Context) error {
	return httpx.BindJSON[api.AttachmentListParams](ctx, func(params api.AttachmentListParams) (any, error) {
		return controller.service.ListAttachments(ctx, params)
	})
}
