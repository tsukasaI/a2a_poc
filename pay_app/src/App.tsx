import type { Component } from 'solid-js';
import { Route, Router } from '@solidjs/router';

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

const Callback: Component = () => {
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <h1>Callback</h1>
      </header>
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
