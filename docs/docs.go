// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-09-26 00:41:12.040029362 +0200 CEST m=+0.178136750

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/agenda": {
            "post": {
                "description": "Get user agenda by username or e-mail",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Old"
                ],
                "summary": "Get user agenda",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username or e-mail",
                        "name": "mail",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of events",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ade.OldFormat"
                            }
                        }
                    }
                }
            }
        },
        "/agenda/{mail}": {
            "get": {
                "description": "Get user agenda by username or e-mail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Old"
                ],
                "summary": "Get user agenda",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username or e-mail",
                        "name": "mail",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of events",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ade.OldFormat"
                            }
                        }
                    }
                }
            }
        },
        "/api/ade-esiee/agenda": {
            "post": {
                "description": "Get user agenda by username or e-mail",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Old"
                ],
                "summary": "Get user agenda",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username or e-mail",
                        "name": "mail",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of events",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ade.OldFormat"
                            }
                        }
                    }
                }
            }
        },
        "/api/ade-esiee/agenda/{mail}": {
            "get": {
                "description": "Get user agenda by username or e-mail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Old"
                ],
                "summary": "Get user agenda",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username or e-mail",
                        "name": "mail",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of events",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ade.OldFormat"
                            }
                        }
                    }
                }
            }
        },
        "/api/ics/{mail}": {
            "get": {
                "description": "Get user agenda by username or e-mail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V2",
                    "Agenda",
                    "ICS"
                ],
                "summary": "Get user agenda in ICS format",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username or e-mail",
                        "name": "mail",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ICS Calendar",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/rooms/{hour}": {
            "get": {
                "description": "Get all the free rooms at now or now + X hours",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rooms"
                ],
                "summary": "Get free rooms",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Hour shift",
                        "name": "hour",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Array of free rooms",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Do ping to check api",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Core"
                ],
                "summary": "Ask for ping get pong",
                "responses": {
                    "200": {
                        "description": "{\"message\": \"pong\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "description": "Got API informations about local files and uptime",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Core"
                ],
                "summary": "Get API status",
                "responses": {
                    "200": {
                        "description": "API informations",
                        "schema": {
                            "$ref": "#/definitions/string"
                        }
                    }
                }
            }
        },
        "/v2/agenda/{mail}": {
            "get": {
                "description": "Get user agenda by username or e-mail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V2",
                    "Agenda"
                ],
                "summary": "Get user agenda",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username or e-mail",
                        "name": "mail",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of events",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ade.Event"
                            }
                        }
                    }
                }
            }
        },
        "/v2/events/{name}": {
            "get": {
                "description": "Get all events of specific unite with its code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V2",
                    "Agenda"
                ],
                "summary": "Get events of specific unite",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Unite Code",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of events",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ade.Event"
                            }
                        }
                    }
                }
            }
        },
        "/v2/groups/{mail}": {
            "get": {
                "description": "Get user groups by username or e-mail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V2",
                    "Aurion"
                ],
                "summary": "Get user groups",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username or e-mail",
                        "name": "mail",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of groups",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/group.Group"
                            }
                        }
                    }
                }
            }
        },
        "/v2/unite/{name}": {
            "get": {
                "description": "Get unite code and label",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V2",
                    "Aurion"
                ],
                "summary": "Get unite information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Unite Code",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Unite informations",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/unite.Unite"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ade.Event": {
            "type": "object",
            "properties": {
                "color": {
                    "description": "color r,g,b",
                    "type": "string"
                },
                "creation": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "duration": {
                    "description": "hour quarters count duration",
                    "type": "integer"
                },
                "endHour": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "info": {
                    "description": "info",
                    "type": "string"
                },
                "lastUpdate": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "resources": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ade.Resource"
                    }
                },
                "startHour": {
                    "type": "string"
                },
                "xmlname": {
                    "type": "string"
                }
            }
        },
        "ade.OldFormat": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "end": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "prof": {
                    "type": "string"
                },
                "rooms": {
                    "type": "string"
                },
                "start": {
                    "type": "string"
                },
                "unite": {
                    "type": "string"
                }
            }
        },
        "ade.Resource": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "xmlname": {
                    "type": "string"
                }
            }
        },
        "group.Group": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "groups": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "unite": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "unite.Unite": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "label": {
                    "type": "string"
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
	Version:     "0.7.0",
	Host:        "ade.wallforfry.fr",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "ESIEE API",
	Description: "API pour ade et aurion",
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
	swag.Register(swag.Name, &s{})
}
