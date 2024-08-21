import type { Component } from 'solid-js';

import styles from './App.module.css';

const App: Component = () => {
  const onClickButton = () => {
    console.log("onClickButton called")
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

export default App;
