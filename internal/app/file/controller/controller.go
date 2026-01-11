package controller

import (
    `dpcms/internal/app/file/internal/dto`
    `dpcms/internal/app/file/service`
    `dpcms/internal/erroz`
    `dpcms/internal/util/contextutil`
    
    `github.com/gin-gonic/gin`
)

type FileController struct {
    fileSrv *service.FileService
}

func NewFileController(fileSrv *service.FileService) *FileController {
    return &FileController{fileSrv: fileSrv}
}

func (s FileController) Setup(engine *gin.Engine) {
    g := engine.Group("/file")
    g.POST("/upload", s.Upload)
}

// Upload
// @x-apifox-folder "文件"
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
    u, err := contextutil.GetAuthorizedUser(ctx)
    if err != nil {
        erroz.ResolveWithAbort(ctx, err)
        return
    }
    result, err := s.fileSrv.Save(ctx, u, dto.FileSaveCommand{
        Filename: fileHeader.Filename,
        Data:     f,
    })
    if err != nil {
        erroz.ResolveWithAbort(ctx, err)
        return
    }
    erroz.OK.WithOption(erroz.WithData(result)).Write(ctx)
}
