package web

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// UEditorWriter ueditor's file writer
type UEditorWriter func(c *gin.Context, name string, buf []byte, size int64) (url string, err error)

// UEditorManager ueditor's file manager
type UEditorManager func(c *gin.Context) (urls []string, err error)

// UEditor ueditor
type UEditor struct {
}

// Upload upload handler
func (p *UEditor) Upload(wrt UEditorWriter, images UEditorManager, files UEditorManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Query("action") {
		case "config":
			p.config(c)
		case "uploadimage":
			p.upload(c, "upfile", wrt)
		case "uploadscrawl":
			p.scrawl(c, wrt)
		case "uploadvideo":
			p.upload(c, "upfile", wrt)
		case "uploadfile":
			p.upload(c, "upfile", wrt)
		case "listimage":
			p.manager(c, images)
		case "listfile":
			p.manager(c, files)
		case "catchimage":
			p.catchImage(c)
		default:
			log.Error("TODO: ueditor action")
		}
	}
}

func (p *UEditor) scrawl(c *gin.Context, f UEditorWriter) {
	if err := c.Request.ParseForm(); err != nil {
		p.fail(c, err)
	}
	buf, err := base64.StdEncoding.DecodeString(c.Request.FormValue("upfile"))
	if err != nil {
		p.fail(c, err)
		return
	}

	name := "scrawl-" + time.Now().Format(time.RFC822) + ".png"
	size := int64(len(buf))
	if size == 0 {
		p.fail(c, errors.New("empty file"))
		return
	}
	url, err := f(c, name, buf, size)
	if err != nil {
		p.fail(c, err)
		return
	}
	p.file(c, name, url)
}

func (p *UEditor) file(c *gin.Context, n, u string) {
	p.success(c, gin.H{
		"url":      u,
		"title":    "",
		"original": n,
	})
}

func (p *UEditor) success(c *gin.Context, d gin.H) {
	d["state"] = "SUCCESS"
	c.JSON(http.StatusOK, d)
}

func (p *UEditor) fail(c *gin.Context, e error) {
	log.Error(e)
	c.JSON(http.StatusOK, gin.H{
		"state": "FAILED",
	})
}

func (p *UEditor) upload(c *gin.Context, name string, fn UEditorWriter) {
	fd, fh, err := c.Request.FormFile(name)
	if err != nil {
		p.fail(c, err)
		return
	}
	defer fd.Close()
	buf := make([]byte, fh.Size)
	if _, err = fd.Read(buf); err != nil {
		p.fail(c, err)
		return
	}
	url, err := fn(c, fh.Filename, buf, int64(len(buf)))
	if err != nil {
		p.fail(c, err)
		return
	}
	p.file(c, name, url)
}

func (p *UEditor) manager(c *gin.Context, f UEditorManager) {
	items, err := f(c)
	if err != nil {
		p.fail(c, err)
		return
	}
	var list []gin.H

	for _, it := range items {
		list = append(list, gin.H{"url": it})
	}
	p.success(c, gin.H{
		"list":  list,
		"start": 0,
		"total": len(list),
	})

}

func (p *UEditor) catchImage(c *gin.Context) {
	// TODO
}

