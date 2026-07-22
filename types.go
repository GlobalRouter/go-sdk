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
	Stream         bool             `json:"stream"`
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
	Resolution string `json:"resolution,omitempty"`
}

type ProviderSelection struct {
	ProviderID string                    `json:"provider_id,omitempty"`
	Options    map[string]map[string]any `json:"options,omitempty"`
}

type ImageGenerationRequest struct {
	Model             string                  `json:"model"`
	Prompt            string                  `json:"prompt"`
	AspectRatio       string                  `json:"aspect_ratio,omitempty"`
	Background        string                  `json:"background,omitempty"`
	InputReferences   []ImageTaskReference    `json:"input_references,omitempty"`
	N                 *int                    `json:"n,omitempty"`
	OutputCompression *int                    `json:"output_compression,omitempty"`
	OutputFormat      string                  `json:"output_format,omitempty"`
	Provider          *ProviderSelection      `json:"provider,omitempty"`
	Quality           string                  `json:"quality,omitempty"`
	Resolution        string                  `json:"resolution,omitempty"`
	Seed              *int                    `json:"seed,omitempty"`
	Size              string                  `json:"size,omitempty"`
	Stream            *bool                   `json:"stream,omitempty"`
	SR                *SuperResolutionRequest `json:"sr,omitempty"`
	User              string                  `json:"user,omitempty"`
}

type ImageTaskStatus string

const (
	ImageTaskStatusQueued         ImageTaskStatus = "queued"
	ImageTaskStatusSubmitting     ImageTaskStatus = "submitting"
	ImageTaskStatusSubmitRetrying ImageTaskStatus = "submit_retrying"
	ImageTaskStatusSubmitted      ImageTaskStatus = "submitted"
	ImageTaskStatusInProgress     ImageTaskStatus = "in_progress"
	ImageTaskStatusPollRetrying   ImageTaskStatus = "poll_retrying"
	ImageTaskStatusCompleted      ImageTaskStatus = "completed"
	ImageTaskStatusFailed         ImageTaskStatus = "failed"
	ImageTaskStatusExpired        ImageTaskStatus = "expired"
)

type ImageTaskReference struct {
	Type     string         `json:"type,omitempty"`
	ImageURL map[string]any `json:"image_url"`
}

type ImageTaskProvider = ProviderSelection

type VideoFrameImage struct {
	Type      string         `json:"type"`
	ImageURL  map[string]any `json:"image_url"`
	FrameType string         `json:"frame_type"`
}

type ImageTaskCreateRequest struct {
	Model           string               `json:"model"`
	Prompt          string               `json:"prompt"`
	InputReferences []ImageTaskReference `json:"input_references,omitempty"`
	N               *int                 `json:"n,omitempty"`
	Size            string               `json:"size,omitempty"`
	AspectRatio     string               `json:"aspect_ratio,omitempty"`
	Resolution      string               `json:"resolution,omitempty"`
	Quality         string               `json:"quality,omitempty"`
	Seed            *int                 `json:"seed,omitempty"`
	OutputFormat    string               `json:"output_format,omitempty"`
	ResponseFormat  string               `json:"response_format,omitempty"`
	Provider        *ImageTaskProvider   `json:"provider,omitempty"`
	WebhookURL      string               `json:"webhook_url,omitempty"`
	Metadata        map[string]any       `json:"metadata,omitempty"`
	Priority        *int                 `json:"priority,omitempty"`
	IdempotencyKey  string               `json:"idempotency_key,omitempty"`
}

