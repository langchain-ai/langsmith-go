// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"time"
)

type UnionTime time.Time

func (UnionTime) ImplementsExampleGetParamsAsOfUnion()                   {}
func (UnionTime) ImplementsExampleListParamsAsOfUnion()                  {}
func (UnionTime) ImplementsExampleGetCountParamsAsOfUnion()              {}
func (UnionTime) ImplementsDatasetCloneParamsAsOfUnion()                 {}
func (UnionTime) ImplementsDatasetUpdateTagsParamsAsOfUnion()            {}
func (UnionTime) ImplementsDatasetVersionGetDiffParamsFromVersionUnion() {}
func (UnionTime) ImplementsDatasetVersionGetDiffParamsToVersionUnion()   {}
func (UnionTime) ImplementsDatasetSplitGetParamsAsOfUnion()              {}

type UnionString string

func (UnionString) ImplementsExampleNewParamsSplitUnion()                  {}
func (UnionString) ImplementsExampleGetParamsAsOfUnion()                   {}
func (UnionString) ImplementsExampleUpdateParamsSplitUnion()               {}
func (UnionString) ImplementsExampleListParamsAsOfUnion()                  {}
func (UnionString) ImplementsExampleGetCountParamsAsOfUnion()              {}
func (UnionString) ImplementsExampleBulkNewParamsBodySplitUnion()          {}
func (UnionString) ImplementsExampleBulkPatchAllParamsBodySplitUnion()     {}
func (UnionString) ImplementsExampleValidationResultSplitUnion()           {}
func (UnionString) ImplementsDatasetUpdateParamsDescriptionUnion()         {}
func (UnionString) ImplementsDatasetUpdateParamsNameUnion()                {}
func (UnionString) ImplementsDatasetUpdateParamsPatchExamplesSplitUnion()  {}
func (UnionString) ImplementsDatasetCloneParamsAsOfUnion()                 {}
func (UnionString) ImplementsDatasetUpdateTagsParamsAsOfUnion()            {}
func (UnionString) ImplementsDatasetVersionGetDiffParamsFromVersionUnion() {}
func (UnionString) ImplementsDatasetVersionGetDiffParamsToVersionUnion()   {}
func (UnionString) ImplementsDatasetGroupRunsResponseGroupsGroupKeyUnion() {}
func (UnionString) ImplementsDatasetSplitGetParamsAsOfUnion()              {}
func (UnionString) ImplementsFeedbackCreateSchemaCorrectionUnionParam()    {}
func (UnionString) ImplementsFeedbackCreateSchemaValueUnionParam()         {}
func (UnionString) ImplementsFeedbackSchemaCorrectionUnion()               {}
func (UnionString) ImplementsFeedbackSchemaValueUnion()                    {}
func (UnionString) ImplementsFeedbackUpdateParamsCorrectionUnion()         {}
func (UnionString) ImplementsFeedbackUpdateParamsValueUnion()              {}
func (UnionString) ImplementsFeedbackTokenGetParamsValueUnion()            {}
func (UnionString) ImplementsFeedbackTokenUpdateParamsCorrectionUnion()    {}
func (UnionString) ImplementsFeedbackTokenUpdateParamsValueUnion()         {}

type UnionBool bool

func (UnionBool) ImplementsFeedbackCreateSchemaScoreUnionParam() {}
func (UnionBool) ImplementsFeedbackCreateSchemaValueUnionParam() {}
func (UnionBool) ImplementsFeedbackSchemaScoreUnion()            {}
func (UnionBool) ImplementsFeedbackSchemaValueUnion()            {}
func (UnionBool) ImplementsFeedbackUpdateParamsScoreUnion()      {}
func (UnionBool) ImplementsFeedbackUpdateParamsValueUnion()      {}
func (UnionBool) ImplementsFeedbackTokenGetParamsScoreUnion()    {}
func (UnionBool) ImplementsFeedbackTokenGetParamsValueUnion()    {}
func (UnionBool) ImplementsFeedbackTokenUpdateParamsScoreUnion() {}
func (UnionBool) ImplementsFeedbackTokenUpdateParamsValueUnion() {}

type UnionInt int64

func (UnionInt) ImplementsAnnotationQueueUpdateParamsNumReviewersPerItemUnion() {}

type UnionFloat float64

func (UnionFloat) ImplementsCustomChartsSectionChartsDataValueUnion()            {}
func (UnionFloat) ImplementsCustomChartsSectionSubSectionsChartsDataValueUnion() {}
func (UnionFloat) ImplementsDatasetGroupRunsResponseGroupsGroupKeyUnion()        {}
func (UnionFloat) ImplementsFeedbackCreateSchemaScoreUnionParam()                {}
func (UnionFloat) ImplementsFeedbackCreateSchemaValueUnionParam()                {}
func (UnionFloat) ImplementsFeedbackSchemaScoreUnion()                           {}
func (UnionFloat) ImplementsFeedbackSchemaValueUnion()                           {}
func (UnionFloat) ImplementsFeedbackUpdateParamsScoreUnion()                     {}
func (UnionFloat) ImplementsFeedbackUpdateParamsValueUnion()                     {}
func (UnionFloat) ImplementsFeedbackTokenGetParamsScoreUnion()                   {}
func (UnionFloat) ImplementsFeedbackTokenGetParamsValueUnion()                   {}
func (UnionFloat) ImplementsFeedbackTokenUpdateParamsScoreUnion()                {}
func (UnionFloat) ImplementsFeedbackTokenUpdateParamsValueUnion()                {}
