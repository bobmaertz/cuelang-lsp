// Simple configuration example
package config

// Server configuration
server: {
	host: string | *"localhost"
	port: int & >=1 & <=65535 | *8080
	timeout: int | *30
}

// Database configuration
database: {
	driver: "postgres" | "mysql" | "sqlite"
	host:   string
	port:   int
	name:   string
}

// Example instance
myApp: {
	server: {
		host: "api.example.com"
		port: 443
	}
	database: {
		driver: "postgres"
		host:   "db.example.com"
		port:   5432
		name:   "myapp_db"
	}
}
