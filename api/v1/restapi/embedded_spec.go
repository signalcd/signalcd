// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "https",
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Swagger 2.0 specification for SignalCD",
    "title": "SignalCD Swagger Spec",
    "termsOfService": "http://swagger.io/terms/",
    "version": "v0.0.0"
  },
  "host": "localhost:6660",
  "basePath": "/api/v1",
  "paths": {
    "/deployments": {
      "get": {
        "tags": [
          "deployments"
        ],
        "summary": "Returns the history of deployments",
        "operationId": "deployments",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/deployment"
              }
            }
          },
          "500": {
            "description": "internal server error"
          }
        }
      }
    },
    "/deployments/current": {
      "get": {
        "tags": [
          "deployments"
        ],
        "summary": "Returns the currently active deployment",
        "operationId": "currentDeployment",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/deployment"
            }
          }
        }
      },
      "post": {
        "tags": [
          "deployments"
        ],
        "summary": "Schedule a new deployment",
        "operationId": "setCurrentDeployment",
        "parameters": [
          {
            "type": "string",
            "name": "pipeline",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/deployment"
            }
          },
          "500": {
            "description": "internal server error"
          }
        }
      }
    },
    "/pipelines": {
      "get": {
        "tags": [
          "pipeline"
        ],
        "summary": "returns a list of all pipelines",
        "operationId": "pipelines",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/pipeline"
              }
            }
          },
          "400": {
            "description": "bad request"
          },
          "500": {
            "description": "internal server error"
          }
        }
      },
      "post": {
        "tags": [
          "pipeline"
        ],
        "summary": "creates a new pipeline",
        "operationId": "create",
        "parameters": [
          {
            "description": "The pipeline to be created",
            "name": "pipeline",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/pipeline"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/pipeline"
            }
          },
          "500": {
            "description": "internal server error"
          }
        }
      }
    },
    "/pipelines/{id}": {
      "get": {
        "tags": [
          "pipeline"
        ],
        "summary": "returns a pipeline by id",
        "operationId": "pipeline",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/pipeline"
            }
          },
          "400": {
            "description": "bad request"
          },
          "500": {
            "description": "internal server error"
          }
        }
      }
    }
  },
  "definitions": {
    "check": {
      "required": [
        "name",
        "image"
      ],
      "properties": {
        "duration": {
          "type": "number"
        },
        "environment": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "key": {
                "type": "string"
              },
              "value": {
                "type": "string"
              }
            }
          }
        },
        "image": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      }
    },
    "deployment": {
      "required": [
        "number"
      ],
      "properties": {
        "created": {
          "type": "string",
          "format": "date-time"
        },
        "finished": {
          "type": "string",
          "format": "date-time"
        },
        "number": {
          "type": "number",
          "format": "int64"
        },
        "pipeline": {
          "$ref": "#/definitions/pipeline"
        },
        "started": {
          "type": "string",
          "format": "date-time"
        },
        "status": {
          "$ref": "#/definitions/deploymentstatus"
        }
      }
    },
    "deploymentstatus": {
      "properties": {
        "phase": {
          "type": "string",
          "enum": [
            "unknown",
            "success",
            "failure",
            "progress"
          ]
        }
      }
    },
    "pipeline": {
      "properties": {
        "checks": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/check"
          }
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "name": {
          "type": "string"
        },
        "steps": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/step"
          }
        }
      }
    },
    "step": {
      "required": [
        "name",
        "image"
      ],
      "properties": {
        "commands": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "image": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      }
    }
  },
  "security": [
    {
      "basicAuth": null
    }
  ],
  "tags": [
    {
      "description": "manage pipelines",
      "name": "pipeline"
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "https",
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Swagger 2.0 specification for SignalCD",
    "title": "SignalCD Swagger Spec",
    "termsOfService": "http://swagger.io/terms/",
    "version": "v0.0.0"
  },
  "host": "localhost:6660",
  "basePath": "/api/v1",
  "paths": {
    "/deployments": {
      "get": {
        "tags": [
          "deployments"
        ],
        "summary": "Returns the history of deployments",
        "operationId": "deployments",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/deployment"
              }
            }
          },
          "500": {
            "description": "internal server error"
          }
        }
      }
    },
    "/deployments/current": {
      "get": {
        "tags": [
          "deployments"
        ],
        "summary": "Returns the currently active deployment",
        "operationId": "currentDeployment",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/deployment"
            }
          }
        }
      },
      "post": {
        "tags": [
          "deployments"
        ],
        "summary": "Schedule a new deployment",
        "operationId": "setCurrentDeployment",
        "parameters": [
          {
            "type": "string",
            "name": "pipeline",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/deployment"
            }
          },
          "500": {
            "description": "internal server error"
          }
        }
      }
    },
    "/pipelines": {
      "get": {
        "tags": [
          "pipeline"
        ],
        "summary": "returns a list of all pipelines",
        "operationId": "pipelines",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/pipeline"
              }
            }
          },
          "400": {
            "description": "bad request"
          },
          "500": {
            "description": "internal server error"
          }
        }
      },
      "post": {
        "tags": [
          "pipeline"
        ],
        "summary": "creates a new pipeline",
        "operationId": "create",
        "parameters": [
          {
            "description": "The pipeline to be created",
            "name": "pipeline",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/pipeline"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/pipeline"
            }
          },
          "500": {
            "description": "internal server error"
          }
        }
      }
    },
    "/pipelines/{id}": {
      "get": {
        "tags": [
          "pipeline"
        ],
        "summary": "returns a pipeline by id",
        "operationId": "pipeline",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/pipeline"
            }
          },
          "400": {
            "description": "bad request"
          },
          "500": {
            "description": "internal server error"
          }
        }
      }
    }
  },
  "definitions": {
    "check": {
      "required": [
        "name",
        "image"
      ],
      "properties": {
        "duration": {
          "type": "number"
        },
        "environment": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "key": {
                "type": "string"
              },
              "value": {
                "type": "string"
              }
            }
          }
        },
        "image": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      }
    },
    "deployment": {
      "required": [
        "number"
      ],
      "properties": {
        "created": {
          "type": "string",
          "format": "date-time"
        },
        "finished": {
          "type": "string",
          "format": "date-time"
        },
        "number": {
          "type": "number",
          "format": "int64"
        },
        "pipeline": {
          "$ref": "#/definitions/pipeline"
        },
        "started": {
          "type": "string",
          "format": "date-time"
        },
        "status": {
          "$ref": "#/definitions/deploymentstatus"
        }
      }
    },
    "deploymentstatus": {
      "properties": {
        "phase": {
          "type": "string",
          "enum": [
            "unknown",
            "success",
            "failure",
            "progress"
          ]
        }
      }
    },
    "pipeline": {
      "properties": {
        "checks": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/check"
          }
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "name": {
          "type": "string"
        },
        "steps": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/step"
          }
        }
      }
    },
    "step": {
      "required": [
        "name",
        "image"
      ],
      "properties": {
        "commands": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "image": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      }
    }
  },
  "security": [
    {
      "basicAuth": []
    }
  ],
  "tags": [
    {
      "description": "manage pipelines",
      "name": "pipeline"
    }
  ]
}`))
}
