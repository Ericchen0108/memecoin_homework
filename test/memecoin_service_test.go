package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"memecoin_homework/internal/model"
	"memecoin_homework/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type testEnv struct {
	db     *gorm.DB
	router *gin.Engine
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()

	dsn := "host=localhost user=user password=password dbname=memecoins port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.Memecoin{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	return &testEnv{
		db:     db,
		router: router,
	}
}

func (env *testEnv) cleanup(t *testing.T) {
	t.Helper()
	if err := env.db.Exec("TRUNCATE TABLE memecoins RESTART IDENTITY").Error; err != nil {
		t.Errorf("Failed to cleanup database: %v", err)
	}
}

func (env *testEnv) insertMockData(t *testing.T) []model.Memecoin {
	t.Helper()
	memecoins := []model.Memecoin{
		{Name: "TestCoin1", Description: "This is the first test coin."},
		{Name: "TestCoin2", Description: "This is the second test coin."},
		{Name: "TestCoin3", Description: "This is the third test coin."},
		{Name: "TestCoin4", Description: "This is the fourth test coin."},
		{Name: "TestCoin5", Description: "This is the fifth test coin."},
	}

	for _, coin := range memecoins {
		if err := env.db.Create(&coin).Error; err != nil {
			t.Fatalf("Failed to insert mock data: %v", err)
		}
	}
	return memecoins
}

