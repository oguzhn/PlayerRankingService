package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oguzhn/PlayerRankingService/business"
	"github.com/oguzhn/PlayerRankingService/models"
)

type Controller struct {
	app business.IBusiness
}

func NewController(app business.IBusiness) *Controller {
	return &Controller{app: app}
}

func (c *Controller) RegisterHandlers() http.Handler {
	router := mux.NewRouter()
	router.Path("/leaderboard").Methods(http.MethodGet).HandlerFunc(c.GetLeaderBoard)
	router.Path("/leaderboard/{country_iso_code}").Methods(http.MethodGet).HandlerFunc(c.GetLeaderBoardByCountryCode)
	router.Path("/score/submit").Methods(http.MethodPost).HandlerFunc(c.SubmitScore)
	router.Path("/user/profile/{user_guid}").Methods(http.MethodGet).HandlerFunc(c.GetUserById)
	router.Path("/user/create").Methods(http.MethodPost).HandlerFunc(c.CreateUser)
	return router
}

func (c *Controller) GetLeaderBoard(w http.ResponseWriter, r *http.Request) {
	rs, err := c.app.GetLeaderBoard()
	if err != nil {
		log.Printf("Failed to load leaderboard: err: %s\n", err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		log.Printf("Failed to load marshal data: %v, err: %s\n", rs, err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func (c *Controller) GetLeaderBoardByCountryCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["country_iso_code"]

	rs, err := c.app.GetLeaderBoardByCountryCode(code)
	if err != nil {
		log.Printf("Failed to load leaderboard: err: %s\n", err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		log.Printf("Failed to load marshal data: %v, err: %s\n", rs, err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func (c *Controller) SubmitScore(w http.ResponseWriter, r *http.Request) {
	var score models.ScoreDTO
	err := json.NewDecoder(r.Body).Decode(&score)
	if err != nil {
		log.Printf("Failed to decode json form: %+v, err: %s\n", r.Body, err)
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	if err := score.IsValid(); err != nil {
		log.Printf("Invalid score info: %+v, err: %s\n", r.Body, err)
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	if err = c.app.AddScore(score); err != nil {
		log.Printf("Failed to save data: %v, err: %s\n", r, err)

		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(score)
	if err != nil {
		log.Printf("Failed to load marshal data: %v, err: %s\n", r, err)
	}
}

func (c *Controller) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["user_guid"]

	err := models.IsValidGuid(guid)
	if err != nil {
		log.Printf("Invalid GUID: err: %s\n", err)
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	rs, err := c.app.GetUserById(guid)
	if err != nil {
		if err == models.ErrNotFound {
			log.Printf("User not found: err: %s\n", err)
			http.Error(w,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound)
			return
		}
		log.Printf("Failed to load user detail: err: %s\n", err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		log.Printf("Failed to load marshal data: %v, err: %s\n", rs, err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Failed to decode json form: %+v, err: %s\n", r.Body, err)
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	if err := user.IsValid(); err != nil {
		log.Printf("Invalid user info: %+v, err: %s\n", r.Body, err)
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	if err = c.app.CreateUser(user); err != nil {
		log.Printf("Failed to save data: %v, err: %s\n", r, err)

		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Failed to load marshal data: %v, err: %s\n", r, err)
	}
}
