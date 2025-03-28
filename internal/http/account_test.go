package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vlone310/bss/internal/db/mock"
	db "github.com/vlone310/bss/internal/db/sqlc"
	"github.com/vlone310/bss/testutil"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {

				// build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				testutil.RequireBodyMatch(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Validation Failed",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		// TODO: add more cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:        testutil.RandomInt(1, 1000),
		Owner:     testutil.RandomOwner(),
		Balance:   testutil.RandomMoney(),
		Currency:  testutil.RandomCurrency(),
		CreatedAt: pgtype.Timestamptz{Valid: true},
	}
}
