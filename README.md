# Wolf Workouts

Wolf Workouts is fitness platform that connects trainers with their clients. It's developed in Go (Golang). It manages training sessions, user accounts, and scheduling through a clean microservices architecture. The platform makes it easy for trainers to organize their sessions and for clients to book and track their workouts.

The system is built using modern Go practices, featuring Domain-Driven Design, CQRS for data operations, JWT authentication, and PostgreSQL for storage. All functionality is exposed through REST APIs (and some using gRPC) for easy integration. 

The goal is to be a example of microservices architecture using modern technologies for cloud native environment.

This project is inspired in the requirements of the book [Go With The Domain](https://threedots.tech/go-with-the-domain/) from [Three Dots Labs](https://threedots.tech/).

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.23.5 or higher
- Make (optional, for using Makefile commands)

### Setup
1. Clone the repository
```bash
git clone https://github.com/davidalecrim1/wolf-workouts.git
cd wolf-workouts
```

2. Make sure the environment variables are set in `.env`.

3. Run the containers:
```bash
docker compose up -d
```

## Architecture

### Overall Design
- **Microservice Architecture**: Separate services for users and trainings
- **DDD Lite**: Domain-Driven Design principles applied in a lightweight manner
- **CQRS**: Command Query Responsibility Segregation pattern
- **Clean Architecture**: Separation of concerns with clear boundaries

### Technologies
- **Backend**:
  - Go (Golang) for microservices
  - PostgreSQL for data persistence
  - Docker for containerization
  - gRPC and REST APIs for communication
  - JWT for authentication
  - Swagger for API documentation
- **Frontend**:
  - React with TypeScript
  - Modern UI components

## Microservices

### Users Service
1. **Purpose**: Manages user accounts and authentication
- User registration and login
- Profile management
- Authentication and authorization

2. **Key Features**:
- JWT-based authentication
- Password hashing
- User profile CRUD operations

### Trainings Service
1. **Purpose**: Handles training session management
- Training scheduling and management
- User-training relationships
- Training history tracking

2. **Key Features**:
- Schedule new trainings
- Cancel/reschedule sessions
- View training history
- Training approval workflow

### Trainer Service
1. **Purpose**: Manages trainer availability and scheduling
- Trainer calendar management
- Availability settings
- Schedule coordination

2. **Key Features**:
- Manage available hours
- Handle trainer schedules
- Time slot management
- Calendar operations

## Security
The microservices implement JWT (JSON Web Token) authentication using HMAC (Keyed-Hashing for Message Authentication). This symmetric algorithm relies on a secret key to encode and validate tokens between services.

Future security improvements:
- Implement RSA asymmetric encryption

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```