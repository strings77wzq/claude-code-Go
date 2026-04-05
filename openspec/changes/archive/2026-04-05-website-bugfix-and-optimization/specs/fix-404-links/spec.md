## MODIFIED Requirements

### Requirement: Fix all 404 and wrong-language links
All links in the Chinese navigation and sidebar SHALL point to Chinese pages.

#### Scenario: Chinese sidebar links
- **WHEN** a user clicks any link in the Chinese sidebar
- **THEN** they are taken to the Chinese version of that page (not 404, not English)

#### Scenario: Chinese navbar links
- **WHEN** a user clicks any link in the Chinese navbar dropdowns
- **THEN** they are taken to the Chinese version of that page

#### Scenario: Missing Chinese pages
- **WHEN** a Chinese page is linked but doesn't exist
- **THEN** the Chinese page is created with translated content
