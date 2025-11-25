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
func (UnionString) ImplementsFeedbackTokenGetParamsValueUnion()            {}
func (UnionString) ImplementsFeedbackTokenUpdateParamsValueUnion()         {}

type UnionBool bool

func (UnionBool) ImplementsFeedbackCreateSchemaScoreUnionParam() {}
func (UnionBool) ImplementsFeedbackSchemaScoreUnion()            {}
func (UnionBool) ImplementsFeedbackUpdateParamsScoreUnion()      {}
func (UnionBool) ImplementsFeedbackTokenGetParamsScoreUnion()    {}
func (UnionBool) ImplementsFeedbackTokenGetParamsValueUnion()    {}
func (UnionBool) ImplementsFeedbackTokenUpdateParamsScoreUnion() {}
func (UnionBool) ImplementsFeedbackTokenUpdateParamsValueUnion() {}
func (UnionBool) ImplementsCommitUpdateParamsSkipWebhooksUnion() {}

type UnionInt int64

func (UnionInt) ImplementsAnnotationQueueUpdateParamsNumReviewersPerItemUnion() {}

type UnionFloat float64

func (UnionFloat) ImplementsDatasetGroupRunsResponseGroupsGroupKeyUnion() {}
func (UnionFloat) ImplementsFeedbackCreateSchemaScoreUnionParam()         {}
func (UnionFloat) ImplementsFeedbackSchemaScoreUnion()                    {}
func (UnionFloat) ImplementsFeedbackUpdateParamsScoreUnion()              {}
func (UnionFloat) ImplementsFeedbackTokenGetParamsScoreUnion()            {}
func (UnionFloat) ImplementsFeedbackTokenGetParamsValueUnion()            {}
func (UnionFloat) ImplementsFeedbackTokenUpdateParamsScoreUnion()         {}
func (UnionFloat) ImplementsFeedbackTokenUpdateParamsValueUnion()         {}
