# Snippetbox

Snippetbox is a **production-style web application built in Go**, based on the Let’s Go book project and deliberately extended to reflect **real-world backend engineering practices**.

Rather than stopping at tutorial completion, this project was treated as a **deployable backend system**, with emphasis on:

* Clean architecture
* Testability
* Secure configuration
* CI/CD automation
* Production deployment concerns

It serves as a public backend engineering portfolio showcasing how I build Go services in real environments, since my primary commercial systems (gaming platforms and ERP systems) are private repositories.

---

## Project Overview

Snippetbox allows users to create, view, and manage text snippets via a web interface.

Core features include:

* User registration and authentication
* Secure session management
* Creation and retrieval of text snippets
* Server-side rendered HTML templates
* Persistent relational database storage
* Automated tests with mocks
* Fully automated CI/CD deployment

Live demo:
👉 https://snippetbox.high-la.dev
  
---

## Architecture & Design

The project follows idiomatic Go structure and a strict separation of concerns:

```
cmd/web/           Application entry point
internal/handlers  HTTP handlers & business logic
internal/models    Database models & interfaces
internal/mocks     Mock implementations for testing
ui/                HTML templates & static assets
migrations/        Database schema migrations
```


Key design principles:

* Dependency injection via interfaces
* Explicit error handling and structured logging
* Minimal global state
* Context-aware database operations
* Testability as a first-class concern

This structure mirrors how production Go services are commonly organized.

---

## Testing Strategy

Testing is a first-class concern in this project.

* Database access is abstracted behind interfaces
* Mock models are used to test handlers without a real DB
* HTTP handlers tested using net/http/httptest
* Tests validate routing, responses, status codes, and edge cases

This approach matches production backend testing practices, where isolation and reliability matter more than integration-only tests.

---

## Database

* MySQL (schema is PostgreSQL-compatible)
* Explicit schema with indexes
* Defensive SQL error handling
* Context-aware queries
* Migrations managed explicitly (no manual schema edits)

---

## Security

* Password hashing with bcrypt
* Secure cookie-based sessions
* CSRF protection
* Proper HTTP security headers
* Input validation and error sanitization
* Secrets stored outside the repository
---

## Tech Stack

* **Language:** Go
* **HTTP:** net/http
* **Templates:** Go HTML templates
* **Database:** MySQL
* **Testing:** testing, httptest, mocks
* **Auth & Security:** bcrypt, secure cookies
* **Deployment:** Docker, Nginx, GitHub Actions

---

## Production Setup & Deployment (Stages 1–8)

This project is deployed and maintained using a production-safe workflow, fully documented below.

**Stage 1 — Server & OS Setup**

* Linux VPS
* Non-root user
* Firewall enabled
* SSH key authentication
* Docker & Docker Compose installed

**Stage 2 — Reverse Proxy & HTTPS**

* Nginx as reverse proxy
* Let’s Encrypt TLS certificates
* HTTPS enforced
* Secure headers configured

**Stage 3 — Dockerized Application**

* Go app built as a static binary
* Docker image created for the application
* MySQL container managed via Docker Compose
* Internal Docker networking (no public DB exposure)

**Stage 4 — Database Initialization**

* On first deployment:
* MySQL container starts with empty volume
* Database and user created automatically
* App waits for DB readiness before starting

**Stage 5 — Database Migrations**

* Schema managed via versioned migration files
* Migrations run automatically during deployment
* Prevents manual schema drift
* Ensures reproducible environments (local → CI → prod)

**Stage 6 — Environment Variables & Secrets**

All configuration is externalized:

* .env file (never committed)
* Database credentials
* Session & CSRF keys
* Application environment flags (dev / prod)

Outcome:
* No secrets in GitHub
* Same image works across environments

**Stage 7 — CI/CD Pipeline (GitHub Actions)**

A full CI/CD pipeline is configured:

1. On git push to main
2. GitHub Actions:
* Builds Go binary
* Builds Docker image
* Runs tests

3. Deployment:
* Image pulled or updated on server
* Database migrations executed
* Only Snippetbox containers restarted

Outcome:
    git push → live production deployment

**Stage 8 — Zero-Downtime-Oriented Updates**

* Only application containers restart
* Database container remains untouched
* Persistent volumes preserved
* Safe incremental updates

## Screenshots

> Screenshots of the Snippetbox UI are included to demonstrate the working application and user flow.

* User signup & login
![Home Page](screenshots/signup.png)
![Home Page](screenshots/signin.png)
* Home page with recent snippets
![Home Page](screenshots/home.png)
* Snippet creation form
![Home Page](screenshots/create.png)
* Snippet detail view
![Home Page](screenshots/detail.png)


For an interactive experience, see the live demo: [snippetbox.high-la.dev](https://snippetbox.high-la.dev)

---

##  Author

**Haile Berhaneselassie**
Backend Engineer (Go)
🌐 [https://high-la.dev](https://high-la.dev)

---

## 📄 Notes for Reviewers

This repository is intended to be reviewed as a backend engineering sample:

Reviewers are encouraged to focus on:

* Project structure and package boundaries
* Test strategy and use of mocks
* Error handling and logging
* Configuration and secret management
* CI/CD and deployment automation
* The same patterns demonstrated here are used in my larger production systems.

Thank you for taking the time to review this project.
