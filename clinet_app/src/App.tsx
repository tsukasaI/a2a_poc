import type { Component } from 'solid-js';
import { Route, Router, useSearchParams } from '@solidjs/router';
import styles from './App.module.css';
import { generateCodeVerifier, OAuth2Client, OAuth2Fetch } from '@badgateway/oauth2-client';
import { getCookie, setCookie } from 'typescript-cookie'
import axios from 'axios';

const clientId = "d03a0fa9-3059-4d16-bf51-1decc9bae9b0";
type CodeQuery = { code: string }

let codeVerifier = ""

const client = new OAuth2Client({
  server: 'http://localhost:4444/',
  clientId: clientId,
});

const Callback: Component = () => {
  const [searchParams, _] = useSearchParams<CodeQuery>()
  let r
  let e

  const exchangeCode = () => {
    axios.post("http://localhost:4444/oauth2/token", {
      "grant_type": "authorization_code",
      "code": searchParams.code,
      "redirect_uri": "http://localhost:3002/callback",
      "code_verifier": codeVerifier,
      "client_id": clientId,
    }).then(res => { console.log(res); r = res }).catch(er => { console.log(er); e = er })
  }
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <h1>Client App -Callback-</h1>
      </header>
      <p>code: {searchParams.code}</p>
      <button onClick={exchangeCode}>Exchange</button>
      <p>r: {r}</p>
      <p>e: {e}</p>
    </div>
  );
};


const Home: Component = () => {
  const onClickButton = () => {
    window.open(`http://localhost:4444/oauth2/auth?client_id=${clientId}&redirect_uri=http%3A%2F%2Flocalhost%3A3002%2Fcallback&response_type=code&state=aaaaaaaaaaaaaaaaaa&code_challenge_method=S256&code_challenge=${codeVerifier}`)
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
      console.log(codeVerifier)
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
