package handlers

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/model"
	usecaseMock "MAXPUMP1/pkg/usecase/mock"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	testCase := map[string]struct {
		userInput     interface{}
		usecaseFunc   func(usecaseMock *usecaseMock.MockUserUsecaseInterface, signupData model.Signup)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		"validSignup": {
			userInput: model.Signup{
				FirstName: "test1",
				LastName:  "test2",
				Email:     "test@example.com",
				Phone:     "1234567890",
				Password:  "password123",
			},
			usecaseFunc: func(usecaseMock *usecaseMock.MockUserUsecaseInterface, signupData model.Signup) {
				expectedUser := entity.User{
					FirstName: "test1",
					LastName:  "test2",
					Email:     "test@example.com",
					Phone:     "1234567890",
				}
				usecaseMock.EXPECT().ExecuteSignup(gomock.Any()).Return(&expectedUser, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {

				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.Contains(t, rec.Body.String(), "test@example.com")
				assert.Contains(t, rec.Body.String(), "test1")
			},
		},
		"Signup Conflict": {
			userInput: model.Signup{
				FirstName: "test1",
				LastName:  "test2",
				Email:     "test@example.com",
				Phone:     "1234567890",
				Password:  "password123",
			},
			usecaseFunc: func(usecaseMock *usecaseMock.MockUserUsecaseInterface, signupData model.Signup) {
				usecaseMock.EXPECT().ExecuteSignup(gomock.Any()).Return(nil, errors.New("user already exists"))
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Body.String(), "user already exists")
			},
		},
		"Invalid Input": {
			userInput: model.Signup{
				Email: "invalid email",
			},
			usecaseFunc: nil, //in this "Invalid Input "case there is no need f Execute Signup call so usecaseFunc set as nil
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), "error")

			},
		},
		"Malformed JSON for Binding Error": {
			userInput: `{
				"FirstName": "test1",
				"LastName":  "test2",
				"Email":     "test@example.com",
				"Phone":     "1234567890",
				"Password":  "password123"`,
			usecaseFunc: nil,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), "error")
			},
		},
	}

	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := usecaseMock.NewMockUserUsecaseInterface(ctrl)
			if tc.usecaseFunc != nil {
				tc.usecaseFunc(mockUsecase, tc.userInput.(model.Signup))
			}
			uh := &UserHandler{UserUsecase: mockUsecase}

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/signup", uh.Signup)

			var body []byte
			var err error
			if jsonString, ok := tc.userInput.(string); ok {
				body = []byte(jsonString)
			} else {
				body, err = json.Marshal(tc.userInput)
				if err != nil {
					t.Fatalf("could not marshal request: %v", err)
				}
			}

			req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

// func TestSignupWithOtp(t *testing.T){

// }
