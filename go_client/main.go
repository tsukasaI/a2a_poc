package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"context"

	"golang.org/x/oauth2"
)

const hydraAdminURL = "http://localhost:4445"

var cookieStore []*http.Cookie
var verifier = oauth2.GenerateVerifier()
var codeChallenge = oauth2.S256ChallengeOption(verifier)

var (
	loginTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
</head>
<body>
	<h2>SUPAY APP LOGIN</h2>
	<p>The client app will redirect to the specific page in supay APP</p><br>
	<p>Login validation will be done by supay side. If it can validate user with refresh and access token, that's fine too</p><br>
	<p>This form will do an HTTP PUT request to /oauth2/auth/requests/login/accept?login_challenge=</p><br>
	<p>when the user authentication is completed.</p><br>
    <form action="/login" method="post">
        <input type="hidden" name="login_challenge" value="{{.LoginChallenge}}" />
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" />
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" />
        <button type="submit">Login</button>
    </form>
</body>
</html>
`
	consentTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Consent</title>
</head>
<body>
	<h2>SUPAY APP CONSENT PAGE</h2>
	<p>supay app will receive a redirect_to URL. Either use deep link, or redirect internally somehow.</p><br>
	<p>when the user consents, the app will redirect back to the client app with the ID token, refresh and access token</p>
    <form action="/consent" method="post">
        <input type="hidden" name="consent_challenge" value="{{.ConsentChallenge}}" />
        <h2>{{.ClientName}} wants to access your account</h2>
        <p>Requested scopes: {{.RequestedScopes}}</p>
        <button type="submit" name="consent" value="accept">Allow</button>
        <button type="submit" name="consent" value="reject">Deny</button>
    </form>
</body>
</html>
`
)

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/consent", consentHandler)
	port := ":3001"
	log.Println("verifier: ", verifier)
	log.Println("oauth2.S256ChallengeOption(verifier): ", oauth2.S256ChallengeOption(verifier))
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

		req, _ := http.NewRequest(http.MethodGet, deeplink, nil)

		getToken(req.FormValue("code"))

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

func getToken(code string) {
	log.Println("code: ", code)
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID: "3a121f8a-9802-4efc-b23d-214a90cda035",
		// ClientSecret: "YOUR_CLIENT_SECRET",
		Scopes: []string{"offline_access", "offline", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:4444/oauth2/auth",
			TokenURL: "http://localhost:4444/oauth2/token",
		},
		RedirectURL: "testDeepLink://mobile",
	}

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Println("Exchange err: ", err)
	}

	log.Printf("tok: %+v\n\n", tok)
}
