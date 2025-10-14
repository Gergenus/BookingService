package tests

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gergenus/bookingService/internal/handler"
	"github.com/Gergenus/bookingService/internal/models"
	"github.com/Gergenus/bookingService/internal/service/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateEquipment(t *testing.T) {
	tests := []struct {
		name                string
		inputBody           func() (*bytes.Buffer, string)
		inputEquipment      models.Equipment
		mockBehavior        func(ctx context.Context, equipment models.Equipment, image *multipart.FileHeader, tst *testing.T) *mocks.MockEquipmentServiceInterface
		expectedStatusCode  int
		expectedResposeBody string
	}{
		{
			name: "OK",
			inputBody: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				_ = writer.WriteField("equipment_name", "мультиинструмент")
				_ = writer.WriteField("manufacturer", `харьковский завод электроники "Электрон"`)
				_ = writer.WriteField("description", "вольтметр 1998 года выпуска")
				fileWriter, _ := writer.CreateFormFile("image", "enot.jpg")
				_, err := fileWriter.Write([]byte("FAKE"))
				if err != nil {
					t.Fatal(err)
				}
				writer.Close()
				return body, writer.FormDataContentType()
			},
			inputEquipment: models.Equipment{
				EquipmentName: "мультиинструмент",
				Manufacturer:  `харьковский завод электроники "Электрон"`,
				Description:   "вольтметр 1998 года выпуска",
			},
			mockBehavior: func(ctx context.Context, equipment models.Equipment, image *multipart.FileHeader, tst *testing.T) *mocks.MockEquipmentServiceInterface {
				mock := mocks.NewMockEquipmentServiceInterface(t)
				mock.EXPECT().CreateEquipment(ctx, equipment, image).Return(1, nil)
				return mock
			},
			expectedStatusCode:  http.StatusOK,
			expectedResposeBody: `{"id":1}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, contentType := tt.inputBody()
			copyBody := bytes.NewBuffer(body.Bytes())

			req, err := http.NewRequest("POST", "/equipment", copyBody)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", contentType)

			err = req.ParseMultipartForm(32 << 20)
			if err != nil {
				t.Fatal(err)
			}

			file, fileHeader, err := req.FormFile("image")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			mock := tt.mockBehavior(context.Background(), tt.inputEquipment, fileHeader, t)

			handler := handler.NewEquipmentHandler(mock)

			e := echo.New()

			e.POST("/api/v1/equipment/create", handler.CreateEquipment)

			w := httptest.NewRecorder()
			req = httptest.NewRequest("POST", "/api/v1/equipment/create", body)
			req.Header.Set("Content-Type", contentType)

			e.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedResposeBody, w.Body.String())
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}

}
