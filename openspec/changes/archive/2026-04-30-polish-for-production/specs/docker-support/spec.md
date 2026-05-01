## ADDED Requirements

### Requirement: Docker image available
The project SHALL provide a Docker image for deployment.

#### Scenario: Dockerfile exists
- **WHEN** a developer lists project root
- **THEN** they see a `Dockerfile`

#### Scenario: Docker build works
- **WHEN** a developer runs `docker build -t go-code .`
- **THEN** the image builds successfully

#### Scenario: Docker run works
- **WHEN** a developer runs `docker run -it go-code`
- **THEN** the CLI starts inside the container

### Requirement: Docker Compose available
The project SHALL provide docker-compose.yml for easy deployment.

#### Scenario: docker-compose.yml exists
- **WHEN** a developer lists project root
- **THEN** they see a `docker-compose.yml`

#### Scenario: Compose up works
- **WHEN** a developer runs `docker-compose up`
- **THEN** the service starts with API key from env
