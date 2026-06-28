package globalrouter

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type ContentPart struct {
	Type     string         `json:"type"`
	Text     string         `json:"text,omitempty"`
	ImageURL any            `json:"image_url,omitempty"`
	Audio    map[string]any `json:"audio,omitempty"`
	Video    map[string]any `json:"video,omitempty"`
}

type Message struct {
	Role       Role             `json:"role"`
	Content    any              `json:"content,omitempty"`
	Name       string           `json:"name,omitempty"`
	ToolCalls  []map[string]any `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

type ChatRequest struct {
	Model          string           `json:"model"`
	Provider       string           `json:"provider,omitempty"`
	Messages       []Message        `json:"messages,omitempty"`
	Temperature    *float64         `json:"temperature,omitempty"`
	TopP           *float64         `json:"top_p,omitempty"`
	MaxTokens      *int             `json:"max_tokens,omitempty"`
	Stream         bool             `json:"stream,omitempty"`
	StreamOptions  map[string]any   `json:"stream_options,omitempty"`
	Tools          []map[string]any `json:"tools,omitempty"`
	ToolChoice     any              `json:"tool_choice,omitempty"`
	ResponseFormat map[string]any   `json:"response_format,omitempty"`
	Stop           any              `json:"stop,omitempty"`
	User           string           `json:"user,omitempty"`
	Router         map[string]any   `json:"_router,omitempty"`
}

type ChoiceDelta struct {
	Role      Role             `json:"role,omitempty"`
	Content   string           `json:"content,omitempty"`
	ToolCalls []map[string]any `json:"tool_calls,omitempty"`
}

type ChatChoice struct {
	Index        int          `json:"index"`
	Message      *Message     `json:"message,omitempty"`
	Delta        *ChoiceDelta `json:"delta,omitempty"`
	FinishReason string       `json:"finish_reason,omitempty"`
}

type Usage struct {
	PromptTokens     int      `json:"prompt_tokens"`
	CompletionTokens int      `json:"completion_tokens"`
	TotalTokens      int      `json:"total_tokens"`
	CachedTokens     *int     `json:"cached_tokens,omitempty"`
	Cost             *float64 `json:"cost,omitempty"`
}

type ChatResponse struct {
	ID              string       `json:"id"`
	Object          string       `json:"object,omitempty"`
	Created         int64        `json:"created,omitempty"`
	Model           string       `json:"model"`
	Choices         []ChatChoice `json:"choices"`
	Usage           *Usage       `json:"usage,omitempty"`
	RouterProvider  string       `json:"router_provider,omitempty"`
	RouterRequestID string       `json:"router_request_id,omitempty"`
}

type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type ModelProviderSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Model struct {
	ID                  string                 `json:"id"`
	Object              string                 `json:"object"`
	OwnedBy             string                 `json:"owned_by"`
	Providers           []ModelProviderSummary `json:"providers,omitempty"`
	DisplayName         string                 `json:"display_name"`
	Category            string                 `json:"category"`
	Subtype             string                 `json:"subtype,omitempty"`
	Modality            string                 `json:"modality"`
	Capabilities        []string               `json:"capabilities"`
	Tags                []string               `json:"tags"`
	InputModalities     []string               `json:"input_modalities"`
	OutputModalities    []string               `json:"output_modalities"`
	SupportedParameters []string               `json:"supported_parameters"`
	ExecutionMode       string                 `json:"execution_mode"`
	BillingUnit         string                 `json:"billing_unit"`
	BillingModelKey     string                 `json:"billing_model_key,omitempty"`
	ContextWindow       *int                   `json:"context_window,omitempty"`
	InputPricePer1M     float64                `json:"input_price_per_1m"`
	OutputPricePer1M    float64                `json:"output_price_per_1m"`
	DefaultTimeoutSecs  *int                   `json:"default_timeout_seconds,omitempty"`
	ProviderStatus      string                 `json:"provider_status,omitempty"`
	ProviderConfigured  bool                   `json:"provider_configured"`
	Routable            bool                   `json:"routable"`
}

type EmbeddingsRequest struct {
	Model          string `json:"model"`
	Input          any    `json:"input"`
	EncodingFormat string `json:"encoding_format,omitempty"`
	Dimensions     *int   `json:"dimensions,omitempty"`
	User           string `json:"user,omitempty"`
}

type SuperResolutionRequest struct {
	Provider   string `json:"provider,omitempty"`
	Mode       string `json:"mode,omitempty"`
	Resolution string `json:"resolution,omitempty"`
}

type ImageGenerationRequest struct {
	Model          string                  `json:"model"`
	Prompt         string                  `json:"prompt"`
	N              *int                    `json:"n,omitempty"`
	Size           string                  `json:"size,omitempty"`
	ResponseFormat string                  `json:"response_format,omitempty"`
	Quality        string                  `json:"quality,omitempty"`
	Style          string                  `json:"style,omitempty"`
	SR             *SuperResolutionRequest `json:"sr,omitempty"`
	User           string                  `json:"user,omitempty"`
}

type AudioSpeechRequest struct {
	Model          string   `json:"model"`
	Input          string   `json:"input"`
	Voice          string   `json:"voice"`
	ResponseFormat string   `json:"response_format,omitempty"`
	Speed          *float64 `json:"speed,omitempty"`
}

type AudioTranscriptionRequest struct {
	Model          string   `json:"model"`
	FileURL        string   `json:"file_url,omitempty"`
	Language       string   `json:"language,omitempty"`
	Prompt         string   `json:"prompt,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"`
	Temperature    *float64 `json:"temperature,omitempty"`
}

type GenerationRequest struct {
	Model          string                  `json:"model"`
	Provider       string                  `json:"provider,omitempty"`
	Prompt         string                  `json:"prompt"`
	Input          map[string]any          `json:"input,omitempty"`
	Routing        map[string]any          `json:"routing,omitempty"`
	SR             *SuperResolutionRequest `json:"sr,omitempty"`
	WebhookURL     string                  `json:"webhook_url,omitempty"`
	Metadata       map[string]any          `json:"metadata,omitempty"`
	Priority       *int                    `json:"priority,omitempty"`
	IdempotencyKey string                  `json:"idempotency_key,omitempty"`
}

type TaskType string

const (
	TaskTypeImageGeneration  TaskType = "image_generation"
	TaskTypeImageEdit        TaskType = "image_edit"
	TaskTypeVideoGeneration  TaskType = "video_generation"
	TaskTypeAudioGeneration  TaskType = "audio_generation"
	TaskTypeThreeDGeneration TaskType = "3d_generation"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusQueued    TaskStatus = "queued"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusSucceeded TaskStatus = "succeeded"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCanceled  TaskStatus = "canceled"
	TaskStatusTimeout   TaskStatus = "timeout"
)

type TaskCreateRequest struct {
	Type           TaskType       `json:"type"`
	Model          string         `json:"model"`
	Capability     string         `json:"capability,omitempty"`
	Input          map[string]any `json:"input"`
	Options        map[string]any `json:"options,omitempty"`
	ResponseMode   map[string]any `json:"response_mode,omitempty"`
	Routing        map[string]any `json:"routing,omitempty"`
	WebhookURL     string         `json:"webhook_url,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	Priority       *int           `json:"priority,omitempty"`
	IdempotencyKey string         `json:"idempotency_key,omitempty"`
}

