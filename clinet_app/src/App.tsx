import type { Component } from 'solid-js';
import { Route, Router, useSearchParams } from '@solidjs/router';
import styles from './App.module.css';

const clientId = "f4d402cd-a3b9-446b-857f-4411ac2bde71";
type CodeQuery = {code: string}

const Callback: Component = () => {
  const [searchParams, _] = useSearchParams<CodeQuery>()

  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <h1>Callback</h1>
        <p>code: {searchParams.code}</p>
      </header>
    </div>
  );
};


const Home: Component = () => {
  const onClickButton = () => {
    window.open(`localhost:4444/oauth2/auth?client_id=${clientId}&redirect_uri=http%3A%2F%2Flocalhost%3A3001%2Fcallback&response_type=code&state=aaaaaaaaaaaaaaaaaa`)
  }
  return (
    <div class={styles.App}>
      <header class={styles.header}>
      <h1>Client App</h1>
      </header>
      <div>
        <p>開くボタン</p>
        <button onClick={onClickButton}>開く</button>
      </div>
    </div>
  );
};

const App: Component = () => {
  return (
    <Router>
      <Route path="/" component={Home} />
      <Route path="/callback" component={Callback} />
    </Router>
  );
};

export default App;
