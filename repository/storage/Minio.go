package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/ysfgrl/gerror"
	"log"
	"mime/multipart"
	"net/url"
	"strings"
	"time"
)

var chars = []string{
	" ",
	"ş",
	"Ş",
	"ç",
	"Ç",
	"ö",
	"Ö",
	"ü",
	"Ü",
	"ğ",
	"Ğ",
	"?",
	",",
	"^",
	"'",
	"!",
}

type Base struct {
	Client *minio.Client
	Bucket string
	Exist  bool
}

func (b *Base) exist(ctx context.Context) *gerror.Error {

	if !b.Client.IsOnline() {
		return &gerror.Error{
			Code: "storage_not_online",
		}

	}
	var err error
	b.Exist, err = b.Client.BucketExists(ctx, b.Bucket)
	if err != nil {
		return gerror.GetError(err)
	}
	if !b.Exist {
		return b.create(ctx)
	}
	return nil
}

func (b *Base) create(ctx context.Context) *gerror.Error {
	err := b.Client.MakeBucket(ctx, b.Bucket, minio.MakeBucketOptions{
		Region: "us-east-1",
	})
	if err != nil {
		return gerror.GetError(err)
	}
	return nil
}

func zip(a1, a2 []string) []string {
	r := make([]string, 2*len(a1))
	for i, e := range a1 {
		r[i*2] = e
		r[i*2+1] = a2[i]
	}
	return r
}
func (b *Base) PubHeaderFile(ctx context.Context, prefix string, fileHeader *multipart.FileHeader) (string, *gerror.Error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", gerror.GetError(err)
	}

	newName := strings.NewReplacer(zip(chars, make([]string, len(chars)))...).Replace(fileHeader.Filename)

	name := prefix + time.Now().Format("2006_01_02-15_04") + "_" + newName
	info, err := b.Client.PutObject(ctx,
		b.Bucket,
		name,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
			Expires:     time.Now().UTC().Add(time.Second * 30),
		})
	if err != nil {
		return "", gerror.GetError(err)
	}
	return info.Key, nil
}
func (b *Base) GenerateUrl(ctx context.Context, key string, duration time.Duration) (*url.URL, *gerror.Error) {
	reqParams := make(url.Values)
	//name := strings.Split(key, "/")
	//h := "attachment; filename=\"" + name[len(name)-1] + "\""
	//reqParams.Set("response-content-disposition", h)
	u, err := b.Client.PresignedGetObject(ctx, b.Bucket, key, duration, reqParams)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	return u, nil
}

func (b *Base) CopyFrom(ctx context.Context, dst *Base, prefix string, key string) (string, *gerror.Error) {

	reader, err := b.Client.GetObject(ctx, dst.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return "", gerror.GetError(err)
	}
	defer reader.Close()
	stat, err := reader.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	info, err := b.Client.PutObject(ctx,
		b.Bucket,
		prefix+key,
		reader,
		stat.Size,
		minio.PutObjectOptions{
			ContentType: stat.ContentType,
			Expires:     time.Now().UTC().Add(time.Second * 30),
		})
	if err != nil {
		return "", gerror.GetError(err)
	}
	return info.Key, nil
}
func (b *Base) DeleteByKey(ctx context.Context, key string) *gerror.Error {
	err := b.Client.RemoveObject(ctx, b.Bucket, key, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
	if err != nil {
		return gerror.GetError(err)
	}
	return nil
}

func (b *Base) ListObject(ctx context.Context) {
	opts := minio.ListObjectsOptions{
		Recursive: true,
	}
	// List all objects from a bucket-name with a matching prefix.
	for object := range b.Client.ListObjects(ctx, b.Bucket, opts) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println(object)
	}
}
