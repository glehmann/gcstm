package lib

import (
	"context"
	"errors"
	"reflect"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func ListObjectsAtRestoreTime(ctx *context.Context, bucket *storage.BucketHandle, restoreTime time.Time) (map[string]*storage.ObjectAttrs, error) {
	query := &storage.Query{Prefix: "", Versions: true}
	it := bucket.Objects(*ctx, query)
	objects := map[string]*storage.ObjectAttrs{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return map[string]*storage.ObjectAttrs{}, err
		}
		if (attrs.Updated.Before(restoreTime) || attrs.Updated.Equal(restoreTime)) &&
			(restoreTime.Before(attrs.Deleted) || attrs.Deleted.IsZero()) {
			objects[attrs.Name] = attrs
		}
	}
	return objects, nil
}

func ListCurrentObjects(ctx *context.Context, bucket *storage.BucketHandle) (map[string]*storage.ObjectAttrs, error) {
	query := &storage.Query{Prefix: ""}
	it := bucket.Objects(*ctx, query)
	objects := map[string]*storage.ObjectAttrs{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return map[string]*storage.ObjectAttrs{}, err
		}
		objects[attrs.Name] = attrs
	}
	return objects, nil
}

func PlanRestore(ctx *context.Context, bucket *storage.BucketHandle, restoreTime time.Time) (map[string]PlanElement, error) {
	restoreObjects, err := ListObjectsAtRestoreTime(ctx, bucket, restoreTime)
	if err != nil {
		return map[string]PlanElement{}, err
	}
	currentObjects, err := ListCurrentObjects(ctx, bucket)
	if err != nil {
		return map[string]PlanElement{}, err
	}
	planElements := map[string]PlanElement{}
	for name, restoreAttrs := range restoreObjects {
		currentAttrs, ok := currentObjects[name]
		if ok {
			// check if the content and metadata are the same
			if restoreAttrs.CRC32C != currentAttrs.CRC32C {
				planElements[name] = PlanElement{RestoreObject, restoreAttrs, currentAttrs}
			} else if !FullMetadataEqual(restoreAttrs, currentAttrs) {
				planElements[name] = PlanElement{RestoreMetadata, restoreAttrs, currentAttrs}
			}
		} else {
			planElements[name] = PlanElement{RestoreObject, restoreAttrs, nil}
		}
	}
	for name, currentAttrs := range currentObjects {
		if _, ok := restoreObjects[name]; !ok {
			planElements[name] = PlanElement{Delete, nil, currentAttrs}
		}
	}
	return planElements, nil
}

var timeFormats = []string{
	"2006-01-02 15:04:05.999999999 -0700 MST",
	"2006-01-02",
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

func ParseTime(timeStr string) (time.Time, error) {
	for _, timeFormat := range timeFormats {
		t, err := time.Parse(timeFormat, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("No matching format found")
}

type Action int

const (
	RestoreObject Action = iota
	RestoreMetadata
	Delete
)

type PlanElement struct {
	Action       Action
	RestoreAttrs *storage.ObjectAttrs
	CurrentAttrs *storage.ObjectAttrs
}

//go:generate stringer -type=Action

func FullMetadataEqual(attrs1 *storage.ObjectAttrs, attrs2 *storage.ObjectAttrs) bool {
	if !reflect.DeepEqual(attrs1.Metadata, attrs2.Metadata) {
		return false
	}
	if attrs1.ContentType != attrs2.ContentType {
		return false
	}
	if attrs1.ContentLanguage != attrs2.ContentLanguage {
		return false
	}
	if attrs1.CacheControl != attrs2.CacheControl {
		return false
	}
	if !reflect.DeepEqual(attrs1.ACL, attrs2.ACL) {
		return false
	}
	if attrs1.Owner != attrs2.Owner {
		return false
	}
	if attrs1.ContentEncoding != attrs2.ContentEncoding {
		return false
	}
	if attrs1.CustomerKeySHA256 != attrs2.CustomerKeySHA256 {
		return false
	}
	if attrs1.KMSKeyName != attrs2.KMSKeyName {
		return false
	}
	if attrs1.CustomTime != attrs2.CustomTime {
		return false
	}
	return true
}
