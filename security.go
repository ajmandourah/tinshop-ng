package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/DblK/tinshop/config"
	"github.com/DblK/tinshop/utils"
)

// Middleware to ensure not forged query and real tinfoil client
func tinfoilMiddleware(next http.Handler) http.Handler {
	shopTemplate, _ := template.ParseFS(assetData, "assets/shop.tmpl")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify all headers
		headers := r.Header

		if config.GetConfig().DebugNoSecurity() {
			next.ServeHTTP(w, r)
			return
		}

		if r.RequestURI == "/" || utils.IsValidFilter(cleanPath(r.RequestURI)) {
			// Check for blacklist/whitelist
			var uid = strings.Join(headers["Uid"], "")
			if config.GetConfig().IsBlacklisted(uid) {
				log.Println("[Security] Blacklisted switch detected...", uid)
				_ = shopTemplate.Execute(w, config.GetConfig().ShopTemplateData())
				return
			}

			// Check for banned theme
			var theme = strings.Join(headers["Theme"], "")
			if config.GetConfig().IsBannedTheme(theme) {
				log.Println("[Security] Banned theme detected...", uid, theme)
				_ = shopTemplate.Execute(w, config.GetConfig().ShopTemplateData())
				return
			}

			// No User-Agent for tinfoil app
			if headers["User-Agent"] != nil {
				log.Println("[Security] User-Agent detected...")
				_ = shopTemplate.Execute(w, config.GetConfig().ShopTemplateData())
				return
			}

			// Be sure all tinfoil headers are present
			if headers["Theme"] == nil || headers["Uid"] == nil || headers["Version"] == nil || headers["Language"] == nil || headers["Hauth"] == nil || headers["Uauth"] == nil {
				log.Println("[Security] Missing some expected headers...")
				_ = shopTemplate.Execute(w, config.GetConfig().ShopTemplateData())
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
