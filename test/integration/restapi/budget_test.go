package restapi

// func TestGetBudgetRoute(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		description      string
// 		token            string
// 		expectedCode     int
// 		expectedResponse *dto.GetBudgetResponse
// 	}{
// 		{
// 			description:      "Fail to list budgets without token",
// 			token:            "",
// 			expectedCode:     http.StatusBadRequest,
// 			expectedResponse: nil,
// 		},
// 		{
// 			description:  "Get budget",
// 			token:        mockoauth.DefaultMockToken,
// 			expectedCode: http.StatusOK,
// 			expectedResponse: func() *dto.GetBudgetResponse {
// 				date, _ := time.Parse(
// 					time.RFC3339,
// 					"2025-03-09 18:35:32.556-03",
// 				)

// 				return &dto.GetBudgetResponse{
// 					GetBudgetOutput: usecase.GetBudgetOutput{
// 						Budget: entity.Budget{
// 							ID: uuid.MustParse(
// 								"8aa317f8-702c-43b1-897b-e24a4285d2d2",
// 							),
// 							Amount: 500000,
// 							Date:   date,
// 						},
// 						Spent:                              0,
// 						Available:                          0,
// 						AvailablePercentageVariation:       0,
// 						AvailablePerDay:                    0,
// 						AvailablePerDayPercentageVariation: 0,
// 						ComparisonDates:                    usecase.ComparisonDates{},
// 						BudgetCategories:                   []usecase.GetBudgetBudgetCategories{},
// 					},
// 				}
// 			}(),
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.description, func(t *testing.T) {
// 			t.Parallel()

// 			app, cleanUp := NewTestApp(t)
// 			defer func() {
// 				err := cleanUp(context.Background())
// 				assert.Nil(t, err)
// 			}()

// 			accessToken := ""
// 			if test.token != "" {
// 				signInRes := app.SignIn(test.token)
// 				accessToken = signInRes.AccessToken
// 			}

// 			var out dto.GetBudgetResponse
// 			statusCode, rawBody, err := app.MakeRequest(
// 				http.MethodGet,
// 				"/api/v1/budgets",
// 				WithBearerToken(accessToken),
// 				WithResponse(&out),
// 			)
// 			assert.Nil(t, err)

// 			assert.Equal(
// 				t,
// 				test.expectedCode,
// 				statusCode,
// 				rawBody,
// 			)
// 		})
// 	}
// }