func (p *UEditor) config(c *gin.Context) {
	/* 前后端通信相关的配置,注释只允许使用多行方式 */
	cfg := gin.H{
		/* 上传图片配置项 */
		"imageActionName":     "uploadimage",                                     /* 执行上传图片的action名称 */
		"imageFieldName":      "upfile",                                          /* 提交的图片表单名称 */
		"imageMaxSize":        2048000,                                           /* 上传大小限制，单位B */
		"imageAllowFiles":     []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 上传图片格式显示 */
		"imageCompressEnable": true,                                              /* 是否压缩图片,默认是true */
		"imageCompressBorder": 1600,                                              /* 图片压缩最长边限制 */
		"imageInsertAlign":    "none",                                            /* 插入的图片浮动方式 */
		"imageUrlPrefix":      "",                                                /* 图片访问路径前缀 */
		"imagePathFormat":     "/images/{yyyy}{mm}{dd}/{time}{rand:6}",           /* 上传保存路径,可以自定义保存路径和文件名格式 */
		/* {filename} 会替换成原文件名,配置这项需要注意中文乱码问题 */
		/* {rand:6} 会替换成随机数,后面的数字是随机数的位数 */
		/* {time} 会替换成时间戳 */
		/* {yyyy} 会替换成四位年份 */
		/* {yy} 会替换成两位年份 */
		/* {mm} 会替换成两位月份 */
		/* {dd} 会替换成两位日期 */
		/* {hh} 会替换成两位小时 */
		/* {ii} 会替换成两位分钟 */
		/* {ss} 会替换成两位秒 */
		/* 非法字符 \ : * ? " < > | */
		/* 具请体看线上文档: fex.baidu.com/ueditor/#use-format_upload_filename */
		/* 涂鸦图片上传配置项 */
		"scrawlActionName":  "uploadscrawl",                          /* 执行上传涂鸦的action名称 */
		"scrawlFieldName":   "upfile",                                /* 提交的图片表单名称 */
		"scrawlPathFormat":  "/images/{yyyy}{mm}{dd}/{time}{rand:6}", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"scrawlMaxSize":     2048000,                                 /* 上传大小限制，单位B */
		"scrawlUrlPrefix":   "",                                      /* 图片访问路径前缀 */
		"scrawlInsertAlign": "none",
		/* 截图工具上传 */
		"snapscreenActionName":  "uploadimage",                           /* 执行上传截图的action名称 */
		"snapscreenPathFormat":  "/images/{yyyy}{mm}{dd}/{time}{rand:6}", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"snapscreenUrlPrefix":   "",                                      /* 图片访问路径前缀 */
		"snapscreenInsertAlign": "none",                                  /* 插入的图片浮动方式 */
		/* 抓取远程图片配置 */
		"catcherLocalDomain": []string{"127.0.0.1", "localhost", "image.baidu.com"},
		"catcherActionName":  "catchimage",                                      /* 执行抓取远程图片的action名称 */
		"catcherFieldName":   "source",                                          /* 提交的图片列表表单名称 */
		"catcherPathFormat":  "/images/{yyyy}{mm}{dd}/{time}{rand:6}",           /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"catcherUrlPrefix":   "",                                                /* 图片访问路径前缀 */
		"catcherMaxSize":     2048000,                                           /* 上传大小限制，单位B */
		"catcherAllowFiles":  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 抓取图片格式显示 */
		/* 上传视频配置 */
		"videoActionName": "uploadvideo",                           /* 执行上传视频的action名称 */
		"videoFieldName":  "upfile",                                /* 提交的视频表单名称 */
		"videoPathFormat": "/videos/{yyyy}{mm}{dd}/{time}{rand:6}", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"videoUrlPrefix":  "",                                      /* 视频访问路径前缀 */
		"videoMaxSize":    102400000,                               /* 上传大小限制，单位B，默认100MB */
		"videoAllowFiles": []string{
			".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
			".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
		}, /* 上传视频格式显示 */
		/* 上传文件配置 */
		"fileActionName": "uploadfile",                           /* controller里,执行上传视频的action名称 */
		"fileFieldName":  "upfile",                               /* 提交的文件表单名称 */
		"filePathFormat": "/files/{yyyy}{mm}{dd}/{time}{rand:6}", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"fileUrlPrefix":  "",                                     /* 文件访问路径前缀 */
		"fileMaxSize":    51200000,                               /* 上传大小限制，单位B，默认50MB */
		"fileAllowFiles": []string{
			".png", ".jpg", ".jpeg", ".gif", ".bmp",
			".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
			".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
			".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
			".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
		}, /* 上传文件格式显示 */
		/* 列出指定目录下的图片 */
		"imageManagerActionName":  "listimage",                                       /* 执行图片管理的action名称 */
		"imageManagerListPath":    "/images/",                                        /* 指定要列出图片的目录 */
		"imageManagerListSize":    20,                                                /* 每次列出文件数量 */
		"imageManagerUrlPrefix":   "",                                                /* 图片访问路径前缀 */
		"imageManagerInsertAlign": "none",                                            /* 插入的图片浮动方式 */
		"imageManagerAllowFiles":  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 列出的文件类型 */
		/* 列出指定目录下的文件 */
		"fileManagerActionName": "listfile", /* 执行文件管理的action名称 */
		"fileManagerListPath":   "/files/",  /* 指定要列出文件的目录 */
		"fileManagerUrlPrefix":  "",         /* 文件访问路径前缀 */
		"fileManagerListSize":   20,         /* 每次列出文件数量 */
		"fileManagerAllowFiles": []string{
			".png", ".jpg", ".jpeg", ".gif", ".bmp",
			".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
			".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
			".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
			".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
		}, /* 列出的文件类型 */
	}

	c.JSON(http.StatusOK, cfg)
}
