## ADDED Requirements

### Requirement: API client response body cleanup
The API client SHALL close response bodies on all error paths.

#### Scenario: HTTP error response
- **WHEN** an HTTP error occurs
- **THEN** the response body is properly closed before returning the error
