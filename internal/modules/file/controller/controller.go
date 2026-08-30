package controller

import (
    "fmt"
    
    `cms/internal/pkg/hashid`
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
    
    "cms/internal/erroz"
    "cms/internal/httpx"
    erroz2 "cms/internal/httpx/erroz"
    "cms/internal/middleware/authz"
    "cms/internal/modules/file/internal/dto"
    "cms/internal/modules/file/service"
    "cms/internal/util/authzutil"
    
    "github.com/gin-gonic/gin"
)

type FileController struct {
    fileSrv *service.FileService
    auth    *authz.Factory
    hashID  *hashid.HashID
}

func NewFileController(fileSrv *service.FileService, auth *authz.Factory, hashID *hashid.HashID) *FileController {
    return &FileController{fileSrv: fileSrv, auth: auth, hashID: hashID}
}

func (s FileController) Setup(router httpx.Router) {
    acl := s.auth.auth("file")
    
    g := router.Group("/file", acl.Middleware())
    g.POST("/upload", s.Upload)
    g.POST("/delete", s.Delete)
    g.GET("/download", s.Download)
    
    acl.WithRouterOption(
        g,
        authz.WithRouterPermission("/upload", "upload"),
        authz.WithRouterPermission("/delete", "delete"),
        authz.WithRouterWhitelist("/download"),
    )
}

// Upload
// @x-apifox-folder "文件"
// @Security ApiKeyAuth
// @Summary 文件上传
// @Tags 文件
// @Accept mpfd
// @Produce json
// @Param file formData file true "请求参数"
// @Success 200 {object} erroz.Result{data=dto.FileInfo}
// @Router /file/upload [post]
func (s FileController) Upload(ctx *gin.Context) {
    fileHeader, err := ctx.FormFile("file")
    if err != nil {
        erroz.ResolveWithAbort(ctx, err)
        return
    }
    f, err := fileHeader.Open()
    if err != nil {
        erroz.ResolveWithAbort(ctx, err)
        return
    }
    u, err := authzutil.GetAuthorizedUser(ctx)
    if err != nil {
        erroz.ResolveWithAbort(ctx, err)
        return
    }
    result, err := s.fileSrv.Save(ctx, u, dto.FileSaveCommand{
        Name: fileHeader.Filename,
        Data: f,
    })
    if err != nil {
        erroz.ResolveWithAbort(ctx, err)
        return
    }
    httpx.OK.WithOption(erroz2.WithData(result)).Write(ctx)
}

// Download
// @x-apifox-folder "文件"
// @Security ApiKeyAuth
// @Summary 文件下载
// @Tags 文件
// @Accept json
// @Produce application/octet-stream
// @Param id gquery string true "文件ID"
// @Success 200 {file} file
// @Router /file/download [get]
func (s FileController) Download(ctx *gin.Context) {
    hashID := ctx.Query("id")
    if hashID == "" {
        httpx.DataNotFound.WriteWithAbort(ctx)
        return
    }
    id, err := s.hashID.Decode(hashID)
    if err != nil {
        httpx.DataNotFound.WriteWithAbort(ctx)
        return
    }
    reader, meta, err := s.fileSrv.OpenFileRecordByID(ctx, jsontype.SafeUint64(id[0]))
    if err != nil {
        erroz.ResolveWithAbort(ctx, err)
        return
    }
    ctx.DataFromReader(
        200,
        meta.Size.Raw(),
        "application/octet-stream",
        reader,
        map[string]string{
            "Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, meta.Name),
        },
    )
    _ = reader.Close()
}

// Delete
// @x-apifox-folder "文件"
// @Security ApiKeyAuth
// @Summary 文件删除
// @Tags 文件
// @Accept json
// @Produce json
// @Param body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router /file/delete [post]
func (s FileController) Delete(ctx *gin.Context) {
    httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return nil, s.fileSrv.Delete(ctx, params.ID)
    })
}
