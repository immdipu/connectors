package scrapper

import (
	"errors"
	"fmt"

	"github.com/amp-labs/connectors/common"
)

var ErrObjectNotFound = errors.New("object not found")

func (r *ObjectMetadataResult) Select(objectNames []string) (*common.ListObjectMetadataResult, error) {
	if len(objectNames) == 0 {
		return nil, common.ErrMissingObjects
	}

	list := &common.ListObjectMetadataResult{
		Result: make(map[string]common.ObjectMetadata),
		Errors: nil,
	}

	// Convert and return only listed objects
	for _, name := range objectNames {
		if v, ok := r.Result[name]; ok {
			// move metadata from scrapper object to common object
			list.Result[name] = common.ObjectMetadata{
				DisplayName: v.DisplayName,
				FieldsMap:   v.FieldsMap,
			}
		} else {
			return nil, fmt.Errorf("%w: unknown object [%v]", ErrObjectNotFound, name)
		}
	}

	return list, nil
}

// SelectOne reads one object metadata from the static file.
func (r *ObjectMetadataResult) SelectOne(objectName string) (*common.ObjectMetadata, error) {
	mtd := common.ObjectMetadata{
		FieldsMap: make(map[string]string),
	}

	if v, ok := r.Result[objectName]; ok {
		mtd = common.ObjectMetadata{
			DisplayName: v.DisplayName,
			FieldsMap:   v.FieldsMap,
		}
	} else {
		return nil, fmt.Errorf("%w: unknown object [%v]", ErrObjectNotFound, objectName)
	}

	return &mtd, nil
}
