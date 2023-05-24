package route

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/WebEngrChild/go-graphql-server/pkg/adapter/http/handler"
	"github.com/WebEngrChild/go-graphql-server/pkg/domain/model"
	"github.com/WebEngrChild/go-graphql-server/pkg/domain/model/graph"
	repository "github.com/WebEngrChild/go-graphql-server/pkg/lib/mock"
	"github.com/WebEngrChild/go-graphql-server/pkg/usecase"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestQueryRoute(t *testing.T) {
	// Generate Mock Repository
	ctrlMessage := gomock.NewController(t)
	defer ctrlMessage.Finish()
	mmr := repository.NewMockMessage(ctrlMessage)
	mRes := []*model.Message{
		{ID: "1", Text: "testMessage1", UserID: "1", CreatedAt: "", UpdatedAt: ""},
	}
	mmr.EXPECT().GetMessages().Return(mRes, nil)

	ctrlUser := gomock.NewController(t)
	defer ctrlUser.Finish()
	mur := repository.NewMockUser(ctrlUser)
	ids := []string{"1"}
	uRes := map[string]*graph.User{
		"1": {ID: "1", Name: "testUser1"},
	}
	mur.EXPECT().GetMapInIDs(gomock.Any(), ids).Return(uRes, nil)

	// DI
	mu := usecase.NewMsgUseCase(mmr)
	uu := usecase.NewUserUseCase(mur)
	gh := handler.NewGraphHandler(mu, uu)

	// Routing
	e := echo.New()
	e.POST("/query", gh.QueryHandler())

	// Create a new http request to test the server
	reqFile := "testdata/ok_req.json.golden"
	reqbody := loadFile(t, reqFile)
	body := strings.NewReader(string(reqbody))
	req := httptest.NewRequest(http.MethodPost, "/query", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Serve the http request
	e.ServeHTTP(rec, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rec.Code)
	rspFile := "testdata/ok_res.json.golden"
	wnt := loadFile(t, rspFile)
	assert.JSONEq(t, string(wnt), string(rec.Body.Bytes()))
}

func loadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return bt
}
