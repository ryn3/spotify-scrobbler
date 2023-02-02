// This example demonstrates how to authenticate with Spotify using the
// client credentials flow.  Note that this flow does not include authorization
// and can't be used to access a user's private data.
//
// Make sure you set the SPOTIFY_ID and SPOTIFY_SECRET environment variables
// prior to running this example.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zmb3/spotify"
	spotifyauth "github.com/zmb3/spotify/auth"
)

var (
	redirectURI = "http://localhost:8080/scrobble"
	auth        = spotifyauth.New(
		spotifyauth.WithClientID("f8a8019a95a84c9490936c08f868e379"),
		spotifyauth.WithClientSecret("506ca99490684833bac7f7df9ae7bd35"),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadRecentlyPlayed),
	)
	state = "abc123"
)

func main() {
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:\n", url)
	http.HandleFunc("/scrobble", scrobble)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func scrobble(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	client := spotify.New(auth.Client(r.Context(), tok))

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	var tid string
	for {
		listOfTracks, err := client.PlayerRecentlyPlayed(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		firstTracksID := listOfTracks[0].Track.ID.String()

		if tid != firstTracksID {
			tid = firstTracksID
			log.Println(listOfTracks[0].Track.String())
		}
		time.Sleep(10 * time.Second)
	}
}
