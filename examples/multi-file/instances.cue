// Example instances using shared types
package app

// Users
users: {
	admin: #User & {
		id:        1
		username:  "admin"
		email:     "admin@example.com"
		firstName: "System"
		lastName:  "Administrator"
		roles: ["admin", "user"]
	}

	john: #User & {
		id:        2
		username:  "john_doe"
		email:     "john@example.com"
		firstName: "John"
		lastName:  "Doe"
		roles: ["user"]
	}
}

// Services
services: {
	api: #Service & {
		name:    "user-api"
		version: "v1.2.3"
		replicas: 5
		resources: {
			cpu:    "1000m"
			memory: "2Gi"
		}
		endpoints: [
			{path: "/api/users", method:       "GET"},
			{path: "/api/users", method:       "POST"},
			{path: "/api/users/:id", method:   "GET"},
			{path: "/api/users/:id", method:   "PUT"},
			{path: "/api/users/:id", method:   "DELETE"},
			{path: "/api/health", method:      "GET", auth: false},
		]
	}

	auth: #Service & {
		name:    "auth-service"
		version: "v2.0.1"
		resources: {
			cpu:    "500m"
			memory: "1Gi"
		}
		endpoints: [
			{path: "/auth/login", method:  "POST", auth: false},
			{path: "/auth/logout", method: "POST"},
			{path: "/auth/refresh", method: "POST"},
		]
	}
}
