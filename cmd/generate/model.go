package main

// This file contains the structures that are defined by the VSCode metamodel for LSP.

// Metamodel is the top level specification for LSP.
type MetaModel struct {
	// Metadata contains the vrsions information about the document.
	Metadata Metadata `json:"metaData"`
	// Requests defined the request parameters
	Request []Request `json:"requests"`
	// Structures handle the models
	Structures []Structure `json:"structures"`
	// Notifications handle the async notifications from the LSP
	Notifications []Request `json:"notifications"`
	// Enumerations <TODO>
	Enumerations []Enumeration `json:"enumerations"`
	// TypeAliases <TODO>
	TypeAliases []Type `json:"typeAliases"`
}
type Metadata struct {
	Version string `json:"version"`
}

type Enumeration struct {
	Name   string  `json:"name"`
	Type   Type    `json:"type"`
	Values []Value `json:"values"`
}

type Value struct {
	Name          string      `json:"name"`
	Value         interface{} `json:"value"`
	Documentation string      `json:"documentation"`
}

type Option struct {
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
}

type Structure struct {
	Name       string       `json:"name"`
	Properties []Properties `json:"properties"`
	Kind       []Option     `json:"kind,omitempty"`
	Mixins     []Option     `json:"mixins,omitempty"`
	Extends    []Option     `json:"extends,omitempty"`
}

type Request struct {
	Method              string `json:"method"`
	TypeName            string `json:"typeName"`
	Type                Result `json:"type,omitempty"`
	Result              Result `json:"result,omitempty"`
	MessageDirection    string `json:"messageDirection"`
	Params              Option `json:"params,omitempty"` // Notif
	PartialResult       Result `json:"partialResult,omitempty"`
	RegistrationOptions Option `json:"registrationOptions,omitempty"`
	Documentation       string `json:"documentation"`
}

type Notification struct {
	Method              string `json:"method"`
	TypeName            string `json:"typeName,omitempty"`
	MessageDirection    string `json:"messageDirection"`
	Params              Option `json:"params,omitempty"`
	RegistrationMethod  string `json:"registrationMethod,omitempty"`
	RegistrationOptions Option `json:"registrationOptions,omitempty"`
	Documentation       string `json:"documentation"`
	Since               string `json:"since,omitempty"`
}

type Result struct {
	Kind  string    `json:"kind,omitempty"`
	Name  string    `json:"name,omitempty"`
	Items []Element `json:"items,omitempty"`
}

type Element struct {
	Option
	Element []Option `json:"items,omitempty"`
}

type Properties struct {
	Name          string `json:"name,omitempty"`
	Type          Type   `json:"type,omitempty"`
	Optional      *bool  `json:"optional,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

type Type struct {
	Name          string      `json:"name,omitempty"`
	Kind          string      `json:"kind"`
	Element       *Option     `json:"element,omitempty"`
	Items         []Option    `json:"items,omitempty"`
	Value         interface{} `json:"value,omitempty"` // Or Element or string ....
	Since         string      `json:"single,omitempty"`
	Documentation string      `json:"documentation,omitempty"`
}
