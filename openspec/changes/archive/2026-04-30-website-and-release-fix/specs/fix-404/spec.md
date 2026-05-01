## MODIFIED Requirements

### Requirement: No 404 links
All navigation and sidebar links SHALL point to existing pages.

#### Scenario: Chinese sidebar matches English
- **WHEN** a user views the Chinese sidebar
- **THEN** every link points to an existing Chinese page

#### Scenario: Language switching works
- **WHEN** a user switches between English and Chinese
- **THEN** they are taken to the corresponding page in the other language (not 404)
