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
  "schemes": [
    "https",
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Isakura-manager",
    "title": "isakura-manager",
    "termsOfService": "http://swagger.io/terms/",
    "contact": {
      "email": "apiteam@swagger.io"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0.0"
  },
  "paths": {
    "/channel": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get all channels",
        "operationId": "getChannels",
        "parameters": [
          {
            "type": "string",
            "description": "Search Param to look for in programming",
            "name": "search",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "woot!",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Channel"
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Channel": {
      "type": "object",
      "properties": {
        "channelName": {
          "type": "string"
        },
        "guide": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Guide"
          }
        },
        "id": {
          "type": "integer",
          "format": "int64"
        }
      }
    },
    "Guide": {
      "type": "object",
      "properties": {
        "airdate": {
          "type": "string",
          "format": "date-time"
        },
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        }
      },
      "xml": {
        "name": "Category"
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "https",
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Isakura-manager",
    "title": "isakura-manager",
    "termsOfService": "http://swagger.io/terms/",
    "contact": {
      "email": "apiteam@swagger.io"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0.0"
  },
  "paths": {
    "/channel": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "get all channels",
        "operationId": "getChannels",
        "parameters": [
          {
            "type": "string",
            "description": "Search Param to look for in programming",
            "name": "search",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "woot!",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Channel"
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Channel": {
      "type": "object",
      "properties": {
        "channelName": {
          "type": "string"
        },
        "guide": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Guide"
          }
        },
        "id": {
          "type": "integer",
          "format": "int64"
        }
      }
    },
    "Guide": {
      "type": "object",
      "properties": {
        "airdate": {
          "type": "string",
          "format": "date-time"
        },
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        }
      },
      "xml": {
        "name": "Category"
      }
    }
  }
}`))
}