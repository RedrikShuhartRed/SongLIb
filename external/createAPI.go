package external

import (
	"encoding/json"
	"net/http"

	"github.com/RedrikShuhartRed/EfMobSongLib/config"
)

func SongHandler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	if group == "" || song == "" {
		http.Error(w, "Missing group or song parameters", http.StatusBadRequest)
		return
	}

	details := DataAPI{
		ReleaseDate: "16.07.2015",
		VersesEn: `Ooh baby, don't you know I suffer?
Ooh baby, can you hear me moan?
You caught me under false pretenses
How long before you let me go?
Ooh
You set my soul alight
Ooh
You set my soul alight
Glaciers melting in the dead of night (ooh)
And the superstars sucked into the supermassive (you set my soul alight)
Glaciers melting in the dead of night
And the superstars sucked into the (you set my soul)
(Into the supermassive)
I thought I was a fool for no one
Ooh baby, I'm a fool for you
You're the queen of the superficial
And how long before you tell the truth?
Ooh
You set my soul alight
Ooh
You set my soul alight
Glaciers melting in the dead of night (ooh)
And the superstars sucked into the supermassive (you set my soul alight)
Glaciers melting in the dead of night
And the superstars sucked into the (you set my soul)
(Into the supermassive)
Supermassive black hole
Supermassive black hole
Supermassive black hole
Supermassive black hole
Glaciers melting in the dead of night
And the superstars sucked into the supermassive
Glaciers melting in the dead of night
And the superstars sucked into the supermassive
Glaciers melting in the dead of night (ooh)
And the superstars sucked into the supermassive (you set my soul alight)
Glaciers melting in the dead of night
And the superstars sucked into the (you set my soul)
(Into the supermassive)
Supermassive black hole
Supermassive black hole
Supermassive black hole
Supermassive black hole`,
		VersesRu: `Ох, детка, знаешь ли ты, что я страдаю?
Ох, детка, можешь ли ты услышать мой стон?
Ты поймала меня на ложном притворстве,
Когда же ты меня отпустишь?
Ох... 
Ты поджигаешь мою душу
Ох... 
Ты поджигаешь мою душу
(Ледники тают в глухую ночь,
И звезды засасываются массами)
Ох... Ты поджигаешь мою душу
(Ледники тают в глухую ночь
И звезды засасываются в...)
Поджигаешь душу...
Я думал, что не сходил с ума ни по кому,
Ох, детка, я без ума от тебя.
Ты королева всего искусственного,
И когда же ты мне скажешь правду?
Ох... 
Ты поджигаешь мою душу
Ох... 
Ты поджигаешь мою душу
(Ледники тают в глухую ночь
И звезды засасываются массами)
Ох, ты поджигаешь мою душу
(Ледники тают в глухую ночь
И звезды засасываются в...)
Поджигаешь душу...
Супермассивная черная дыра
Супермассивная черная дыра
Супермассивная черная дыра
Супермассивная черная дыра
Ледники тают в глухую ночь
И звезды засасываются массами.
Ледники тают в глухую ночь
И звезды засасываются массами.
(Ледники тают в глухую ночь,
И звезды засасываются массами)
Ох... Ты поджигаешь мою душу
(Ледники тают в глухую ночь
И звезды засасываются в...)
Поджигаешь душу...
Супермассивная черная дыра
Супермассивная черная дыра
Супермассивная черная дыра
Супермассивная черная дыра`,
		Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartServer(cfg *config.Config) error {
	http.HandleFunc("/info", SongHandler)
	return http.ListenAndServe(":"+cfg.ExternalPort, nil)
}
