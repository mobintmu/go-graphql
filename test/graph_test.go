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

func TestProductsGraphQLFilters(t *testing.T) {
	WithHttpTestServer(t, func() {
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		addr := fmt.Sprintf("http://%s:%d", cfg.HTTPAddress, cfg.HTTPPort)

		// Seed product
		var product dto.ProductResponse
		token := ""
		adminCreateProduct(t, &product, addr+"/api/v1/admin/products", token)

		// Run filter tests
		queryProductsByFilter(t, addr, fmt.Sprintf(`id: %d`, product.ID), func(p map[string]interface{}) bool {
			return int32(p["id"].(float64)) == product.ID
		})

		queryProductsByFilter(t, addr, fmt.Sprintf(`name: \"%s\"`, product.Name), func(p map[string]interface{}) bool {
			return p["name"].(string) == product.Name
		})

		queryProductsByFilter(t, addr, fmt.Sprintf(`description: \"%s\"`, product.Description), func(p map[string]interface{}) bool {
			return p["description"].(string) == product.Description
		})

		queryProductsByFilter(t, addr, fmt.Sprintf(`minPrice: %d`, product.Price-1), func(p map[string]interface{}) bool {
			return int64(p["price"].(float64)) >= product.Price-1
		})

		queryProductsByFilter(t, addr, fmt.Sprintf(`maxPrice: %d`, product.Price+1), func(p map[string]interface{}) bool {
			return int64(p["price"].(float64)) <= product.Price+1
		})

		queryProductsByFilter(t, addr, fmt.Sprintf(`isActive: %t`, true), func(p map[string]interface{}) bool {
			return p["isActive"].(bool) == true
		})

		// Cleanup
		adminDeleteProduct(t, product, addr+"/api/v1/admin/products", token)
	})
}

// Helper to run a GraphQL query with a given filter and validate results
func queryProductsByFilter(t *testing.T, addr, filter string, validate func(map[string]interface{}) bool) {
	t.Run("GraphQL: Filter "+filter, func(t *testing.T) {
		query := fmt.Sprintf(`{"query":"query { products(filter: { %s }) { products { id name description price isActive } total } }"}`, filter)
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
		conn := data["products"].(map[string]interface{})
		products := conn["products"].([]interface{})

		if len(products) == 0 {
			t.Errorf("Expected at least one product for filter %s", filter)
			return
		}

		for _, p := range products {
			pm := p.(map[string]interface{})
			if !validate(pm) {
				t.Errorf("Product %+v did not satisfy filter %s", pm, filter)
			}
		}
	})
}