type TaskArtifact struct {
	Kind       string `json:"kind"`
	URL        string `json:"url"`
	ExpiresAt  *int64 `json:"expires_at,omitempty"`
	MimeType   string `json:"mime_type"`
	SizeBytes  *int64 `json:"size_bytes,omitempty"`
	W          *int   `json:"w,omitempty"`
	H          *int   `json:"h,omitempty"`
	DurationMS *int   `json:"duration_ms,omitempty"`
}

type TaskResponse struct {
	ID              string            `json:"id"`
	Object          string            `json:"object"`
	Status          TaskStatus        `json:"status"`
	Progress        float64           `json:"progress"`
	ProgressDetail  map[string]any    `json:"progress_detail,omitempty"`
	Type            TaskType          `json:"type"`
	Capability      string            `json:"capability,omitempty"`
	Model           string            `json:"model"`
	Input           map[string]any    `json:"input,omitempty"`
	Output          []TaskArtifact    `json:"output,omitempty"`
	Result          map[string]any    `json:"result,omitempty"`
	Error           map[string]any    `json:"error,omitempty"`
	Usage           map[string]any    `json:"usage,omitempty"`
	Attempt         map[string]int    `json:"attempt,omitempty"`
	Links           map[string]string `json:"links,omitempty"`
	Provider        string            `json:"provider,omitempty"`
	CostUSD         *float64          `json:"cost_usd,omitempty"`
	BillingModelKey string            `json:"billing_model_key,omitempty"`
	DiscountID      string            `json:"discount_id,omitempty"`
	DiscountFactor  float64           `json:"discount_factor,omitempty"`
	ListCost        float64           `json:"list_cost,omitempty"`
	DiscountAmount  float64           `json:"discount_amount,omitempty"`
	ActualCost      float64           `json:"actual_cost,omitempty"`
	WebhookURL      string            `json:"webhook_url,omitempty"`
	Metadata        map[string]any    `json:"metadata,omitempty"`
	Priority        int               `json:"priority,omitempty"`
	StartedAt       string            `json:"started_at,omitempty"`
	FinishedAt      string            `json:"finished_at,omitempty"`
	CompletedAt     string            `json:"completed_at,omitempty"`
	CanceledAt      string            `json:"canceled_at,omitempty"`
	CreatedAt       string            `json:"created_at,omitempty"`
	UpdatedAt       string            `json:"updated_at,omitempty"`
}

type TaskListResponse struct {
	Items      []TaskResponse `json:"items"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	NextCursor string         `json:"next_cursor,omitempty"`
}

type TaskBatchResponse struct {
	BatchID  string         `json:"batch_id"`
	Items    []TaskResponse `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Summary  map[string]any `json:"summary,omitempty"`
}

type TaskEvent struct {
	EventID   string         `json:"event_id"`
	TaskID    string         `json:"task_id"`
	Attempt   int            `json:"attempt"`
	EventType string         `json:"event_type"`
	Status    TaskStatus     `json:"status"`
	Progress  float64        `json:"progress"`
	Payload   map[string]any `json:"payload,omitempty"`
	CreatedAt string         `json:"created_at"`
}
