// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Daniel Mesquita",
            "email": "danielmesquitta123@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/health": {
            "get": {
                "description": "Health check",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.HealthResponseDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    }
                }
            }
        },
        "/v1/auth/sign-in": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Authenticate user through Google or Apple token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Sign in",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.SignInRequestDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SignInResponseDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    }
                }
            }
        },
        "/v1/calculator/compound-interest": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Calculate compound interest",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Calculator"
                ],
                "summary": "Calculate compound interest",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CompoundInterestRequestDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.CompoundInterestResponseDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    }
                }
            }
        },
        "/v1/calculator/emergency-reserve": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Calculate emergency reserve",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Calculator"
                ],
                "summary": "Calculate emergency reserve",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.EmergencyReserveRequestDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.EmergencyReserveResponseDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    }
                }
            }
        },
        "/v1/calculator/retirement": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Calculate investments needed for retirement",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Calculator"
                ],
                "summary": "Calculate retirement",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RetirementRequestDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.RetirementResponseDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    }
                }
            }
        },
        "/v1/calculator/simple-interest": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Calculate simple interest",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Calculator"
                ],
                "summary": "Calculate simple interest",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.SimpleInterestRequestDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SimpleInterestResponseDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponseDTO"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.CompoundInterestRequestDTO": {
            "type": "object",
            "required": [
                "interest",
                "interest_type",
                "period_in_months"
            ],
            "properties": {
                "initial_deposit": {
                    "type": "number",
                    "minimum": 0
                },
                "interest": {
                    "type": "number",
                    "maximum": 100,
                    "minimum": 0
                },
                "interest_type": {
                    "enum": [
                        "MONTHLY",
                        "ANNUAL"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/entity.InterestType"
                        }
                    ]
                },
                "monthly_deposit": {
                    "type": "number"
                },
                "period_in_months": {
                    "type": "integer",
                    "minimum": 1
                }
            }
        },
        "dto.CompoundInterestResponseDTO": {
            "type": "object",
            "properties": {
                "by_month": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/usecase.CompoundInterestResult"
                    }
                },
                "total_amount": {
                    "type": "number"
                },
                "total_deposit": {
                    "type": "number"
                },
                "total_interest": {
                    "type": "number"
                }
            }
        },
        "dto.EmergencyReserveRequestDTO": {
            "type": "object",
            "required": [
                "jobType"
            ],
            "properties": {
                "jobType": {
                    "enum": [
                        "ENTREPRENEUR",
                        "EMPLOYEE",
                        "CIVIL_SERVANT"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/entity.JobType"
                        }
                    ]
                },
                "monthlyExpenses": {
                    "type": "number",
                    "minimum": 0
                },
                "monthlyIncome": {
                    "type": "number",
                    "minimum": 0
                },
                "monthlySavingsPercentage": {
                    "type": "number",
                    "maximum": 100,
                    "minimum": 0
                }
            }
        },
        "dto.EmergencyReserveResponseDTO": {
            "type": "object",
            "properties": {
                "monthsToAchieveEmergencyReserve": {
                    "type": "integer"
                },
                "recommendedReserveInMonths": {
                    "type": "integer"
                },
                "recommendedReserveInValue": {
                    "type": "number"
                }
            }
        },
        "dto.ErrorResponseDTO": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.HealthResponseDTO": {
            "type": "object",
            "properties": {
                "ok": {
                    "type": "boolean"
                }
            }
        },
        "dto.RetirementRequestDTO": {
            "type": "object",
            "required": [
                "age",
                "goal_income",
                "goal_patrimony",
                "income_investment_percentage",
                "initial_deposit",
                "interest",
                "interest_type",
                "monthly_income",
                "retirement_age"
            ],
            "properties": {
                "age": {
                    "type": "integer",
                    "minimum": 0
                },
                "goal_income": {
                    "type": "number",
                    "minimum": 0
                },
                "goal_patrimony": {
                    "type": "number",
                    "minimum": 0
                },
                "income_investment_percentage": {
                    "type": "number",
                    "maximum": 100,
                    "minimum": 0
                },
                "initial_deposit": {
                    "type": "number",
                    "minimum": 0
                },
                "interest": {
                    "type": "number",
                    "maximum": 100,
                    "minimum": 0
                },
                "interest_type": {
                    "enum": [
                        "MONTHLY",
                        "ANNUAL"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/entity.InterestType"
                        }
                    ]
                },
                "life_expectancy": {
                    "type": "integer",
                    "minimum": 1
                },
                "monthly_income": {
                    "type": "number",
                    "minimum": 0
                },
                "retirement_age": {
                    "type": "integer",
                    "minimum": 1
                }
            }
        },
        "dto.RetirementResponseDTO": {
            "type": "object",
            "properties": {
                "achieved_goal_income": {
                    "type": "boolean"
                },
                "achieved_goal_patrimony": {
                    "type": "boolean"
                },
                "heritage": {
                    "type": "number"
                },
                "max_monthly_expenses": {
                    "type": "number"
                },
                "property_on_retirement": {
                    "type": "number"
                }
            }
        },
        "dto.SignInRequestDTO": {
            "type": "object",
            "properties": {
                "provider": {
                    "type": "string"
                }
            }
        },
        "dto.SignInResponseDTO": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/entity.User"
                }
            }
        },
        "dto.SimpleInterestRequestDTO": {
            "type": "object",
            "required": [
                "interest",
                "interestType",
                "periodInMonths"
            ],
            "properties": {
                "initialDeposit": {
                    "type": "number",
                    "minimum": 0
                },
                "interest": {
                    "type": "number",
                    "maximum": 100,
                    "minimum": 0
                },
                "interestType": {
                    "enum": [
                        "MONTHLY",
                        "ANNUAL"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/entity.InterestType"
                        }
                    ]
                },
                "periodInMonths": {
                    "type": "integer",
                    "minimum": 1
                }
            }
        },
        "dto.SimpleInterestResponseDTO": {
            "type": "object",
            "properties": {
                "by_month": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/usecase.SimpleInterestResult"
                    }
                },
                "total_amount": {
                    "type": "number"
                },
                "total_deposit": {
                    "type": "number"
                },
                "total_interest": {
                    "type": "number"
                }
            }
        },
        "entity.InterestType": {
            "type": "string",
            "enum": [
                "MONTHLY",
                "ANNUAL"
            ],
            "x-enum-varnames": [
                "InterestTypeMonthly",
                "InterestTypeAnnual"
            ]
        },
        "entity.JobType": {
            "type": "string",
            "enum": [
                "ENTREPRENEUR",
                "EMPLOYEE",
                "CIVIL_SERVANT"
            ],
            "x-enum-varnames": [
                "JobTypeEntrepreneur",
                "JobTypeEmployee",
                "JobTypeCivilServant"
            ]
        },
        "entity.User": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "external_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "provider": {
                    "type": "string"
                },
                "subscription_expires_at": {
                    "type": "string"
                },
                "synchronized_at": {
                    "type": "string"
                },
                "tier": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "verified_email": {
                    "type": "boolean"
                }
            }
        },
        "usecase.CompoundInterestResult": {
            "type": "object",
            "properties": {
                "monthly_interest": {
                    "type": "number"
                },
                "total_amount": {
                    "type": "number"
                },
                "total_deposit": {
                    "type": "number"
                },
                "total_interest": {
                    "type": "number"
                }
            }
        },
        "usecase.SimpleInterestResult": {
            "type": "object",
            "properties": {
                "total_amount": {
                    "type": "number"
                },
                "total_deposit": {
                    "type": "number"
                },
                "total_interest": {
                    "type": "number"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        },
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
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "API Finance Manager",
	Description:      "API Finance Manager",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
