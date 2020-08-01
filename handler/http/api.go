package http

import (
	"fmt"
	"github.com/cloud-mesh/object-storage/model"
	"github.com/cloud-mesh/object-storage/model/usecase"
	"github.com/labstack/echo"
	"net/http"
)

const (
	maxUploadBytes = 5 * (1 << 30) // 华为云限制5G
)

type handler struct {
	ucase usecase.UseCase
}

func NewHandler(ucase usecase.UseCase) *handler {
	return &handler{ucase}
}

func (h *handler) Route(e *echo.Echo) {
	e.HEAD("/", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			// TODO: head region
			return c.NoContent(http.StatusNotFound)
		} else {
			return h.headBucket(c)
		}
	})
	e.GET("/", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			return h.listBucket(c)
		} else {
			return c.NoContent(http.StatusNotFound)
		}
	})
	e.PUT("/", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			return c.NoContent(http.StatusNotFound)
		} else {
			return h.createBucket(c)
		}
	})
	e.DELETE("/", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			return c.NoContent(http.StatusNotFound)
		} else {
			return h.deleteBucket(c)
		}
	})
	e.HEAD("/:object_key", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			return c.NoContent(http.StatusNotFound)
		}

		objectKey := c.Param("object_key")
		if objectKey == "" {
			return c.NoContent(http.StatusNotFound)
		}

		return h.headObject(c)
	})
	e.GET("/:object_key", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			return c.NoContent(http.StatusNotFound)
		}

		objectKey := c.Param("object_key")
		if objectKey == "" {
			return c.NoContent(http.StatusNotFound)
		}

		uploadId := c.QueryParam("upload_id")
		if uploadId == "" {
			return c.NoContent(http.StatusNotFound)
		}

		if _, ok := c.QueryParams()["parts"]; !ok {
			return c.NoContent(http.StatusNotFound)
		}

		return h.listParts(c)
	})
	e.POST("/:object_key", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			return c.NoContent(http.StatusNotFound)
		}

		objectKey := c.Param("object_key")
		if objectKey == "" {
			return c.NoContent(http.StatusNotFound)
		}

		if _, ok := c.QueryParams()["uploads"]; ok {
			return h.initMultipartUpload(c)
		}

		uploadId := c.QueryParam("upload_id")
		if uploadId == "" {
			return h.postObject(c)
		}

		if _, ok := c.QueryParams()["eof"]; ok {
			return h.completeMultipartUpload(c)
		}

		partNumber := c.QueryParam("part_number")
		if partNumber == "" {
			return c.NoContent(http.StatusNotFound)
		}

		return h.uploadPart(c)
	})
	e.PUT("/:object_key", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			return c.NoContent(http.StatusNotFound)
		}

		objectKey := c.Param("object_key")
		if objectKey == "" {
			return c.NoContent(http.StatusNotFound)
		}

		return h.putObject(c)
	})
	e.DELETE("/:object_key", func(c echo.Context) error {
		bucketName := getBucketName(c)
		if bucketName == "" {
			return c.NoContent(http.StatusNotFound)
		}

		objectKey := c.Param("object_key")
		if objectKey == "" {
			return c.NoContent(http.StatusNotFound)
		}

		uploadId := c.QueryParam("upload_id")
		if uploadId == "" {
			return h.deleteObject(c)
		} else {
			return h.abortMultipartUpload(c)
		}
	})
}

func (h *handler) headRegion(c echo.Context) error {
	region, err := h.ucase.HeadRegion()
	if err != nil {
		return err
	}
	return jsonOK(c, map[string]interface{}{
		"region": adapterRegion(region),
	})
}

func (h *handler) listBucket(c echo.Context) error {
	page, pageSize := getPaging(c)
	buckets, err := h.ucase.ListBucket(page, pageSize)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"buckets": adapterBuckets(buckets),
	})
}

func (h *handler) headBucket(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	bucket, err := h.ucase.GetBucket(bucketName)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"bucket": adapterBucket(bucket),
	})
}

