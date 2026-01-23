import "./App.css";
import { GameList } from "./features/games/components/GameList";

function App() {
  return (
    <main class="container">
      <h1>Игры</h1>
      <GameList />
    </main>
  );
}

export default App;
