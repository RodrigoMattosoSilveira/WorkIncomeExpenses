# Introduction
Build a `server-rendered HTMX` app as a `BFF`, keep `domain CRUD` behind `JSON/gRPC` services.

```bash
BFF Web Handlers (Adapters)
        ↓
PeopleAPI (Port)
        ↓
HTTP PeopleClient (Adapter)
        ↓
people-svc

```

### Use a “UI/BFF” service as the HTMX + template server
One service renders all HTML (full pages + partials). It talks to backend services over HTTP.
Why this is the best fit:
- the HTMX requests will serve HTML; the microservices will return JSON.
- The UI/BFF will keep a consistent UX, shared layout, form validation rendering, and HTMX behaviors in one place.

**Rule of thumb**
- Backend services: return JSON (or gRPC) + domain rules.
- UI/BFF: returns HTML (full page or partial) and orchestrates calls.

### Design the CRUD UI as “pages + fragments”
With HTMX, each CRUD operation will have:
- A full-page route (normal navigation)
- A fragment route (HTMX swap target)

#### Example resource
Below is a simple example of how to apply this approach to the `person` domain.

##### Full pages:
- GET /person (index page)
- GET /person/new
- GET /person/:id (details)
- GET /person/:id/edit

##### HTMX fragment endpoints 
Often same handlers, different template:
- GET /person returns page or just <tbody> depending on HX-Request
- POST /person returns new row fragment (or validation fragment)
- PATCH /person/:id returns updated row fragment
- DELETE /person/:id returns empty fragment or a flash message fragment

##### Implementation trick
In handlers, detect HTMX via header:
- HX-Request: true
  and then choose index.html vs person_tbody.html.

### Use Post/Redirect/Get for non-HTMX, and direct fragment returns for HTMX
This keeps progressive enhancement: if HTMX is off, the app still works:
- Normal browser POST: validate → store → redirect to GET /person
- HTMX POST: validate → store → return fragment to swap in-place

### Put validation & domain rules in the backend service, render errors in HTML in the BFF
This avoids duplicating business rules in UI. Flow:
1. UI/BFF receives form post.
2. UI/BFF calls backend service CreatePerson(...).
3. Backend returns either:
	- Success DTO
	- Structured validation errors (field → message)
4. UI/BFF renders:
	- success fragment (new row / updated row)
	- or the form fragment with inline errors

### Data ownership
Database-per-service
- Each service owns its DB schema.
- UI/BFF never queries DBs directly.
- Cross-service views: compose in UI/BFF.

### CRUD + microservices
Distributed transactions are where CRUD microservices hurt.
Practical patterns:
- Keep all CRUD within a single service boundary.
- UI/BFF coordinates cross service calls.
- Create in one service, reference by ID in another.
- Ohter possible across services coordination, to be considered in the distant future:
	- Use outbox + events for eventual consistency
	- Use sagas only when truly needed

### Go project structure
```bash
WIE
|--- cmd
     |--- ui_bff (HTML + HTMX server)

	 |--- person_svc (person JSON API)

|--- internal
     |--- config (.env file handler)

     |--- platform/ (logging, config, tracing, db, http middleware)

	 |--- entities
	      |-- person (routes, controller, service, model, repository)

     |--- bff (web, clients, templates, htmx helpers)
```

## Implementation
Below is an outline of the implmentation of a `Person` service using Fiber with the UI/BFF pattern (server-rendered html/template + HTMX), talking to a separate person-svc (JSON).

### Summary
Architecture in one picture: HTMX browser → bff (HTML) → person-svc (JSON) → bff renders partials

- bff (HTML): owns templates + HTMX endpoints, returns pages/fragments
- person-svc (API): owns Person domain + DB, returns JSON + validation errors

### Routes
#### BFF (HTML)
- GET  /person (page or <tbody> fragment)
- GET  /person/new (form page or fragment)
- POST /person (creates; returns new <tr> fragment or form w/errors)
- GET  /person/:id/edit (row editor fragment)
- PATCH /person/:id (updates; returns updated <tr> fragment or editor w/errors)
- DELETE /person/:id (deletes; returns empty)

#### “Person” service (JSON)

- GET    /api/person
- POST   /api/person
- GET    /api/person/:id
- PATCH  /api/person/:id
- DELETE /api/person/:id

### 
Folder layout

```bash
WIE
/cmd
  /bff
    main.go
  /person-svc
    main.go

/internal
  /bff
    /web
      handlers_person.go
      render.go
      htmx.go
    /clients
      person_client.go
    /views
      layout.html
      person_index.html
      person_tbody.html
      person_form.html
      person_row.html
      person_row_edit.html

  /person
    domain.go
    service.go
    repo.go
    api_handlers.go
```