func (h *handler) createBucket(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")

	post := &Bucket{}
	if err := c.Bind(post); err != nil {
		return model.ErrInvalidParam
	}

	if err := h.ucase.CreateBucket(bucketName); err != nil {
		return err
	}
	bucket, err := h.ucase.GetBucket(bucketName)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"bucket": adapterBucket(bucket),
	})
}

func (h *handler) deleteBucket(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")

	if err := h.ucase.DeleteBucket(bucketName); err != nil {
		return err
	}

	return jsonOK(c, nil)
}

func (h *handler) headObject(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")

	object, err := h.ucase.HeadObject(bucketName, objectKey)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"object": adapterObject(object),
	})
}

// Form表单文件当超过32MB时，会保存到文件系统，Unix上临时目录由环境变量`$TMPDIR`指定，如果`$TMPDIR`为空，则为`/tmp`。
// 生成的临时文件会在请求结束时自动删除。
// 需确保临时文件目录有足够的存储空间，否则上传将失败。
func (h *handler) postObject(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")

	f, fh, err := c.Request().FormFile("file")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	defer f.Close()

	if fh.Size <= 0 || fh.Size > maxUploadBytes {
		return c.String(http.StatusBadRequest, fmt.Sprintf("文件大小为%d", fh.Size))
	}

	if err := h.ucase.PutObject(bucketName, objectKey, f); err != nil {
		return err
	}

	object, err := h.ucase.HeadObject(bucketName, objectKey)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"object": adapterObject(object),
	})
}

// 流式上传，不消耗内存及磁盘，即服务器从请求流读取数据，并将读取到的数据直接上传到对象存储。
// 相比表单上传方式，流式上传的失败概率会增加。因为只要上传中，读取连接或上传连接失败，即上传失败。
func (h *handler) putObject(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	body := c.Request().Body
	defer body.Close()

	if err := h.ucase.PutObject(bucketName, objectKey, body); err != nil {
		return err
	}

	object, err := h.ucase.HeadObject(bucketName, objectKey)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"object": adapterObject(object),
	})
}

func (h *handler) deleteObject(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")

	if err := h.ucase.DeleteObject(bucketName, objectKey); err != nil {
		return err
	}

	return jsonOK(c, nil)
}

func (h *handler) initMultipartUpload(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")

	uploadID, err := h.ucase.InitMultipartUpload(bucketName, objectKey)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"upload_id": uploadID,
	})
}

func (h *handler) completeMultipartUpload(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	uploadID := c.QueryParam("upload_id")

	post := struct {
		CompleteParts []model.ObjectCompletePart `json:"complete_parts"`
	}{}
	if err := c.Bind(&post); err != nil {
		return model.ErrInvalidParam
	}

	if err := h.ucase.CompleteUploadPart(bucketName, objectKey, uploadID, post.CompleteParts); err != nil {
		return err
	}

	return jsonOK(c, nil)
}

func (h *handler) abortMultipartUpload(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	uploadID := c.QueryParam("upload_id")

	if err := h.ucase.AbortMultipartUpload(bucketName, objectKey, uploadID); err != nil {
		return err
	}

	return jsonOK(c, nil)
}

func (h *handler) uploadPart(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	uploadID := c.QueryParam("upload_id")
	partNumber := getInt(c.QueryParam("part_id"), 0)

	f, fh, err := c.Request().FormFile("file")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	defer f.Close()

	if fh.Size <= 0 || fh.Size > maxUploadBytes {
		return c.String(http.StatusBadRequest, fmt.Sprintf("文件大小为%d", fh.Size))
	}

	etag, err := h.ucase.UploadPart(bucketName, objectKey, uploadID, partNumber, f)
	if err != nil {
		return err
	}

	part := &model.ObjectPart{
		PartNumber: partNumber,
		ETag:       etag,
	}
	return jsonOK(c, map[string]interface{}{
		"part": adapterObjectPart(part),
	})
}

func (h *handler) listParts(c echo.Context) error {
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	uploadID := c.QueryParam("upload_id")

	parts, err := h.ucase.ListParts(bucketName, objectKey, uploadID)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"parts": adapterObjectParts(parts),
	})
}
