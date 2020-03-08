package server

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
)

func setWeb() {
	err := godotenv.Load()
	chk(err)

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	goth.UseProviders(discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), callbackURL))

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		if _, err := gothic.CompleteUserAuth(res, req); err == nil {
			sendIndex(res)
		} else {
			gothic.BeginAuthHandler(res, req)
		}
	})

	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		chk(err)

		for _, v := range permittedUsers {
			if v == user.UserID {
				sendIndex(res)
				break
			}
		}

		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte("Forbidden"))
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "/auth/discord", http.StatusPermanentRedirect)
	})
}

func sendIndex(res http.ResponseWriter) {
	page, err := ioutil.ReadFile("./views/index.html")
	chk(err)

	res.Write(page)
}
