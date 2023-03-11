package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestAuthMiddlewareWithoutCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "http://testing", nil)
	responseRecorder := httptest.NewRecorder()

	AuthMiddleware()(nil).ServeHTTP(responseRecorder, req)

	expectedCode := 400
	if responseRecorder.Code != expectedCode {
		t.Errorf("response code is %d, expected: %d", responseRecorder.Code, expectedCode)
	}

	expectedBody := "{\"message\":\"JWT is missing.\"}"
	if responseRecorder.Body.String() != expectedBody {
		t.Errorf("response body is %s, expected: %s", responseRecorder.Body, expectedBody)
	}
}

func TestAuthMiddlewareWithInvalidJWT(t *testing.T) {
	req := httptest.NewRequest("GET", "http://testing", nil)
	req.AddCookie(&http.Cookie{
		Name:  "access",
		Value: "foo",
	})

	responseRecorder := httptest.NewRecorder()

	AuthMiddleware()(nil).ServeHTTP(responseRecorder, req)

	expectedCode := 401
	if responseRecorder.Code != expectedCode {
		t.Errorf("response code is %d, expected: %d", responseRecorder.Code, expectedCode)
	}

	expectedBody := "{\"message\":\"JWT is invalid.\"}"
	if responseRecorder.Body.String() != expectedBody {
		t.Errorf("response body is %s, expected: %s", responseRecorder.Body, expectedBody)
	}

}

func TestXxx(t *testing.T) {
	testCases := []struct {
		name           string
		userId         int32
		wantStatusCode int
	}{
		{
			name:           "has userID",
			userId:         7,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "userID is zero",
			userId:         0,
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	handler := func(expectedUID int32) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			id, err := GetUserId(ctx)
			if err != nil {
				t.Error(err)
			}
			if id != expectedUID {
				t.Errorf("wanted userid: %d, but got: %d", expectedUID, id)
			}
		}
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			token := genTokenforTesting(t, test.userId)
			request.AddCookie(&http.Cookie{
				Name:  "access",
				Value: token,
			})

			responseRecorder := httptest.NewRecorder()

			mainHandler := AuthMiddleware()(handler(int32(test.userId)))
			mainHandler.ServeHTTP(responseRecorder, request)

			if want, got := test.wantStatusCode, responseRecorder.Code; want != got {
				t.Log(responseRecorder.Body)
				t.Fatalf("wanted status code %d, but got status code %d", want, got)

			}
		})
	}
}

func genTokenforTesting(t *testing.T, userID int32) string {
	token, err := GenToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	return token
}
