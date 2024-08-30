import type { Component } from 'solid-js';
import { Route, Router, useSearchParams } from '@solidjs/router';
import styles from './App.module.css';

const Home: Component = () => {
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <h1>Payment App</h1>
      </header>
    </div>
  );
};

type LoginChallengeQuery = { login_challenge: string }
type ConsentChallengeQuery = { consent_challenge: string }

const Login: Component = () => {
  const [searchParams, _] = useSearchParams<LoginChallengeQuery>()
  console.log(searchParams.login_challenge)
  return <div class={styles.App}>
    <header class={styles.header}>
      <h1>Payment App -Login-</h1>
    </header>
    <div>
      <form action="http://localhost:3030/login" method="post">
        <input type="hidden" name="login_challenge" value={searchParams.login_challenge} />
        <button type="submit">Login</button>
      </form>
    </div>
  </div>
}

const Consent: Component = () => {
  const [searchParams, _] = useSearchParams<ConsentChallengeQuery>()
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <h1>Payment App</h1>
      </header>
      <p>consent page</p>
      <form action="http://localhost:3030/consent" method="post">
        <input type="hidden" name="consent_challenge" value={searchParams.consent_challenge} />
        <button type="submit" name="consent" value="accept">Allow</button>
        <button type="submit" name="consent" value="reject">Deny</button>
      </form>
    </div>
  );
};

const App: Component = () => {
  return (
    <Router>
      <Route path="/" component={Home} />
      <Route path="/login" component={Login} />
      <Route path="/consent" component={Consent} />
    </Router>
  );
};

export default App;
