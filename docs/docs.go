// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/api/v1/user/{id}/wallets": {
            "delete": {
                "description": "Delete wallet for the user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user wallet"
                ],
                "summary": "Delete wallet for the user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/users/{id}/wallets": {
            "get": {
                "description": "Get all wallets for the user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user wallet"
                ],
                "summary": "Get all wallets for the user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
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
                                "$ref": "#/definitions/wallet.Wallet"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/wallets": {
            "get": {
                "description": "Get all wallets",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "Get all wallets",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by wallet type",
                        "name": "wallet_type",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/wallet.Wallet"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    }
                }
            },
            "put": {
                "description": "Update wallet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "Update wallet",
                "parameters": [
                    {
                        "description": "Wallet object",
                        "name": "wallet",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/wallet.WalletForUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/wallet.Wallet"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    }
                }
            },
            "post": {
                "description": "Create wallet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "Create wallet",
                "parameters": [
                    {
                        "description": "Wallet object",
                        "name": "wallet",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/wallet.WalletForCreate"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/wallet.Wallet"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/wallet.Err"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "wallet.Err": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "wallet.Wallet": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 100
                },
                "created_at": {
                    "type": "string",
                    "example": "2024-03-25T14:19:00.729237Z"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                },
                "user_name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "wallet_name": {
                    "type": "string",
                    "example": "John's Wallet"
                },
                "wallet_type": {
                    "type": "string",
                    "example": "Credit Card"
                }
            }
        },
        "wallet.WalletForCreate": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 100
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                },
                "user_name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "wallet_name": {
                    "type": "string",
                    "example": "John's Wallet"
                },
                "wallet_type": {
                    "type": "string",
                    "example": "Credit Card"
                }
            }
        },
        "wallet.WalletForUpdate": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 100
                },
                "id": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:1323",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Wallet API",
	Description:      "Sophisticated Wallet API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
