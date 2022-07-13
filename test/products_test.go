package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ignaciofalco/new-store/cmd/server/handler"
	"github.com/ignaciofalco/new-store/internal/products"
	"github.com/ignaciofalco/new-store/pkg/store"
	"github.com/stretchr/testify/assert"
)

func copyFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func createServer(pathDB string) *gin.Engine {
	_ = os.Setenv("TOKEN", "123456")
	err := copyFile(pathDB, fmt.Sprintf("tmp_%s", pathDB))
	if err != nil {
		panic(err)
	}
	db := store.New(store.FileType, pathDB)
	repo := products.NewRepository(db)
	service := products.NewService(repo)
	p := handler.NewProduct(service)
	r := gin.Default()

	pr := r.Group("/products")
	pr.POST("/", p.Store())
	pr.GET("/", p.GetAll())
	pr.PATCH("/:id", p.UpdateName())

	return r
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", "123456")

	return req, httptest.NewRecorder()
}

func TestGetAllProducts(t *testing.T) {

	type producto struct {
		ID       int    `json:"id"`
		Nombre   string `json:"nombre"`
		Tipo     string `json:"tipo"`
		Cantidad int    `json:"cantidad"`
		Precio   int    `json:"precio"`
	}

	// crear el Server y definir las Rutas
	r := createServer("products.json")
	// crear Request del tipo GET y Response para obtener el resultado
	req, rr := createRequestTest(http.MethodGet, "/products/", "")

	var objRes []producto

	// indicar al servidor que pueda atender la solicitud
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.True(t, len(objRes) > 0)
}

func TestSaveProduct(t *testing.T) {
	// crear el Server y definir las Rutas
	r := createServer("products.json")
	// crear Request del tipo POST y Response para obtener el resultado
	req, rr := createRequestTest(http.MethodPost, "/products/", `{
        "nombre": "Tester","tipo": "Funcional","cantidad": 10,"precio": 99.99
    }`)

	// indicar al servidor que pueda atender la solicitud
	r.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

}

func TestUpdateNameProduct(t *testing.T) {
	// crear el Server y definir las Rutas
	r := createServer("product_update_name.json")
	// crear Request del tipo POST y Response para obtener el resultado
	req, rr := createRequestTest(http.MethodPatch, "/products/3", `{"nombre": "Arroz"}`)

	// indicar al servidor que pueda atender la solicitud
	r.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

}
