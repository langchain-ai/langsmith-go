// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package pagination

import (
	"net/http"
	"strconv"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

type OffsetPaginationTopLevelArray[T any] struct {
	Items []T                               `json:"-,inline"`
	JSON  offsetPaginationTopLevelArrayJSON `json:"-"`
	cfg   *requestconfig.RequestConfig
	res   *http.Response
}

// offsetPaginationTopLevelArrayJSON contains the JSON metadata for the struct
// [OffsetPaginationTopLevelArray[T]]
type offsetPaginationTopLevelArrayJSON struct {
	Items       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *OffsetPaginationTopLevelArray[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r offsetPaginationTopLevelArrayJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *OffsetPaginationTopLevelArray[T]) GetNextPage() (res *OffsetPaginationTopLevelArray[T], err error) {
	if len(r.Items) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)

	q := cfg.Request.URL.Query()
	offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	length := int64(len(r.Items))
	next := offset + length

	if length > 0 && next != 0 {
		err = cfg.Apply(option.WithQuery("offset", strconv.FormatInt(next, 10)))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *OffsetPaginationTopLevelArray[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &OffsetPaginationTopLevelArray[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type OffsetPaginationTopLevelArrayAutoPager[T any] struct {
	page *OffsetPaginationTopLevelArray[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewOffsetPaginationTopLevelArrayAutoPager[T any](page *OffsetPaginationTopLevelArray[T], err error) *OffsetPaginationTopLevelArrayAutoPager[T] {
	return &OffsetPaginationTopLevelArrayAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *OffsetPaginationTopLevelArrayAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Items) == 0 {
		return false
	}
	if r.idx >= len(r.page.Items) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Items) == 0 {
			return false
		}
	}
	r.cur = r.page.Items[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *OffsetPaginationTopLevelArrayAutoPager[T]) Current() T {
	return r.cur
}

func (r *OffsetPaginationTopLevelArrayAutoPager[T]) Err() error {
	return r.err
}

func (r *OffsetPaginationTopLevelArrayAutoPager[T]) Index() int {
	return r.run
}

type OffsetPaginationIssues[T any] struct {
	Items []T                        `json:"-,inline"`
	JSON  offsetPaginationIssuesJSON `json:"-"`
	cfg   *requestconfig.RequestConfig
	res   *http.Response
}

// offsetPaginationIssuesJSON contains the JSON metadata for the struct
// [OffsetPaginationIssues[T]]
type offsetPaginationIssuesJSON struct {
	Items       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *OffsetPaginationIssues[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r offsetPaginationIssuesJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *OffsetPaginationIssues[T]) GetNextPage() (res *OffsetPaginationIssues[T], err error) {
	if len(r.Items) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)

	q := cfg.Request.URL.Query()
	offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	length := int64(len(r.Items))
	next := offset + length

	if length > 0 && next != 0 {
		err = cfg.Apply(option.WithQuery("offset", strconv.FormatInt(next, 10)))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *OffsetPaginationIssues[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &OffsetPaginationIssues[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type OffsetPaginationIssuesAutoPager[T any] struct {
	page *OffsetPaginationIssues[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewOffsetPaginationIssuesAutoPager[T any](page *OffsetPaginationIssues[T], err error) *OffsetPaginationIssuesAutoPager[T] {
	return &OffsetPaginationIssuesAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *OffsetPaginationIssuesAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Items) == 0 {
		return false
	}
	if r.idx >= len(r.page.Items) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Items) == 0 {
			return false
		}
	}
	r.cur = r.page.Items[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *OffsetPaginationIssuesAutoPager[T]) Current() T {
	return r.cur
}

func (r *OffsetPaginationIssuesAutoPager[T]) Err() error {
	return r.err
}

func (r *OffsetPaginationIssuesAutoPager[T]) Index() int {
	return r.run
}

type OffsetPaginationRepos[T any] struct {
	Repos []T                       `json:"repos"`
	Total int64                     `json:"total"`
	JSON  offsetPaginationReposJSON `json:"-"`
	cfg   *requestconfig.RequestConfig
	res   *http.Response
}

// offsetPaginationReposJSON contains the JSON metadata for the struct
// [OffsetPaginationRepos[T]]
type offsetPaginationReposJSON struct {
	Repos       apijson.Field
	Total       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *OffsetPaginationRepos[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r offsetPaginationReposJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *OffsetPaginationRepos[T]) GetNextPage() (res *OffsetPaginationRepos[T], err error) {
	if len(r.Repos) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)

	q := cfg.Request.URL.Query()
	offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	length := int64(len(r.Repos))
	next := offset + length

	if length > 0 && next != 0 {
		err = cfg.Apply(option.WithQuery("offset", strconv.FormatInt(next, 10)))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *OffsetPaginationRepos[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &OffsetPaginationRepos[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type OffsetPaginationReposAutoPager[T any] struct {
	page *OffsetPaginationRepos[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewOffsetPaginationReposAutoPager[T any](page *OffsetPaginationRepos[T], err error) *OffsetPaginationReposAutoPager[T] {
	return &OffsetPaginationReposAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *OffsetPaginationReposAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Repos) == 0 {
		return false
	}
	if r.idx >= len(r.page.Repos) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Repos) == 0 {
			return false
		}
	}
	r.cur = r.page.Repos[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *OffsetPaginationReposAutoPager[T]) Current() T {
	return r.cur
}

func (r *OffsetPaginationReposAutoPager[T]) Err() error {
	return r.err
}

func (r *OffsetPaginationReposAutoPager[T]) Index() int {
	return r.run
}

type OffsetPaginationCommits[T any] struct {
	Commits []T                         `json:"commits"`
	Total   int64                       `json:"total"`
	JSON    offsetPaginationCommitsJSON `json:"-"`
	cfg     *requestconfig.RequestConfig
	res     *http.Response
}

// offsetPaginationCommitsJSON contains the JSON metadata for the struct
// [OffsetPaginationCommits[T]]
type offsetPaginationCommitsJSON struct {
	Commits     apijson.Field
	Total       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *OffsetPaginationCommits[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r offsetPaginationCommitsJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *OffsetPaginationCommits[T]) GetNextPage() (res *OffsetPaginationCommits[T], err error) {
	if len(r.Commits) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)

	q := cfg.Request.URL.Query()
	offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	length := int64(len(r.Commits))
	next := offset + length

	if length > 0 && next != 0 {
		err = cfg.Apply(option.WithQuery("offset", strconv.FormatInt(next, 10)))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *OffsetPaginationCommits[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &OffsetPaginationCommits[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type OffsetPaginationCommitsAutoPager[T any] struct {
	page *OffsetPaginationCommits[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewOffsetPaginationCommitsAutoPager[T any](page *OffsetPaginationCommits[T], err error) *OffsetPaginationCommitsAutoPager[T] {
	return &OffsetPaginationCommitsAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *OffsetPaginationCommitsAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Commits) == 0 {
		return false
	}
	if r.idx >= len(r.page.Commits) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Commits) == 0 {
			return false
		}
	}
	r.cur = r.page.Commits[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *OffsetPaginationCommitsAutoPager[T]) Current() T {
	return r.cur
}

func (r *OffsetPaginationCommitsAutoPager[T]) Err() error {
	return r.err
}

func (r *OffsetPaginationCommitsAutoPager[T]) Index() int {
	return r.run
}

type OffsetPaginationOnlineEvaluators[T any] struct {
	Evaluators []T                                  `json:"evaluators"`
	Total      int64                                `json:"total"`
	JSON       offsetPaginationOnlineEvaluatorsJSON `json:"-"`
	cfg        *requestconfig.RequestConfig
	res        *http.Response
}

// offsetPaginationOnlineEvaluatorsJSON contains the JSON metadata for the struct
// [OffsetPaginationOnlineEvaluators[T]]
type offsetPaginationOnlineEvaluatorsJSON struct {
	Evaluators  apijson.Field
	Total       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *OffsetPaginationOnlineEvaluators[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r offsetPaginationOnlineEvaluatorsJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *OffsetPaginationOnlineEvaluators[T]) GetNextPage() (res *OffsetPaginationOnlineEvaluators[T], err error) {
	if len(r.Evaluators) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)

	q := cfg.Request.URL.Query()
	offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	length := int64(len(r.Evaluators))
	next := offset + length

	if length > 0 && next != 0 {
		err = cfg.Apply(option.WithQuery("offset", strconv.FormatInt(next, 10)))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *OffsetPaginationOnlineEvaluators[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &OffsetPaginationOnlineEvaluators[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type OffsetPaginationOnlineEvaluatorsAutoPager[T any] struct {
	page *OffsetPaginationOnlineEvaluators[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewOffsetPaginationOnlineEvaluatorsAutoPager[T any](page *OffsetPaginationOnlineEvaluators[T], err error) *OffsetPaginationOnlineEvaluatorsAutoPager[T] {
	return &OffsetPaginationOnlineEvaluatorsAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *OffsetPaginationOnlineEvaluatorsAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Evaluators) == 0 {
		return false
	}
	if r.idx >= len(r.page.Evaluators) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Evaluators) == 0 {
			return false
		}
	}
	r.cur = r.page.Evaluators[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *OffsetPaginationOnlineEvaluatorsAutoPager[T]) Current() T {
	return r.cur
}

func (r *OffsetPaginationOnlineEvaluatorsAutoPager[T]) Err() error {
	return r.err
}

func (r *OffsetPaginationOnlineEvaluatorsAutoPager[T]) Index() int {
	return r.run
}

type OffsetPaginationInsightsClusteringJobs[T any] struct {
	ClusteringJobs []T                                        `json:"clustering_jobs"`
	JSON           offsetPaginationInsightsClusteringJobsJSON `json:"-"`
	cfg            *requestconfig.RequestConfig
	res            *http.Response
}

// offsetPaginationInsightsClusteringJobsJSON contains the JSON metadata for the
// struct [OffsetPaginationInsightsClusteringJobs[T]]
type offsetPaginationInsightsClusteringJobsJSON struct {
	ClusteringJobs apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *OffsetPaginationInsightsClusteringJobs[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r offsetPaginationInsightsClusteringJobsJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *OffsetPaginationInsightsClusteringJobs[T]) GetNextPage() (res *OffsetPaginationInsightsClusteringJobs[T], err error) {
	if len(r.ClusteringJobs) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)

	q := cfg.Request.URL.Query()
	offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	length := int64(len(r.ClusteringJobs))
	next := offset + length

	if length > 0 && next != 0 {
		err = cfg.Apply(option.WithQuery("offset", strconv.FormatInt(next, 10)))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *OffsetPaginationInsightsClusteringJobs[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &OffsetPaginationInsightsClusteringJobs[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type OffsetPaginationInsightsClusteringJobsAutoPager[T any] struct {
	page *OffsetPaginationInsightsClusteringJobs[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewOffsetPaginationInsightsClusteringJobsAutoPager[T any](page *OffsetPaginationInsightsClusteringJobs[T], err error) *OffsetPaginationInsightsClusteringJobsAutoPager[T] {
	return &OffsetPaginationInsightsClusteringJobsAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *OffsetPaginationInsightsClusteringJobsAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.ClusteringJobs) == 0 {
		return false
	}
	if r.idx >= len(r.page.ClusteringJobs) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.ClusteringJobs) == 0 {
			return false
		}
	}
	r.cur = r.page.ClusteringJobs[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *OffsetPaginationInsightsClusteringJobsAutoPager[T]) Current() T {
	return r.cur
}

func (r *OffsetPaginationInsightsClusteringJobsAutoPager[T]) Err() error {
	return r.err
}

func (r *OffsetPaginationInsightsClusteringJobsAutoPager[T]) Index() int {
	return r.run
}

type CursorPaginationCursors struct {
	Next string                      `json:"next"`
	JSON cursorPaginationCursorsJSON `json:"-"`
}

// cursorPaginationCursorsJSON contains the JSON metadata for the struct
// [CursorPaginationCursors]
type cursorPaginationCursorsJSON struct {
	Next        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CursorPaginationCursors) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r cursorPaginationCursorsJSON) RawJSON() string {
	return r.raw
}

type CursorPagination[T any] struct {
	Runs    []T                     `json:"runs"`
	Cursors CursorPaginationCursors `json:"cursors"`
	JSON    cursorPaginationJSON    `json:"-"`
	cfg     *requestconfig.RequestConfig
	res     *http.Response
}

// cursorPaginationJSON contains the JSON metadata for the struct
// [CursorPagination[T]]
type cursorPaginationJSON struct {
	Runs        apijson.Field
	Cursors     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CursorPagination[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r cursorPaginationJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *CursorPagination[T]) GetNextPage() (res *CursorPagination[T], err error) {
	if len(r.Runs) == 0 {
		return nil, nil
	}
	next := r.Cursors.Next
	if len(next) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)
	err = cfg.Apply(option.WithQuery("cursor", next))
	if err != nil {
		return nil, err
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *CursorPagination[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &CursorPagination[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type CursorPaginationAutoPager[T any] struct {
	page *CursorPagination[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewCursorPaginationAutoPager[T any](page *CursorPagination[T], err error) *CursorPaginationAutoPager[T] {
	return &CursorPaginationAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *CursorPaginationAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Runs) == 0 {
		return false
	}
	if r.idx >= len(r.page.Runs) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Runs) == 0 {
			return false
		}
	}
	r.cur = r.page.Runs[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *CursorPaginationAutoPager[T]) Current() T {
	return r.cur
}

func (r *CursorPaginationAutoPager[T]) Err() error {
	return r.err
}

func (r *CursorPaginationAutoPager[T]) Index() int {
	return r.run
}

type ItemsCursorPostPagination[T any] struct {
	Items      []T                           `json:"items"`
	NextCursor string                        `json:"next_cursor"`
	JSON       itemsCursorPostPaginationJSON `json:"-"`
	cfg        *requestconfig.RequestConfig
	res        *http.Response
}

// itemsCursorPostPaginationJSON contains the JSON metadata for the struct
// [ItemsCursorPostPagination[T]]
type itemsCursorPostPaginationJSON struct {
	Items       apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ItemsCursorPostPagination[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r itemsCursorPostPaginationJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *ItemsCursorPostPagination[T]) GetNextPage() (res *ItemsCursorPostPagination[T], err error) {
	if len(r.Items) == 0 {
		return nil, nil
	}
	next := r.NextCursor
	if len(next) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)
	err = cfg.Apply(option.WithQuery("cursor", next))
	if err != nil {
		return nil, err
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *ItemsCursorPostPagination[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &ItemsCursorPostPagination[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type ItemsCursorPostPaginationAutoPager[T any] struct {
	page *ItemsCursorPostPagination[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewItemsCursorPostPaginationAutoPager[T any](page *ItemsCursorPostPagination[T], err error) *ItemsCursorPostPaginationAutoPager[T] {
	return &ItemsCursorPostPaginationAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *ItemsCursorPostPaginationAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Items) == 0 {
		return false
	}
	if r.idx >= len(r.page.Items) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Items) == 0 {
			return false
		}
	}
	r.cur = r.page.Items[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *ItemsCursorPostPaginationAutoPager[T]) Current() T {
	return r.cur
}

func (r *ItemsCursorPostPaginationAutoPager[T]) Err() error {
	return r.err
}

func (r *ItemsCursorPostPaginationAutoPager[T]) Index() int {
	return r.run
}

type ItemsCursorGetPagination[T any] struct {
	Items      []T                          `json:"items"`
	NextCursor string                       `json:"next_cursor"`
	JSON       itemsCursorGetPaginationJSON `json:"-"`
	cfg        *requestconfig.RequestConfig
	res        *http.Response
}

// itemsCursorGetPaginationJSON contains the JSON metadata for the struct
// [ItemsCursorGetPagination[T]]
type itemsCursorGetPaginationJSON struct {
	Items       apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ItemsCursorGetPagination[T]) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r itemsCursorGetPaginationJSON) RawJSON() string {
	return r.raw
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *ItemsCursorGetPagination[T]) GetNextPage() (res *ItemsCursorGetPagination[T], err error) {
	if len(r.Items) == 0 {
		return nil, nil
	}
	next := r.NextCursor
	if len(next) == 0 {
		return nil, nil
	}
	cfg := r.cfg.Clone(r.cfg.Context)
	err = cfg.Apply(option.WithQuery("cursor", next))
	if err != nil {
		return nil, err
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *ItemsCursorGetPagination[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &ItemsCursorGetPagination[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type ItemsCursorGetPaginationAutoPager[T any] struct {
	page *ItemsCursorGetPagination[T]
	cur  T
	idx  int
	run  int
	err  error
}

func NewItemsCursorGetPaginationAutoPager[T any](page *ItemsCursorGetPagination[T], err error) *ItemsCursorGetPaginationAutoPager[T] {
	return &ItemsCursorGetPaginationAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *ItemsCursorGetPaginationAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Items) == 0 {
		return false
	}
	if r.idx >= len(r.page.Items) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Items) == 0 {
			return false
		}
	}
	r.cur = r.page.Items[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *ItemsCursorGetPaginationAutoPager[T]) Current() T {
	return r.cur
}

func (r *ItemsCursorGetPaginationAutoPager[T]) Err() error {
	return r.err
}

func (r *ItemsCursorGetPaginationAutoPager[T]) Index() int {
	return r.run
}
