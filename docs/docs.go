// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/book": {
            "get": {
                "description": "This method return books",
                "consumes": [
                    "text/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book Operations"
                ],
                "summary": "Get books",
                "parameters": [
                    {
                        "type": "string",
                        "description": "BookName",
                        "name": "bookname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "GenreName",
                        "name": "genre",
                        "in": "query"
                    }
                ]
            },
            "put": {
                "description": "Update book",
                "consumes": [
                    "text/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book Operations"
                ],
                "summary": "Update book",
                "parameters": [
                    {
                        "description": "Book",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Book"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "int"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a book with data",
                "consumes": [
                    "text/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book Operations"
                ],
                "summary": "Create a book",
                "parameters": [
                    {
                        "description": "Book",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.BookRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "int"
                        }
                    }
                }
            }
        },
        "/book/{id}": {
            "get": {
                "description": "This method gets book via id",
                "consumes": [
                    "text/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book Operations"
                ],
                "summary": "Get book by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Book Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ]
            },
            "delete": {
                "description": "This method delete book by id",
                "consumes": [
                    "text/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book Operations"
                ],
                "summary": "Delete book by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Book Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "int"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Book": {
            "type": "object",
            "required": [
                "amount",
                "genre",
                "name",
                "price"
            ],
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "genre": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                }
            }
        },
        "models.BookRequest": {
            "type": "object",
            "required": [
                "amount",
                "genre",
                "name",
                "price"
            ],
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "genre": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:4000",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Book API",
	Description: "This is a sample server.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
