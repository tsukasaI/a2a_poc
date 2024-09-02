package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const hydraAdminURL = "http://localhost:4445"

var cookieStore []*http.Cookie

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/consent", consentHandler)
	port := ":3030"
	log.Printf("Login & Consent app listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == http.MethodGet {
		log.Println("loginHandler get called")

		cookieStore = append(make([]*http.Cookie, 0), r.Cookies()...)
		log.Printf("cookieStore: %+v\n", cookieStore)

		loginChallenge := r.FormValue("login_challenge")
		redirectURL := fmt.Sprintf("testDeepLink://mobile/login?login_challenge=%s", loginChallenge)
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}

	if r.Method == http.MethodPost {
		log.Println("loginHandler post called")
		loginChallenge := r.FormValue("login_challenge")
		username := r.FormValue("username")
		redirectURL, err := acceptLoginRequest(loginChallenge, username)
		if err != nil {
			log.Printf("acceptLoginRequest failed.; %v", err)
			http.Error(w, "Failed to accept login request", http.StatusInternalServerError)
			return
		}
		consntChallengeRedirectUrl, err := getConsentChallengeRedirect(redirectURL)
		if err != nil {
			log.Printf("getConsentChallengeRedirect failed.; %v", err)
			http.Error(w, "Failed to getConsentChallengeRedirect", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, consntChallengeRedirectUrl, http.StatusFound)
		return
	}
}

func getConsentChallengeRedirect(redirectURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, redirectURL, nil)
	if err != nil {
		log.Printf("http.NewRequest failed.; %v", err)
		return "", err
	}
	noRedirectClient := new(http.Client)
	noRedirectClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	setAllCookies(req)
	redirectResp, err := noRedirectClient.Do(req)
	if err != nil {
		log.Printf("redirectURL get failed: %v", err)
		return "", err
	}
	defer redirectResp.Body.Close()

	loc, err := redirectResp.Location()
	if err != nil {
		log.Printf("redirectResp.Location() failed: %v", err)
		return "", err
	}
	cookieStore = append(cookieStore, redirectResp.Cookies()...)

	return loc.String(), nil
}

func setAllCookies(r *http.Request) {
	for _, c := range cookieStore {
		r.AddCookie(c)
	}
}

func acceptLoginRequest(loginChallenge, subject string) (string, error) {
	url := fmt.Sprintf("%s/oauth2/auth/requests/login/accept?login_challenge=%s", hydraAdminURL, loginChallenge)
	data := map[string]string{"subject": subject}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to accept login request: %s", resp.Status)
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", err
	}

	redirectTo, ok := responseData["redirect_to"].(string)
	if !ok {
		return "", fmt.Errorf("missing redirect_to field in response")
	}

	return redirectTo, nil
}

func consentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("POST consentHandler called")
		consentChallenge := r.FormValue("consent_challenge")
		consent := r.FormValue("consent")

		var redirectURL string
		var err error
		if consent == "accept" {
			redirectURL, err = acceptConsentRequest(consentChallenge)
		} else {
			redirectURL, err = rejectConsentRequest(consentChallenge)
		}

		if err != nil {
			http.Error(w, "Failed to process consent request", http.StatusInternalServerError)
			return
		}

		deeplink, err := getCodeRequest(redirectURL)
		if err != nil {
			log.Printf("failed getCodeRequest: %+v", err)
			http.Error(w, "Failed to process get code request", http.StatusInternalServerError)
			return
		}
		fmt.Println(deeplink)

		http.Redirect(w, r, deeplink, http.StatusFound)
	}
}

func acceptConsentRequest(consentChallenge string) (string, error) {
	url := fmt.Sprintf("%s/oauth2/auth/requests/consent/accept?consent_challenge=%s", hydraAdminURL, consentChallenge)
	data := map[string]interface{}{
		"grant_scope":  []string{"openid", "offline"},
		"remember":     true,
		"remember_for": 3600,
	}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to accept consent request: %s", resp.Status)
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", err
	}

	redirectTo, ok := responseData["redirect_to"].(string)
	if !ok {
		return "", fmt.Errorf("missing redirect_to field in response")
	}
	return redirectTo, nil
}

func rejectConsentRequest(consentChallenge string) (string, error) {
	url := fmt.Sprintf("%s/oauth2/auth/requests/consent/reject?consent_challenge=%s", hydraAdminURL, consentChallenge)
	data := map[string]interface{}{
		"error":             "access_denied",
		"error_description": "The resource owner denied the request",
	}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to reject consent request: %s", resp.Status)
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", err
	}

	redirectTo, ok := responseData["redirect_to"].(string)
	if !ok {
		return "", fmt.Errorf("missing redirect_to field in response")
	}

	return redirectTo, nil
}

func getCodeRequest(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("New Request failed: %+v\n", err)
		return "", err
	}
	setAllCookies(req)

	noRedirectClient := new(http.Client)
	noRedirectClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	res, err := noRedirectClient.Do(req)
	if err != nil {
		log.Printf("fialed noRedirectClient.Do: %+v\n", err)
		return "", err
	}
	defer res.Body.Close()

	redirectTo, ok := res.Header["Location"]
	if !ok {
		return "", fmt.Errorf("missing redirect_to field in response")
	}
	return redirectTo[0], nil
}
