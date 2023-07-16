package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	valorant "github.com/iCodeOfTruth/go-rso"
)

type UserInfoResp struct {
	UserId string `json:"sub"`
}

type EntitlementsResp struct {
	EntitlementsToken string `json:"entitlements_token"`
}

type StoreResponse struct {
	SkinsPanelLayout struct {
		SingleItemOffers []string `json:"SingleItemOffers"`
	} `json:"SkinsPanelLayout"`
}

type WeaponSkin struct {
	Data struct {
		DisplayName string `json:"displayName"`
	} `json:"data"`
}

func main() {
	valorant.RiotUserAgent = "RiotClient/62.0.1.4852117.4789131 rso-auth (Windows;11;;Professional, x64)" // Set your own user agent

	client := valorant.New(nil)
	var username string
	var password string
	var code string

	fmt.Println("Input your username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		return
	}
	fmt.Println("Input your password: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		return
	}

	// Authorize
	data, err := client.Authorize(username, password)
	if err == valorant.ErrorRiotMultifactor {
		fmt.Println("Input your multi-factor code: ")
		_, err := fmt.Scanln(&code)
		if err != nil {
			return
		}

		data, err = client.SubmitTwoFactor(code)
	} else if err != nil {
		panic(err)
	}

	// Get get user id
	req, err := http.NewRequest(http.MethodGet, "https://auth.riotgames.com/userinfo", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body := new(UserInfoResp)
	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		return
	}
	userId := body.UserId

	// Get entitlements token
	req, err = http.NewRequest(http.MethodPost, "https://entitlements.auth.riotgames.com/api/token/v1", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	entitlementsBody := new(EntitlementsResp)
	err = json.NewDecoder(resp.Body).Decode(entitlementsBody)
	if err != nil {
		return
	}
	entitlementsToken := entitlementsBody.EntitlementsToken

	// Get store
	req, err = http.NewRequest(http.MethodGet, "https://pd.na.a.pvp.net/store/v2/storefront/"+userId, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.AccessToken))
	req.Header.Set("X-Riot-Entitlements-JWT", entitlementsToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	storeBody := new(StoreResponse)
	err = json.NewDecoder(resp.Body).Decode(storeBody)
	if err != nil {
		return
	}

	// Get weapon skin display names and print them
	for _, offer := range storeBody.SkinsPanelLayout.SingleItemOffers {
		req, err = http.NewRequest(http.MethodGet, "https://valorant-api.com/v1/weapons/skinlevels/"+offer, nil)
		if err != nil {
			panic(err)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		weaponSkin := new(WeaponSkin)
		err = json.NewDecoder(resp.Body).Decode(weaponSkin)
		resp.Body.Close()
		if err != nil {
			return
		}
		fmt.Println(weaponSkin.Data.DisplayName)
	}

	fmt.Println("Press the enter key to exit...")
	_, _ = fmt.Scanln()

}
