## ADDED Requirements

### Requirement: VitePress site configuration
The project SHALL use VitePress as the documentation site framework.

#### Scenario: Site builds successfully
- **WHEN** `npx vitepress build docs` is run
- **THEN** the site is built to `docs/.vitepress/dist/` with no errors

#### Scenario: Dev server runs
- **WHEN** `npx vitepress dev docs` is run
- **THEN** the site is accessible at localhost:5173

### Requirement: Multi-language support
The site SHALL support English and Chinese languages with a language switcher.

#### Scenario: Language switcher
- **WHEN** a user visits the site
- **THEN** a language switcher is visible in the navbar allowing switching between English and Chinese

#### Scenario: Default language
- **WHEN** a user visits the root URL
- **THEN** the English version is shown by default

#### Scenario: Chinese version
- **WHEN** a user switches to Chinese or visits /zh/
- **THEN** the Chinese version is shown with all content translated

### Requirement: Navigation structure
The site SHALL have a clear navigation with sidebar and navbar.

#### Scenario: Navbar navigation
- **WHEN** a user views the navbar
- **THEN** it contains: Guide, Architecture, Harness links in both languages

#### Scenario: Sidebar navigation
- **WHEN** a user navigates to a documentation page
- **THEN** a sidebar shows the document hierarchy for the current section

### Requirement: GitHub Pages deployment
The site SHALL deploy automatically via GitHub Actions.

#### Scenario: Auto deployment
- **WHEN** changes are pushed to main that affect docs/
- **THEN** the VitePress site is built and deployed to GitHub Pages
