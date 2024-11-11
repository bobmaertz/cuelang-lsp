package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Metamodel is the top level specification for LSP.
type MetaModel struct {
	// Metadata contains the vrsions information about the document.
	Metadata struct {
		Version string `json:"version"`
	} `json:"metaData"`
	// Requests defined the request parameters
	Request []Request `json:"requests"`
	// Structures handle the models
	Structures []Structures `json:"structures"`
	// Notifications handle the async notifications from the LSP
	Notifications []Request `json:"notifications"`
	// Enumerations <TODO>
	Enumerations []interface{} `json:"enumerations"`
	// TypeAliases <TODO>
	TypeAliases []interface{} `json:"typeAliases"`
}

type Structures struct {
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

type Option struct {
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
}

type Properties struct {
	Name          string `json:"name,omitempty"`
	Type          Type   `json:"type,omitempty"`
	Optional      *bool  `json:"optional,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

type Type struct {
	Name    string      `json:"name,omitempty"`
	Kind    string      `json:"kind"`
	Element *Option     `json:"element,omitempty"`
	Items   []Option    `json:"items,omitempty"`
	Value   interface{} `json:"value,omitempty"` // Or Element or string ....
	Since   string      `json:"single,omitempty"`
}

func main() {
    //TOOD: Read in from stdin ..
	b, err := os.ReadFile("testdata/metaModel.json")
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(10)
	}


    //Unmarshall into go represetation 
	model := &MetaModel{}
	err = json.Unmarshal(b, model)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(15)
	}

    //TODO: Specify package. 
    //Create output grammer file
	file, err := os.OpenFile("test.go", os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(16)
	}
	defer file.Close()

	fileWriter := bufio.NewWriter(file)

    //TODO: Specify package 
	fmt.Fprint(fileWriter,  "package main\n")


    // For each structure.. 
	for _, s := range model.Structures {

        //TODO: This is a good spot to use a template instead.. 
		start := "type %s struct {\n"
		end := "}\n"
		fmt.Fprintf(fileWriter,  start, s.Name)

		// TODO: each field
		for _, p := range s.Properties {
			if p.Documentation != "" {
                doc := strings.ReplaceAll(p.Documentation, "\n", "")
				fmt.Fprintf(fileWriter, "\t // %s %s\n", ToTitleCase(p.Name), doc)
			}
			fmt.Fprintf(fileWriter, "\t %s %s\n", ToTitleCase(p.Name), p.Type.Name)
		}

		fmt.Fprint(fileWriter, end)
	}

	// o, _ := json.MarshalIndent(model, "", "    ")
	// fmt.Print(string(o))
}



func ToTitleCase(s string) string {
    if s == ""{
        return s 
    }
    return strings.ToUpper(string(s[0]))+ s[1:]
}
