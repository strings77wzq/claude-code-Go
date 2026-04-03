## ADDED Requirements

### Requirement: Branded SVG icons
The website SHALL use custom SVG icons instead of emoji for feature cards.

#### Scenario: Feature card icons
- **WHEN** a user views the homepage feature cards
- **THEN** each card uses a custom SVG icon in Go blue (#00ADD8) color scheme

#### Scenario: Icon design consistency
- **WHEN** all icons are viewed together
- **THEN** they share a consistent style: stroke-based, rounded corners, Go community aesthetic

#### Scenario: Icon subjects
- **WHEN** the icons are examined individually
- **THEN** they represent: Agent Loop (gopher + cycle arrows), Tools (wrench + code bracket), Permission (shield), MCP (plug/connection), SSE (wave/stream), Context (brain/chip)

### Requirement: Homepage stats section
The homepage SHALL display project core numbers.

#### Scenario: Stats display
- **WHEN** a user scrolls the homepage
- **THEN** they see: source files count, modules count, built-in tools count, LOC count

### Requirement: Homepage learning outcomes section
The homepage SHALL show what users can learn from the project.

#### Scenario: Learning outcomes
- **WHEN** a user views the homepage
- **THEN** they see 5 learning outcome items: AI Agent architecture, streaming API, tool system, security design, engineering practice
