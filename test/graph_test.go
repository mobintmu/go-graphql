package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-graphql/internal/config"
	"go-graphql/internal/product/dto"
	"io"
	"net/http"
	"testing"
)

func TestProductsGraphQL(t *testing.T) {
	WithHttpTestServer(t, func() {
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		addr := fmt.Sprintf("http://%s:%d", cfg.HTTPAddress, cfg.HTTPPort)

		// First, create a product via admin API so it's available to GraphQL
		var product dto.ProductResponse
		token := ""
		adminCreateProduct(t, &product, addr+"/api/v1/admin/products", token)

		// GraphQL queries
		queryProductByID(t, product, addr)
		queryProductsList(t, product, addr)

		// Cleanup
		adminDeleteProduct(t, product, addr+"/api/v1/admin/products", token)
	})
}

func queryProductByID(t *testing.T, product dto.ProductResponse, addr string) {
	t.Run("GraphQL: Get Product By ID", func(t *testing.T) {
		query := fmt.Sprintf(`{"query":"query { product(id: %d) { id name price } }"}`, product.ID)
		resp, err := http.Post(addr+"/query", "application/json", bytes.NewBufferString(query))
		if err != nil {
			t.Fatalf("Failed to send GraphQL request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			t.Logf("Response body: %s", string(body))
			return
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode GraphQL response: %v", err)
		}

		data := result["data"].(map[string]interface{})
		p := data["product"].(map[string]interface{})

		if int32(p["id"].(float64)) != product.ID {
			t.Errorf("Expected product ID %d, got %v", product.ID, p["id"])
		}
		if p["name"].(string) != product.Name {
			t.Errorf("Expected product name %q, got %q", product.Name, p["name"])
		}
	})
}

func queryProductsList(t *testing.T, product dto.ProductResponse, addr string) {
	t.Run("GraphQL: List Products", func(t *testing.T) {
		query := `{"query":"query { products { id name description price } }"}`
		resp, err := http.Post(addr+"/query", "application/json", bytes.NewBufferString(query))
		if err != nil {
			t.Fatalf("Failed to send GraphQL request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			t.Logf("Response body: %s", string(body))
			return
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode GraphQL response: %v", err)
		}

		data := result["data"].(map[string]interface{})
		products := data["products"].([]interface{})

		found := false
		for _, p := range products {
			pm := p.(map[string]interface{})
			if int32(pm["id"].(float64)) == product.ID {
				found = true
				t.Logf("Found product in GraphQL list: %+v", pm)
				break
			}
		}
		if !found {
			t.Errorf("Product not found in GraphQL products list")
		}
	})
}
