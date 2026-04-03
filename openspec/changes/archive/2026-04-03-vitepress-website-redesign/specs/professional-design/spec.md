## ADDED Requirements

### Requirement: Dark/light theme
The site SHALL support dark and light theme with a toggle button.

#### Scenario: Theme toggle
- **WHEN** a user clicks the theme toggle in the navbar
- **THEN** the site switches between dark and light themes

#### Scenario: Default theme
- **WHEN** a user first visits the site
- **THEN** the theme matches their system preference (prefers-color-scheme)

### Requirement: Professional code blocks
Code blocks SHALL use syntax highlighting and have a copy button.

#### Scenario: Code highlighting
- **WHEN** a user views a code block
- **THEN** Go code is highlighted with proper Go syntax coloring

#### Scenario: Copy button
- **WHEN** a user hovers over a code block
- **THEN** a copy button appears that copies the code to clipboard

### Requirement: Responsive design
The site SHALL be fully responsive on mobile, tablet, and desktop.

#### Scenario: Mobile view
- **WHEN** a user views the site on a mobile device
- **THEN** all content is readable, navigation collapses to a hamburger menu

#### Scenario: Desktop view
- **WHEN** a user views the site on a desktop
- **THEN** the full sidebar navigation and feature card grid are displayed
