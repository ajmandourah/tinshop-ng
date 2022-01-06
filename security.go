package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/DblK/tinshop/utils"
)

// CORSMiddleware is a middleware to ensure right CORS headers
func (s *TinShop) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/api/") {
			w.Header().Set("Access-Control-Allow-Origin", s.Shop.Config.RootShop())
			w.Header().Set("Vary", "Origin")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// TinfoilMiddleware is a middleware to ensure not forged query and real tinfoil client
func (s *TinShop) TinfoilMiddleware(next http.Handler) http.Handler {
	shopTemplate, _ := template.ParseFS(assetData, "assets/shop.tmpl")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify all headers
		headers := r.Header

		if s.Shop.Config.DebugNoSecurity() {
			next.ServeHTTP(w, r)
			return
		}

		if r.RequestURI == "/" || utils.IsValidFilter(cleanPath(r.RequestURI)) {
			// Check for blacklist/whitelist
			var uid = strings.Join(headers["Uid"], "")
			if s.Shop.Config.IsBlacklisted(uid) {
				log.Println("[Security] Blacklisted switch detected...", uid)
				_ = shopTemplate.Execute(w, s.Shop.Config.ShopTemplateData())
				return
			}

			// Check for banned theme
			var theme = strings.Join(headers["Theme"], "")
			if s.Shop.Config.IsBannedTheme(theme) {
				log.Println("[Security] Banned theme detected...", uid, theme)
				_ = shopTemplate.Execute(w, s.Shop.Config.ShopTemplateData())
				return
			}

			// No User-Agent for tinfoil app
			if headers["User-Agent"] != nil {
				log.Println("[Security] User-Agent detected...")
				_ = shopTemplate.Execute(w, s.Shop.Config.ShopTemplateData())
				return
			}

			// Be sure all tinfoil headers are present
			if headers["Theme"] == nil || headers["Uid"] == nil || headers["Version"] == nil || headers["Language"] == nil || headers["Hauth"] == nil || headers["Uauth"] == nil {
				log.Println("[Security] Missing some expected headers...")
				_ = shopTemplate.Execute(w, s.Shop.Config.ShopTemplateData())
				return
			}

			// Enforce true tinfoil queries
			// TODO: Check Uauth and Hauth headers
			log.Printf("Switch %s, %s, %s, %s, %s, %s requesting %s", headers["Theme"], headers["Uid"], headers["Version"], headers["Language"], headers["Hauth"], headers["Uauth"], r.RequestURI)
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func cleanPath(path string) string {
	actualPath := path[1:]
	if path[len(path)-1:] == "/" {
		actualPath = path[1 : len(path)-1]
	}
	return actualPath
}
