package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"webinar-testing/internal/api/mocks"
	"webinar-testing/pkg/models"
)

func Test_server_Add(t *testing.T) {
	type fields struct {
		serv    *echo.Echo
		service *mocks.Service
	}
	type args struct {
		reqRawBody []byte
		order      models.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				serv:    echo.New(),
				service: mocks.NewService(t),
			},
			args: args{
				reqRawBody: nil,
				order: models.Order{
					UserID: "user1",
					Goods:  make(map[models.GoodID]int),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				serv:    tt.fields.serv,
				service: tt.fields.service,
			}
			reqBody := tt.args.reqRawBody
			if len(reqBody) == 0 {
				var err error
				reqBody, err = json.Marshal(&tt.args.order)
				assert.NoError(t, err)
			}

			t.Log(string(reqBody))

			req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := s.serv.NewContext(req, rec)

			tt.fields.service.EXPECT().Add(c.Request().Context(), tt.args.order).Return(nil)

			if err := s.Add(c); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Zero(t, rec.Body.Len())
		})
	}
}
