## ADDED Requirements

### Requirement: Fix bilingual language switching 404
The website SHALL support correct language switching on all pages without 404 errors.

#### Scenario: English to Chinese switch
- **WHEN** a user is on any English page and clicks the language switcher
- **THEN** they are redirected to the corresponding Chinese page (not a 404)

#### Scenario: Chinese to English switch
- **WHEN** a user is on any Chinese page and clicks the language switcher
- **THEN** they are redirected to the corresponding English page (not a 404)

#### Scenario: File structure matches VitePress locale config
- **WHEN** the site is built
- **THEN** English files are at docs/ root level, Chinese files are at docs/zh/
