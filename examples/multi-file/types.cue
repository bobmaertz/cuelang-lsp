// Shared type definitions
package app

// User defines a user entity
#User: {
	id:        int & >0
	username:  string & =~"^[a-zA-Z0-9_]+$"
	email:     string & =~"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	firstName: string
	lastName:  string
	active:    bool | *true
	roles:     [...string]
}

// Service defines a microservice
#Service: {
	name:    string
	version: string & =~"^v[0-9]+\\.[0-9]+\\.[0-9]+$"
	replicas: int & >=1 & <=100 | *3
	resources: {
		cpu:    string
		memory: string
	}
	endpoints: [...#Endpoint]
}

// Endpoint defines an API endpoint
#Endpoint: {
	path:   string
	method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH"
	auth:   bool | *true
}
