## ADDED Requirements

### Requirement: Detailed error classification
The system SHALL classify API errors and show user-friendly messages.

#### Scenario: Invalid API key (401)
- **WHEN** the API returns 401
- **THEN** the message is: "Invalid API key. Please check your ANTHROPIC_API_KEY."

#### Scenario: Access denied (403)
- **WHEN** the API returns 403
- **THEN** the message is: "API access denied. Check your API key permissions."

#### Scenario: Rate limited (429)
- **WHEN** the API returns 429
- **THEN** the message is: "Rate limited. Retrying automatically..."

#### Scenario: Server error (500+)
- **WHEN** the API returns 5xx
- **THEN** the message is: "Server error. Please try again later."

#### Scenario: Network error
- **WHEN** a network error occurs
- **THEN** the message is: "Network error. Please check your internet connection."
