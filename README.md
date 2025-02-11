# Wolf Workouts

This project is inspired in the requirements of the book [Go With The Domain](https://threedots.tech/go-with-the-domain/) from Three Dots Labs to provide a workout scheduling system.

## Getting Started

// TODO

## Overall Design

- Microservice Architecture
- DDD Lite
- CQRS

## Technologies

- Postgres Database
- Docker
- gRPC and REST API
- Swagger for REST API
- React with Typescript for the frontend (UI)

## Requirements (Business Logic)

### Microservices

#### Trainings Application
1. **Purpose**: Handles the management of training sessions from a business perspective
- Manages training scheduling, cancellations, and rescheduling
- Handles user-training relationships
- Maintains training records and states

2. **Key Functionalities**:
- Schedule new trainings
- Cancel trainings
- Reschedule trainings
- View training history
- Manage training approvals
- Handle training requests

3. **Domain Focus**:
- Training entity management
- Training business rules
- User-training relationships
- Training state management

#### Trainer Application
1. **Purpose**: Manages trainer availability and schedule
- Handles trainer's calendar
- Manages available/unavailable hours
- Provides trainer schedule information

2. **Key Functionalities**:
- Make hours available/unavailable
- Check trainer availability
- Manage trainer's calendar
- Handle trainer schedule updates

3. **Domain Focus**:
- Trainer availability
- Schedule management
- Time slot management
- Calendar operations

#### Users

## Security
The microservices rely on a issue of a JWT (JSON Web Token) using the HMAC (Keyed-Hashing for Message Authentication). This is a symmetric algorithm which relies on a secret to encode and validate the issued token between the microservices.

A future improvement could be to use the RSA asymmetric encrypt.

// TODO