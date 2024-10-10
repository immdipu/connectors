package atlassian

import (
	"github.com/amp-labs/connectors/common/paramsbuilder"
	"github.com/amp-labs/connectors/internal/deep/dpvars"
	"github.com/amp-labs/connectors/internal/deep/requirements"
)

const cloudIdKey = "cloudId"

// AuthMetadataVars is a complete list of authentication metadata associated with connector.
// This model serves as a documentation of map[string]string contents.
type AuthMetadataVars struct {
	CloudID string
}

func (v *AuthMetadataVars) FromMap(dictionary map[string]string) {
	v.CloudID = dictionary[cloudIdKey]
}

func (v *AuthMetadataVars) ToMap() map[string]string {
	return map[string]string{
		cloudIdKey: v.CloudID,
	}
}

func (v *AuthMetadataVars) GetSubstitutionPlans() []paramsbuilder.SubstitutionPlan {
	return nil
}

func NewAuthMetadataVars() *AuthMetadataVars {
	return &AuthMetadataVars{CloudID: ""}
}

var _ dpvars.MetadataVariables = &AuthMetadataVars{}

func (v *AuthMetadataVars) Satisfies() requirements.Dependency {
	return requirements.Dependency{
		ID:          "metadataVariables",
		Constructor: NewAuthMetadataVars,
		Interface:   new(dpvars.MetadataVariables),
	}
}
