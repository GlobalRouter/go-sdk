package globalrouter

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type TasksResource struct {
	client *Client
}

type ListTasksOptions struct {
	Status            TaskStatus
	Type              TaskType
	Capability        string
	Model             string
	CreatedAfter      string
	CreatedBefore     string
	RequestID         string
	MetadataProjectID string
	TenantID          string
	Cursor            string
	Page              *int
	PageSize          *int
}

type GetTaskOptions struct {
	Wait bool
}

func (r *TasksResource) Create(ctx context.Context, request TaskCreateRequest, opts ...RequestOption) (*TaskResponse, error) {
	var out TaskResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/tasks", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *TasksResource) List(ctx context.Context, options *ListTasksOptions, opts ...RequestOption) (*TaskListResponse, error) {
	params := url.Values{}
	if options != nil {
		addString(params, "status", string(options.Status))
		addString(params, "type", string(options.Type))
		addString(params, "capability", options.Capability)
		addString(params, "model", options.Model)
		addString(params, "created_after", options.CreatedAfter)
		addString(params, "created_before", options.CreatedBefore)
		addString(params, "request_id", options.RequestID)
		addString(params, "metadata.project_id", options.MetadataProjectID)
		addString(params, "tenant_id", options.TenantID)
		addString(params, "cursor", options.Cursor)
		addInt(params, "page", options.Page)
		addInt(params, "page_size", options.PageSize)
	}
	var out TaskListResponse
	if err := r.client.doJSON(ctx, http.MethodGet, "/v1/tasks", params, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *TasksResource) Get(ctx context.Context, taskID string, options *GetTaskOptions, opts ...RequestOption) (*TaskResponse, error) {
	params := url.Values{}
	if options != nil && options.Wait {
		params.Set("wait", "1")
	}
	var out TaskResponse
	if err := r.client.doJSON(ctx, http.MethodGet, "/v1/tasks/"+cleanPathValue(taskID), params, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *TasksResource) CreateBatch(ctx context.Context, requests []TaskCreateRequest, opts ...RequestOption) (*TaskBatchResponse, error) {
	var out TaskBatchResponse
	body := map[string]any{"tasks": requests}
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/tasks/batch", nil, body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *TasksResource) GetBatch(ctx context.Context, batchID string, opts ...RequestOption) (*TaskBatchResponse, error) {
	var out TaskBatchResponse
	if err := r.client.doJSON(ctx, http.MethodGet, "/v1/tasks/batch/"+cleanPathValue(batchID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *TasksResource) Events(ctx context.Context, taskID string, opts ...RequestOption) (*SSEStream[TaskEvent], error) {
	res, err := r.client.doStream(ctx, http.MethodGet, "/v1/tasks/"+cleanPathValue(taskID)+"/events", nil, nil, opts...)
	if err != nil {
		return nil, err
	}
	return newSSEStream[TaskEvent](res), nil
}

func (r *TasksResource) EventsMany(ctx context.Context, taskIDs []string, opts ...RequestOption) (*SSEStream[TaskEvent], error) {
	params := url.Values{}
	params.Set("task_id", strings.Join(taskIDs, ","))
	res, err := r.client.doStream(ctx, http.MethodGet, "/v1/tasks/events", params, nil, opts...)
	if err != nil {
		return nil, err
	}
	return newSSEStream[TaskEvent](res), nil
}

func (r *TasksResource) Cancel(ctx context.Context, taskID string, opts ...RequestOption) (*TaskResponse, error) {
	var out TaskResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/tasks/"+cleanPathValue(taskID)+"/cancel", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *TasksResource) Retry(ctx context.Context, taskID string, opts ...RequestOption) (*TaskResponse, error) {
	var out TaskResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/tasks/"+cleanPathValue(taskID)+"/retry", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
