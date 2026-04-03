## ADDED Requirements

### Requirement: Blocking API call
The system SHALL support blocking (non-streaming) API calls to the Anthropic Messages API.

#### Scenario: Successful blocking call
- **WHEN** a message is sent with stream=false
- **THEN** the full response is returned after the API completes

### Requirement: Streaming API call
The system SHALL support streaming API calls with real-time text delta delivery.

#### Scenario: Streaming text response
- **WHEN** a message is sent with stream=true
- **THEN** text_delta events are delivered to the onTextDelta callback in real-time as they arrive

#### Scenario: Streaming tool_use response
- **WHEN** the model returns tool_use blocks via streaming
- **THEN** input_json_delta events are accumulated and the complete tool_use block is assembled

### Requirement: SSE event parsing
The system SHALL parse SSE (Server-Sent Events) format from the API response stream.

#### Scenario: Parse message_start event
- **WHEN** an SSE event with type "message_start" is received
- **THEN** the message metadata (id, model, usage) is extracted

#### Scenario: Parse content_block_delta event
- **WHEN** an SSE event with type "content_block_delta" is received
- **THEN** the delta (text_delta or input_json_delta) is extracted and processed

#### Scenario: Parse message_delta event
- **WHEN** an SSE event with type "message_delta" is received
- **THEN** the stop_reason and final usage are extracted

#### Scenario: Handle partial chunks
- **WHEN** an SSE event spans multiple network chunks
- **THEN** the parser buffers incomplete events and assembles them correctly

### Requirement: HTTP error handling
The system SHALL handle HTTP errors gracefully.

#### Scenario: Invalid API key (401)
- **WHEN** the API returns HTTP 401
- **THEN** a clear "API Key is invalid" error is returned

#### Scenario: Rate limited (429)
- **WHEN** the API returns HTTP 429
- **THEN** the request is retried with exponential backoff, up to 3 retries

#### Scenario: Server error (5xx)
- **WHEN** the API returns HTTP 500/502/503
- **THEN** a descriptive server error message is returned

### Requirement: Request headers
The system SHALL set required headers on all API requests.

#### Scenario: Required headers
- **WHEN** an API request is made
- **THEN** x-api-key, anthropic-version (2023-06-01), and content-type headers are included

### Requirement: API request construction
The system SHALL construct valid API requests with all required fields.

#### Scenario: Full request with tools
- **WHEN** a request includes tools
- **THEN** each tool's name, description, and input_schema are serialized correctly