type ImageTaskArtifact struct {
	Index     int     `json:"index"`
	URL       string  `json:"url,omitempty"`
	B64JSON   string  `json:"b64_json,omitempty"`
	MediaType string  `json:"media_type,omitempty"`
	Width     *int    `json:"width,omitempty"`
	Height    *int    `json:"height,omitempty"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

type ImageTaskUsage struct {
	ImageCount    *int           `json:"image_count,omitempty"`
	BillingStatus string         `json:"billing_status,omitempty"`
	CostUSD       *float64       `json:"cost_usd,omitempty"`
	UsageRecordID string         `json:"usage_record_id,omitempty"`
	Details       map[string]any `json:"details,omitempty"`
}

type ImageTaskError struct {
	Code          string         `json:"code"`
	Message       string         `json:"message"`
	ProviderError map[string]any `json:"provider_error,omitempty"`
}

type ImageTaskResponse struct {
	ID          string              `json:"id"`
	Object      string              `json:"object"`
	Status      ImageTaskStatus     `json:"status"`
	Model       string              `json:"model"`
	Provider    string              `json:"provider,omitempty"`
	PollingURL  string              `json:"polling_url,omitempty"`
	CreatedAt   int64               `json:"created_at,omitempty"`
	UpdatedAt   int64               `json:"updated_at,omitempty"`
	CompletedAt *int64              `json:"completed_at,omitempty"`
	Data        []ImageTaskArtifact `json:"data,omitempty"`
	Usage       *ImageTaskUsage     `json:"usage,omitempty"`
	Error       *ImageTaskError     `json:"error,omitempty"`
	Metadata    map[string]any      `json:"metadata,omitempty"`
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
	Model           string                  `json:"model"`
	Prompt          string                  `json:"prompt"`
	AspectRatio     string                  `json:"aspect_ratio,omitempty"`
	CallbackURL     string                  `json:"callback_url,omitempty"`
	Duration        *int                    `json:"duration,omitempty"`
	FrameImages     []VideoFrameImage       `json:"frame_images,omitempty"`
	GenerateAudio   *bool                   `json:"generate_audio,omitempty"`
	InputReferences []map[string]any        `json:"input_references,omitempty"`
	Provider        *ProviderSelection      `json:"provider,omitempty"`
	Resolution      string                  `json:"resolution,omitempty"`
	SR              *SuperResolutionRequest `json:"sr,omitempty"`
	Seed            *int                    `json:"seed,omitempty"`
	Size            string                  `json:"size,omitempty"`
}

type VideoGenerationResponse struct {
	ID           string         `json:"id"`
	PollingURL   string         `json:"polling_url"`
	Status       string         `json:"status"`
	Error        string         `json:"error,omitempty"`
	GenerationID string         `json:"generation_id,omitempty"`
	UnsignedURLs []string       `json:"unsigned_urls,omitempty"`
	Usage        map[string]any `json:"usage,omitempty"`
}

// SeedanceModel identifies a Seedance model supported by the compatibility API.
type SeedanceModel string

const (
	SeedanceModelStandard SeedanceModel = "doubao-seedance-2-0-260128"
	SeedanceModelFast     SeedanceModel = "doubao-seedance-2-0-fast-260128"
)

type SeedanceVideoContentType string

const (
	SeedanceVideoContentText     SeedanceVideoContentType = "text"
	SeedanceVideoContentImageURL SeedanceVideoContentType = "image_url"
	SeedanceVideoContentVideoURL SeedanceVideoContentType = "video_url"
	SeedanceVideoContentAudioURL SeedanceVideoContentType = "audio_url"
)

type SeedanceVideoContentRole string

const (
	SeedanceVideoRoleFirstFrame     SeedanceVideoContentRole = "first_frame"
	SeedanceVideoRoleReferenceImage SeedanceVideoContentRole = "reference_image"
	SeedanceVideoRoleReferenceVideo SeedanceVideoContentRole = "reference_video"
	SeedanceVideoRoleReferenceAudio SeedanceVideoContentRole = "reference_audio"
)

type SeedanceVideoRatio string

const (
	SeedanceVideoRatio16x9     SeedanceVideoRatio = "16:9"
	SeedanceVideoRatio4x3      SeedanceVideoRatio = "4:3"
	SeedanceVideoRatio1x1      SeedanceVideoRatio = "1:1"
	SeedanceVideoRatio3x4      SeedanceVideoRatio = "3:4"
	SeedanceVideoRatio9x16     SeedanceVideoRatio = "9:16"
	SeedanceVideoRatio21x9     SeedanceVideoRatio = "21:9"
	SeedanceVideoRatioAdaptive SeedanceVideoRatio = "adaptive"
)

type SeedanceVideoResolution string

const (
	SeedanceVideoResolution480p  SeedanceVideoResolution = "480p"
	SeedanceVideoResolution720p  SeedanceVideoResolution = "720p"
	SeedanceVideoResolution1080p SeedanceVideoResolution = "1080p"
)

// SeedanceMediaURL is used for image, video, and audio content URLs.
type SeedanceMediaURL struct {
	URL string `json:"url"`
}

// SeedanceVideoContent is one compatibility API content item. Set Text for
// text items, or the matching URL field and Role for media items.
type SeedanceVideoContent struct {
	Type     SeedanceVideoContentType `json:"type"`
	Text     string                   `json:"text,omitempty"`
	ImageURL *SeedanceMediaURL        `json:"image_url,omitempty"`
	VideoURL *SeedanceMediaURL        `json:"video_url,omitempty"`
	AudioURL *SeedanceMediaURL        `json:"audio_url,omitempty"`
	Role     SeedanceVideoContentRole `json:"role,omitempty"`
}

// SeedanceVideoGenerationRequest is the strict request shape for the
// singular /v1/video/generations compatibility endpoint.
type SeedanceVideoGenerationRequest struct {
	Model                 SeedanceModel           `json:"model"`
	Content               []SeedanceVideoContent  `json:"content"`
	Ratio                 SeedanceVideoRatio      `json:"ratio,omitempty"`
	Resolution            SeedanceVideoResolution `json:"resolution,omitempty"`
	Duration              *int                    `json:"duration,omitempty"`
	GenerateAudio         *bool                   `json:"generate_audio,omitempty"`
	ReturnLastFrame       *bool                   `json:"return_last_frame,omitempty"`
	ExecutionExpiresAfter *int                    `json:"execution_expires_after,omitempty"`
	Seed                  *int64                  `json:"seed,omitempty"`
	Watermark             *bool                   `json:"watermark,omitempty"`
	SafetyIdentifier      string                  `json:"safety_identifier,omitempty"`
	CallbackURL           string                  `json:"callback_url,omitempty"`
}

type SeedanceVideoArtifact struct {
	URL string `json:"url"`
}

type SeedanceVideoUsage struct {
	PromptTokens           *int     `json:"prompt_tokens,omitempty"`
	CompletionTokens       *int     `json:"completion_tokens,omitempty"`
	TotalTokens            *int     `json:"total_tokens,omitempty"`
	Cost                   *float64 `json:"cost,omitempty"`
	Currency               string   `json:"currency,omitempty"`
	CostUSD                *float64 `json:"cost_usd,omitempty"`
	CostCNY                *float64 `json:"cost_cny,omitempty"`
	VideoGenerationCostCNY *float64 `json:"video_generation_cost_cny,omitempty"`
	SuperResolutionCostCNY *float64 `json:"sr_cost_cny,omitempty"`
}

type SeedanceVideoError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SeedanceVideoGenerationData struct {
	TaskID    string                  `json:"task_id"`
	Status    string                  `json:"status"`
	Artifacts []SeedanceVideoArtifact `json:"artifacts,omitempty"`
	Usage     *SeedanceVideoUsage     `json:"usage,omitempty"`
	Error     *SeedanceVideoError     `json:"error,omitempty"`
}

// SeedanceVideoGenerationCreateResponse is the flat response returned when a
// compatibility video generation is created.
type SeedanceVideoGenerationCreateResponse struct {
	ID        string        `json:"id"`
	TaskID    string        `json:"task_id"`
	Object    string        `json:"object"`
	Model     SeedanceModel `json:"model"`
	Status    string        `json:"status"`
	Progress  float64       `json:"progress"`
	CreatedAt string        `json:"created_at"`
}

// SeedanceVideoGenerationResponse is the envelope returned when querying a
// compatibility video generation.
type SeedanceVideoGenerationResponse struct {
	Code    string                      `json:"code"`
	Message string                      `json:"message"`
	Data    SeedanceVideoGenerationData `json:"data"`
}

type SeedanceAssetType string

const (
	SeedanceAssetTypeImage SeedanceAssetType = "Image"
	SeedanceAssetTypeVideo SeedanceAssetType = "Video"
	SeedanceAssetTypeAudio SeedanceAssetType = "Audio"
)

type SeedanceAssetGroupCreateRequest struct {
	Model       SeedanceModel `json:"model"`
	Name        string        `json:"Name"`
	Description string        `json:"Description,omitempty"`
}

type SeedanceAssetCreateRequest struct {
	Model     SeedanceModel     `json:"model"`
	URL       string            `json:"URL"`
	AssetType SeedanceAssetType `json:"AssetType"`
	GroupID   string            `json:"GroupId,omitempty"`
	Name      string            `json:"Name,omitempty"`
}

type SeedanceAssetGetRequest struct {
	Model SeedanceModel `json:"model"`
	ID    string        `json:"Id"`
}

// SeedanceAssetGroupCreateResponse is the flat response returned when an
// asset group is created.
type SeedanceAssetGroupCreateResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// SeedanceAssetCreateResponse is the flat response returned when an asset is
// created.
type SeedanceAssetCreateResponse struct {
	ID string `json:"id"`
}

// SeedanceAssetGetResponse is the flat response returned when an asset is
// queried.
type SeedanceAssetGetResponse struct {
	ID        string            `json:"id"`
	URL       string            `json:"url,omitempty"`
	Name      string            `json:"name,omitempty"`
	AssetType SeedanceAssetType `json:"asset_type"`
	GroupID   string            `json:"group_id,omitempty"`
	Status    string            `json:"status"`
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
