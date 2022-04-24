package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vendingmachine/config"
	"vendingmachine/domains"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testDeposit struct {
	name                 string
	userID               string
	testID               string
	mockMakeUserDeposit  func()
	mockResetUserDeposit func()
}

func TestMakeDepositEmptyBody(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/deposit", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer-user-test")
	MakeDeposit(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestMakeDepositInvalidBody(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader("{\"deposit\":\"aaa\"}"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer-user-test")
	MakeDeposit(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestMakeDepositBadCoin(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader("{\"deposit\": 2}"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer-user-test")
	MakeDeposit(w, r)
	assert.Equal(t, config.BadDepositError("").Code, w.Code)
}

func TestMakeDepositDataError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader("{\"deposit\": 10}"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer-user-test")
	testDepositDataError := testDeposit{
		name:   "Data Error Expected",
		userID: "buyer-user-test",
		testID: "MakeDepositDataError",
		mockMakeUserDeposit: func() {
			MakeUserDeposit = func(userDeposit *domains.UserDeposit) (*domains.UserDeposit, *config.HttpError) {
				return nil, config.DataAccessLayerError(gorm.ErrRecordNotFound.Error())
			}
		},
	}
	testDepositDataError.mockMakeUserDeposit()
	MakeDeposit(w, r)
	assert.Equal(t, config.DataAccessLayerError("").Code, w.Code)

}

func TestMakeDepositValid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader("{\"deposit\": 10}"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer-user-test")
	testValidDeposit := testDeposit{
		name:   "Valid Deposit Expected",
		userID: "buyer-user-test",
		testID: "ValidDeposit",
		mockMakeUserDeposit: func() {
			MakeUserDeposit = func(userDeposit *domains.UserDeposit) (*domains.UserDeposit, *config.HttpError) {
				return &domains.UserDeposit{
					Deposit: 10,
					Balance: 10,
					UserId:  "buyer-user-test",
				}, nil
			}
		},
	}
	testValidDeposit.mockMakeUserDeposit()
	MakeDeposit(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestResetUserDepositInvalidUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader("{\"deposit\": 10}"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	DeleteUserDeposit(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestResetDepositDataError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader("{\"deposit\": 10}"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer-user-test")
	testDepositDataError := testDeposit{
		name:   "Data Error Expected",
		userID: "buyer-user-test",
		testID: "ResetDepositDataError",
		mockResetUserDeposit: func() {
			ResetUserDeposit = func(id string) *config.HttpError {
				return config.NotFoundError()
			}
		},
	}
	testDepositDataError.mockResetUserDeposit()
	DeleteUserDeposit(w, r)
	assert.Equal(t, config.NotFoundError().Code, w.Code)

}

func TestResetUserDepositValid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader("{\"deposit\": 10}"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer-user-test")
	testResetValid := testDeposit{
		name:   "Valid Deposit Expected",
		userID: "buyer-user-test",
		testID: "ValidDeposit",
		mockResetUserDeposit: func() {
			ResetUserDeposit = func(id string) *config.HttpError {
				return nil
			}
		},
	}
	testResetValid.mockResetUserDeposit()
	DeleteUserDeposit(w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)

}
