package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/ajmandourah/tinshop-ng/utils"
	"golang.org/x/crypto/bcrypt"
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

		//Show Hauth for the specefied host
		//tinfoil sends requests appending "/" at the end
		if r.RequestURI == "/hauth/" && r.Header.Get("Hauth") != "" {
			log.Println("HAUTH for ", s.Shop.Config.Host(), " is: ", headers["Hauth"])
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

			// TODO: Here implement usage of IsWhitelisted

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
				log.Println("[Security] Missing some expected headers...an access attempted from a client other than tinfoil")
				_ = shopTemplate.Execute(w, s.Shop.Config.ShopTemplateData())
				return
			}


			//Hauth check
			if s.Shop.Config.Get_Hauth() != "" && r.Header.Get("Hauth") != s.Shop.Config.Get_Hauth(){
				log.Println("Hauth header mismatch. Possible attempt to access shop from a possible forged request. ", r.RemoteAddr)
				return
			}


			// Enforce true tinfoil queries
			// TODO: Check Uauth and Hauth headers
			log.Printf("Switch %s requesting %s", headers["Uid"], r.RequestURI)

			// Check user password
			if s.Shop.Config.ForwardAuthURL() != "" && headers["Authorization"] != nil {
				log.Println("[Security] Forwarding auth to", s.Shop.Config.ForwardAuthURL())
				client := &http.Client{}
				req, _ := http.NewRequest("GET", s.Shop.Config.ForwardAuthURL(), nil)
				req.Header.Add("Authorization", strings.Join(headers["Authorization"], ""))
				req.Header.Add("Device-Id", strings.Join(headers["Uid"], ""))
				resp, err := client.Do(req)
				if err != nil {
					log.Print(err)
					_ = shopTemplate.Execute(w, s.Shop.Config.ShopTemplateData())
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					log.Println("Wrong credentials enterd from switch ",r.Header.Get("Uid"), " " ,r.RemoteAddr)
					_ = shopTemplate.Execute(w, s.Shop.Config.ShopTemplateData())
					return
				}
			}
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

// HttpAuthCheck function checks for correct credentials
func HttpAuthCheck(user ,pass string, r *http.Request) bool {
	for _,cred := range creds {
		splitted := strings.Split(cred,":")
		if splitted[0] == user {
			err := bcrypt.CompareHashAndPassword([]byte(splitted[1]),[]byte(pass))
			if err == nil {
				return true
			}
		
		}
	}
	log.Println("An attempt to access the shop with username: ", user)
	return false

}
