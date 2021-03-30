package lib

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func ListObjectAtRestoreTime(ctx *context.Context, bucket *storage.BucketHandle, restoreTime time.Time) (map[string]*storage.ObjectAttrs, error) {
	query := &storage.Query{Prefix: "", Versions: true}
	it := bucket.Objects(*ctx, query)
	objects := map[string]*storage.ObjectAttrs{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return map[string]*storage.ObjectAttrs{}, nil
		}
		if (attrs.Updated.Before(restoreTime) || attrs.Updated.Equal(restoreTime)) &&
			(restoreTime.Before(attrs.Deleted) || attrs.Deleted.IsZero()) {
			objects[attrs.Name] = attrs
		}
	}
	return objects, nil
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
