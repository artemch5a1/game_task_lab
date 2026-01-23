// features/games/stores/game.store.ts
import { createStore } from 'solid-js/store';
import type {CreateGameDto, GameDto} from '../types/game.types.ts';
import { gameApi } from '../api/game.api';

interface GameState {
    games: GameDto[];
    data: GameDto[];
    isLoading: boolean;
    error: string | null;
}

const initialState: GameState = {
    games: [],
    data:[],
    isLoading: false,
    error: null,
};

export type GameStore = {
    state: GameState;
    actions: {
        loadGames: () => Promise<void>;
        createGame: (dto: CreateGameDto) => Promise<GameDto>;
        deleteGame: (id: string) => Promise<void>;
        filterGames: (search: string) => Promise<void>
    };
};

export const createGameStore  = () : GameStore => {
    const [state, setState] = createStore(initialState);

    const actions = {
        async loadGames() {
            setState('isLoading', true);
            setState('error', null);

            try {
                const games = await gameApi.getAllGames();
                setState('data', games)
                setState('games', games);
            } catch (error) {
                console.error('[loadGames] error', error);
                setState('error', error instanceof Error ? error.message : 'Unknown error');
            } finally {
                setState('isLoading', false);
            }
        },

        async createGame(dto: CreateGameDto) {
            setState('isLoading', true);

            try {
                const newGame = await gameApi.createGame(dto);
                setState('games', [...state.games, newGame]);
                setState('games', [...state.data, newGame]);
                return newGame;
            } catch (error) {
                setState('error', error instanceof Error ? error.message : 'Failed to create game');
                throw error;
            } finally {
                setState('isLoading', false);
            }
        },

        async deleteGame(id: string) {
            setState('isLoading', true);

            try {
                await gameApi.deleteGame(id);
                setState('games', state.games.filter(game => game.id !== id));
                setState('games', state.data.filter(game => game.id !== id));
            } catch (error) {
                setState('error', error instanceof Error ? error.message : 'Failed to delete game');
                throw error;
            } finally {
                setState('isLoading', false);
            }
        },

        async filterGames(search: string): Promise<void> {
            const games = await gameApi.getAllGames();

            const gamesFiltered = games.filter(game => game.title.includes(search));

            setState('data', games)
            setState('games', gamesFiltered);
        }
    };

    return { state, actions };
};



// Глобальный store
export const gameStore = createGameStore();