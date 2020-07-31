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
	bucketMux := NewQueryMux()
	e.Any("/", bucketMux.Handle())
	// 桶列表
	bucketMux.GET("?vendor=:vendor_name", h.listBucket)
	// 创建桶
	bucketMux.POST("?vendor=:vendor_name", h.createBucket)
	// 桶详情
	bucketMux.GET("?vendor=:vendor_name&bucket=:bucket_name", h.getBucket)
	// 删除桶
	bucketMux.DELETE("?vendor=:vendor_name&bucket=:bucket_name", h.deleteBucket)

	objectMux := NewQueryMux()
	e.Any("/:object_key", objectMux.Handle())
	// 对象详情
	objectMux.HEAD("?vendor=:vendor_name&bucket=:bucket_name", h.headObject)
	// 删除对象
	objectMux.DELETE("?vendor=:vendor_name&bucket=:bucket_name", h.deleteObject)
	// 表单上传
	objectMux.POST("?vendor=:vendor_name&bucket=:bucket_name", h.postObject)
	// 流式上传
	objectMux.PUT("?vendor=:vendor_name&bucket=:bucket_name", h.putObject)
	// 初始化分段上传
	objectMux.POST("?vendor=:vendor_name&bucket=:bucket_name&uploads", h.initMultipartUpload)
	// 完成分段上传
	objectMux.POST("?vendor=:vendor_name&bucket=:bucket_name&upload_id=:upload_id&eof", h.completeMultipartUpload)
	// 取消分段上传
	objectMux.DELETE("?vendor=:vendor_name&bucket=:bucket_name&upload_id=:upload_id", h.abortMultipartUpload)
	// 上传分段
	objectMux.POST("?vendor=:vendor_name&bucket=:bucket_name&upload_id=:upload_id&part_id=:part_id", h.uploadPart)
	// 查看分段列表
	objectMux.GET("?vendor=:vendor_name&bucket=:bucket_name&upload_id=:upload_id&parts", h.listParts)
}

func (h *handler) listBucket(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	page, pageSize := getPaging(c)

	buckets, err := h.ucase.ListBucket(vendor, page, pageSize)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"buckets": adapterBuckets(buckets),
	})
}

func (h *handler) getBucket(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")

	bucket, err := h.ucase.GetBucket(vendor, bucketName)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"bucket": adapterBucket(bucket),
	})
}

func (h *handler) createBucket(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")

	post := &Bucket{}
	if err := c.Bind(post); err != nil {
		return model.ErrInvalidParam
	}

	if err := h.ucase.CreateBucket(vendor, bucketName); err != nil {
		return err
	}
	bucket, err := h.ucase.GetBucket(vendor, bucketName)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"bucket": adapterBucket(bucket),
	})
}

func (h *handler) deleteBucket(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")

	if err := h.ucase.DeleteBucket(vendor, bucketName); err != nil {
		return err
	}

	return jsonOK(c, nil)
}

func (h *handler) headObject(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")

	object, err := h.ucase.HeadObject(vendor, bucketName, objectKey)
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
	vendor := c.QueryParam("vendor_name")
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

	if err := h.ucase.PutObject(vendor, bucketName, objectKey, f); err != nil {
		return err
	}

	object, err := h.ucase.HeadObject(vendor, bucketName, objectKey)
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
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	body := c.Request().Body
	defer body.Close()

	if err := h.ucase.PutObject(vendor, bucketName, objectKey, body); err != nil {
		return err
	}

	object, err := h.ucase.HeadObject(vendor, bucketName, objectKey)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"object": adapterObject(object),
	})
}

func (h *handler) deleteObject(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")

	if err := h.ucase.DeleteObject(vendor, bucketName, objectKey); err != nil {
		return err
	}

	return jsonOK(c, nil)
}

func (h *handler) initMultipartUpload(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")

	uploadID, err := h.ucase.InitMultipartUpload(vendor, bucketName, objectKey)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"upload_id": uploadID,
	})
}

func (h *handler) completeMultipartUpload(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	uploadID := c.QueryParam("upload_id")

	post := struct {
		CompleteParts []model.ObjectCompletePart `json:"complete_parts"`
	}{}
	if err := c.Bind(&post); err != nil {
		return model.ErrInvalidParam
	}

	if err := h.ucase.CompleteUploadPart(vendor, bucketName, objectKey, uploadID, post.CompleteParts); err != nil {
		return err
	}

	return jsonOK(c, nil)
}

func (h *handler) abortMultipartUpload(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	uploadID := c.QueryParam("upload_id")

	if err := h.ucase.AbortMultipartUpload(vendor, bucketName, objectKey, uploadID); err != nil {
		return err
	}

	return jsonOK(c, nil)
}

func (h *handler) uploadPart(c echo.Context) error {
	vendor := c.QueryParam("vendor_name")
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

	etag, err := h.ucase.UploadPart(vendor, bucketName, objectKey, uploadID, partNumber, f)
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
	vendor := c.QueryParam("vendor_name")
	bucketName := c.QueryParam("bucket_name")
	objectKey := c.Param("object_key")
	uploadID := c.QueryParam("upload_id")

	parts, err := h.ucase.ListParts(vendor, bucketName, objectKey, uploadID)
	if err != nil {
		return err
	}

	return jsonOK(c, map[string]interface{}{
		"parts": adapterObjectParts(parts),
	})
}
