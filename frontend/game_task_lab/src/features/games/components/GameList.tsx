// features/games/components/GameList.tsx
import { createEffect } from 'solid-js';
import { gameStore } from '../store/game.store';

export const GameList = () => {
    const { state, actions } = gameStore;

    createEffect(() => {
        actions.loadGames();
    });

    return (
        <div>
            {state.isLoading && <div>Loading...</div>}
            {state.error && <div class="text-red-500">Error: {state.error}</div>}

            <ul>
                {state.games.map(game => (
                    <li id={game.id}>
                        <h3>{game.title}</h3>
                        <p>{game.description}</p>
                        <button
                            onClick={() => actions.deleteGame(game.id)}
                            disabled={state.isLoading}
                        >
                            Delete
                        </button>
                    </li>
                ))}
            </ul>
        </div>
    );
};