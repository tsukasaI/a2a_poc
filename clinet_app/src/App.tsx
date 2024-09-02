import { createSignal, type Component } from 'solid-js';
import { Route, Router, useSearchParams } from '@solidjs/router';
import styles from './App.module.css';
import { generateCodeVerifier, OAuth2Client, OAuth2Fetch } from '@badgateway/oauth2-client';
import { getCookie, setCookie } from 'typescript-cookie'

const clientId = "4d9b630b-abfa-4aaf-82c7-9a50518cb68b";
type CodeQuery = { code: string }

let codeVerifier = ""

const client = new OAuth2Client({
  server: 'http://localhost:4444/',
  clientId: clientId,
  tokenEndpoint: '/oauth2/token',
  authorizationEndpoint: '/oauth2/auth',
});

const Callback: Component = () => {
  const [searchParams, _] = useSearchParams<CodeQuery>()
  const [accessToken, setAccessToken] = createSignal("")
  const [refreshToken, setRefreshToken] = createSignal("")
  const [expiresAt, setExpiresAt] = createSignal(0)

  const exchangeCode = () => {
    client.authorizationCode.getTokenFromCodeRedirect(
      document.location.href,
      {
        redirectUri: "http://localhost:3002/callback",
        state: "aaaaaaaaaaaaaaaaaa",
        codeVerifier: codeVerifier,
      })

      .then(res => {
        setAccessToken(res.accessToken)
        if (res.refreshToken) { setRefreshToken(res.refreshToken) }
        if (res.expiresAt) { setExpiresAt(res.expiresAt) }

      }).catch(e => { console.error(e) })
  }
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <h1>Client App -Callback-</h1>
      </header>
      <p>code: {searchParams.code}</p>
      <button onClick={exchangeCode}>Exchange</button>

      <div>
      {accessToken() !== "" && <p>access_token: {accessToken()}</p>}
      {refreshToken() !== "" && <p>access_token: {refreshToken()}</p>}
      {expiresAt() !== 0 && <p>access_token: {expiresAt()}</p>}
      </div>
    </div>
  );
};


const Home: Component = () => {
  const onClickButton = () => {
    client.authorizationCode.getAuthorizeUri({

      // URL in the app that the user should get redirected to after authenticating
      redirectUri: 'http://localhost:3002/callback',

      // Optional string that can be sent along to the auth server. This value will
      // be sent along with the redirect back to the app verbatim.
      state: 'aaaaaaaaaaaaaaaaaa',

      codeVerifier,

      // scope: ['scope1', 'scope2'],

    }).then((res) => { document.location = res });
  }
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <h1>Client App</h1>
      </header>
      <div>
        <button onClick={onClickButton}>Get token</button>
      </div>
    </div>
  );
};

const setCodeVerifier = () => {
  const codeVerifierCookie = getCookie("codeVerifier")
  if (typeof codeVerifierCookie === "undefined") {
    generateCodeVerifier().then((res) => {
      codeVerifier = res
      setCookie("codeVerifier", codeVerifier)
    });
  } else {
    codeVerifier = codeVerifierCookie
  }
}

const App: Component = () => {
  setCodeVerifier()
  return (
    <Router>
      <Route path="/" component={Home} />
      <Route path="/callback" component={Callback} />
    </Router>
  );
};

export default App;
