# auto-recovery Specification

## Purpose
TBD - created by archiving change short-mid-term-roadmap. Update Purpose after archive.
## Requirements
### Requirement: Auto-recovery
The system SHALL automatically recover from agent crashes.

#### Scenario: Recovery recipes
- **WHEN** the agent encounters an error
- **THEN** a recovery recipe is applied based on error type

#### Scenario: Error-specific recovery
- **WHEN** API timeout occurs
- **THEN** retry with exponential backoff
- **WHEN** rate limit occurs
- **THEN** wait and retry