func TestCreateMemecoin(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful creation",
			payload: `{
				"name": "NewTestCoin",
				"description": "This is a test coin created by test case"
			}`,
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "duplicate name",
			payload: `{
				"name": "TestCoin1",
				"description": "This is a duplicate coin"
			}`,
			wantStatus: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "empty name",
			payload: `{
				"name": "",
				"description": "This coin has no name"
			}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setupTestEnv(t)
			defer env.cleanup(t)
			env.insertMockData(t)

			env.router.POST("/memecoins", func(c *gin.Context) {
				service.CreateMemecoin(c, env.db)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/memecoins", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			env.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateMemecoin() status = %v, want %v", w.Code, tt.wantStatus)
				t.Logf("Response body: %s", w.Body.String())
			}

			if !tt.wantErr && w.Code == http.StatusCreated {
				var response model.Memecoin
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response.ID == 0 {
					t.Error("Expected ID to be set")
				}
			}
		})
	}
}

func TestGetMemecoin(t *testing.T) {
	tests := []struct {
		name       string
		coinID     string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "existing memecoin",
			coinID:     "1",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existing memecoin",
			coinID:     "999",
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "invalid id format",
			coinID:     "invalid",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setupTestEnv(t)
			defer env.cleanup(t)
			mockCoins := env.insertMockData(t)

			env.router.GET("/memecoins/:id", func(c *gin.Context) {
				service.GetMemecoin(c, env.db)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/memecoins/"+tt.coinID, nil)
			env.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetMemecoin() status = %v, want %v", w.Code, tt.wantStatus)
				t.Logf("Response body: %s", w.Body.String())
			}

			if !tt.wantErr && w.Code == http.StatusOK {
				var response model.Memecoin
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				id, _ := strconv.Atoi(tt.coinID)
				if id <= len(mockCoins) {
					expectedCoin := mockCoins[id-1]
					if response.Name != expectedCoin.Name {
						t.Errorf("Expected name %s, got %s", expectedCoin.Name, response.Name)
					}
				}
			}
		})
	}
}

func TestUpdateMemecoin(t *testing.T) {
	tests := []struct {
		name       string
		coinID     string
		payload    string
		wantStatus int
		wantErr    bool
	}{
		{
			name:   "successful update",
			coinID: "1",
			payload: `{
				"description": "Updated description"
			}`,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "non-existing memecoin",
			coinID: "999",
			payload: `{
				"description": "Updated description"
			}`,
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:   "invalid id format",
			coinID: "invalid",
			payload: `{
				"description": "Updated description"
			}`,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setupTestEnv(t)
			defer env.cleanup(t)
			env.insertMockData(t)

			env.router.PUT("/memecoins/:id", func(c *gin.Context) {
				service.UpdateMemecoin(c, env.db)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/memecoins/"+tt.coinID, strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			env.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("UpdateMemecoin() status = %v, want %v", w.Code, tt.wantStatus)
				t.Logf("Response body: %s", w.Body.String())
			}

			if !tt.wantErr && w.Code == http.StatusOK {
				var updatedCoin model.Memecoin
				id, _ := strconv.Atoi(tt.coinID)
				if err := env.db.First(&updatedCoin, id).Error; err != nil {
					t.Fatalf("Failed to fetch updated coin: %v", err)
				}
				if updatedCoin.Description != "Updated description" {
					t.Errorf("Expected description to be updated, got %s", updatedCoin.Description)
				}
			}
		})
	}
}

func TestDeleteMemecoin(t *testing.T) {
	tests := []struct {
		name       string
		coinID     string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			coinID:     "1",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existing memecoin",
			coinID:     "999",
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "invalid id format",
			coinID:     "invalid",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setupTestEnv(t)
			defer env.cleanup(t)
			env.insertMockData(t)

			env.router.DELETE("/memecoins/:id", func(c *gin.Context) {
				service.DeleteMemecoin(c, env.db)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/memecoins/"+tt.coinID, nil)
			env.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("DeleteMemecoin() status = %v, want %v", w.Code, tt.wantStatus)
				t.Logf("Response body: %s", w.Body.String())
			}

			if !tt.wantErr && w.Code == http.StatusOK {
				var deletedCoin model.Memecoin
				id, _ := strconv.Atoi(tt.coinID)
				err := env.db.First(&deletedCoin, id).Error
				if err == nil {
					t.Error("Expected record to be deleted, but it still exists")
				}
			}
		})
	}
}

func TestPokeMemecoin(t *testing.T) {
	tests := []struct {
		name            string
		coinID          string
		concurrentPokes int
		wantStatus      int
		wantErr         bool
	}{
		{
			name:            "successful poke",
			coinID:          "1",
			concurrentPokes: 1,
			wantStatus:      http.StatusOK,
			wantErr:         false,
		},
		{
			name:            "concurrent pokes",
			coinID:          "1",
			concurrentPokes: 10,
			wantStatus:      http.StatusOK,
			wantErr:         false,
		},
		{
			name:            "non-existing memecoin",
			coinID:          "999",
			concurrentPokes: 1,
			wantStatus:      http.StatusNotFound,
			wantErr:         true,
		},
		{
			name:            "invalid id format",
			coinID:          "invalid",
			concurrentPokes: 1,
			wantStatus:      http.StatusBadRequest,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setupTestEnv(t)
			defer env.cleanup(t)
			env.insertMockData(t)

			env.router.POST("/memecoins/:id/poke", func(c *gin.Context) {
				service.PokeMemecoin(c, env.db)
			})

			if tt.concurrentPokes > 1 {
				var wg sync.WaitGroup
				wg.Add(tt.concurrentPokes)

				for i := 0; i < tt.concurrentPokes; i++ {
					go func() {
						defer wg.Done()
						w := httptest.NewRecorder()
						req := httptest.NewRequest("POST", "/memecoins/"+tt.coinID+"/poke", nil)
						env.router.ServeHTTP(w, req)

						if w.Code != tt.wantStatus {
							t.Errorf("PokeMemecoin() status = %v, want %v", w.Code, tt.wantStatus)
						}
					}()
				}

				wg.Wait()

				if !tt.wantErr {
					var pokedCoin model.Memecoin
					id, _ := strconv.Atoi(tt.coinID)
					if err := env.db.First(&pokedCoin, id).Error; err != nil {
						t.Fatalf("Failed to fetch poked coin: %v", err)
					}
					if pokedCoin.PopularityScore != tt.concurrentPokes {
						t.Errorf("Expected popularity score %d, got %d", tt.concurrentPokes, pokedCoin.PopularityScore)
					}
				}
			} else {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/memecoins/"+tt.coinID+"/poke", nil)
				env.router.ServeHTTP(w, req)

				if w.Code != tt.wantStatus {
					t.Errorf("PokeMemecoin() status = %v, want %v", w.Code, tt.wantStatus)
					t.Logf("Response body: %s", w.Body.String())
				}

				if !tt.wantErr && w.Code == http.StatusOK {
					var pokedCoin model.Memecoin
					id, _ := strconv.Atoi(tt.coinID)
					if err := env.db.First(&pokedCoin, id).Error; err != nil {
						t.Fatalf("Failed to fetch poked coin: %v", err)
					}
					if pokedCoin.PopularityScore != 1 {
						t.Errorf("Expected popularity score 1, got %d", pokedCoin.PopularityScore)
					}
				}
			}
		})
	}
}
