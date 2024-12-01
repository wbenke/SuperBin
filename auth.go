/*
This file is part of GigaPaste.

GigaPaste is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

GigaPaste is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with GigaPaste. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

type Session struct {
	Timer         *time.Timer
	SessionString string
}

var Sessions []Session = []Session{}

func AuthHandler(w http.ResponseWriter, r *http.Request) {

	var jsonData struct {
		Key string `json:"key"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&jsonData)
	if err != nil {
		fmt.Println(err)
		return
	}

	if jsonData.Key == Global.Password {

		randomBytes := make([]byte, 32)
		_, err := rand.Read(randomBytes)
		if err != nil {
			fmt.Println(err)
		}

		randomBytesString := base64.StdEncoding.EncodeToString(randomBytes)

		timer := time.AfterFunc(24*time.Hour, func() {

			for i, v := range Sessions {

				if v.SessionString == randomBytesString {

					Sessions = slices.Delete(Sessions, i, i+1)

				}

			}

		})

		Sessions = append(Sessions, Session{timer, randomBytesString})
		http.SetCookie(w, &http.Cookie{

			Name:     "session",
			Value:    randomBytesString,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		_, err = io.WriteString(w, "done")
		if err != nil {
			fmt.Println(err)
		}

	} else {

		_, err = io.WriteString(w, "wrong")
		if err != nil {
			fmt.Println(err)
		}

	}

}

func ValidateSession(w http.ResponseWriter, r *http.Request) bool {

	if !Global.EnablePassword {

		return true

	}

	cookie, err := r.Cookie("session")
	if err != nil {

		return false

	}

	for _, v := range Sessions {

		if v.SessionString == cookie.Value {

			return true

		}

	}

	return false

}

func DeleteSession(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	for i, v := range Sessions {

		if v.SessionString == cookie.Value {

			v.Timer.Stop()
			Sessions = slices.Delete(Sessions, i, i+1)

		}

	}

}
