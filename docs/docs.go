// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/projects": {
            "post": {
                "description": "Creates a new project with the provided data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Create a new project",
                "parameters": [
                    {
                        "description": "Project Data",
                        "name": "project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Project"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully created project",
                        "schema": {
                            "$ref": "#/definitions/models.Project"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/projects/student/{student_id}": {
            "get": {
                "description": "Fetches a project by its student ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Get a project by Student ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Student ID",
                        "name": "student_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved project",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.ProjectWithDetails"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid student ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Project not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/projects/{id}": {
            "get": {
                "description": "Fetches a project by its unique ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Get a project by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved project",
                        "schema": {
                            "$ref": "#/definitions/models.Project"
                        }
                    },
                    "400": {
                        "description": "Invalid project ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Project not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "put": {
                "description": "Updates a project by its ID with the provided data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Update an existing project",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated Project Data",
                        "name": "project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Project"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully updated project",
                        "schema": {
                            "$ref": "#/definitions/models.Project"
                        }
                    },
                    "400": {
                        "description": "Invalid project ID or request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes the specified project using its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Delete a project by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Project deleted successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid project ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/resource": {
            "post": {
                "description": "Upload a file to MinIO and save its information as a resource",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Resource"
                ],
                "summary": "Upload a file and create a resource",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "project_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Resource Type ID",
                        "name": "resource_type_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Resource Title",
                        "name": "title",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Resource"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/resource/project/{project_id}": {
            "get": {
                "description": "Get all resources associated with a project",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Resource"
                ],
                "summary": "Get resources by project ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "project_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Resource"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/resource/{id}": {
            "get": {
                "description": "Get a resource by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Resource"
                ],
                "summary": "Get a resource by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Resource ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Resource"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a resource and its file",
                "tags": [
                    "Resource"
                ],
                "summary": "Delete a resource",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Resource ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.ProjectWithDetails": {
            "type": "object",
            "properties": {
                "employees": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Employee"
                    }
                },
                "project": {
                    "$ref": "#/definitions/models.Project"
                },
                "students": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Student"
                    }
                }
            }
        },
        "models.Course": {
            "type": "object",
            "properties": {
                "course_name": {
                    "type": "string"
                },
                "course_no": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "models.Employee": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "employee_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "role": {
                    "$ref": "#/definitions/models.Role"
                },
                "role_id": {
                    "type": "integer"
                }
            }
        },
        "models.Major": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "major_name": {
                    "type": "string"
                }
            }
        },
        "models.Project": {
            "type": "object",
            "properties": {
                "abstract": {
                    "type": "string"
                },
                "academic_year": {
                    "type": "integer"
                },
                "advisor": {
                    "$ref": "#/definitions/models.Employee"
                },
                "advisor_id": {
                    "type": "integer"
                },
                "committees": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Employee"
                    }
                },
                "course": {
                    "$ref": "#/definitions/models.Course"
                },
                "course_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "major": {
                    "$ref": "#/definitions/models.Major"
                },
                "major_id": {
                    "type": "integer"
                },
                "members": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Student"
                    }
                },
                "old_project_no": {
                    "type": "string"
                },
                "project_no": {
                    "type": "string"
                },
                "project_status": {
                    "type": "string"
                },
                "relation_description": {
                    "type": "string"
                },
                "section": {
                    "$ref": "#/definitions/models.Section"
                },
                "section_id": {
                    "type": "integer"
                },
                "semester": {
                    "type": "integer"
                },
                "title_en": {
                    "type": "string"
                },
                "title_th": {
                    "type": "string"
                }
            }
        },
        "models.Resource": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "project": {
                    "$ref": "#/definitions/models.Project"
                },
                "project_id": {
                    "type": "integer"
                },
                "resource_type": {
                    "$ref": "#/definitions/models.ResourceType"
                },
                "resource_type_id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "models.ResourceType": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "resource_type": {
                    "type": "string"
                }
            }
        },
        "models.Role": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "role_name": {
                    "type": "string"
                }
            }
        },
        "models.Section": {
            "type": "object",
            "properties": {
                "course": {
                    "$ref": "#/definitions/models.Course"
                },
                "course_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "section_number": {
                    "type": "string"
                },
                "semester": {
                    "type": "integer"
                }
            }
        },
        "models.Student": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "major": {
                    "$ref": "#/definitions/models.Major"
                },
                "major_id": {
                    "type": "integer"
                },
                "student_name": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
