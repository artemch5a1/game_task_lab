import { createStore } from "solid-js/store";
import { authApi } from "../api/auth.api";

const AUTH_TOKEN_KEY = "auth_token";

interface AuthState {
  token: string | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: AuthState = {
  token: localStorage.getItem(AUTH_TOKEN_KEY),
  isLoading: false,
  error: null,
};

export type AuthStore = {
  state: AuthState;
  actions: {
    login: (username: string, password: string) => Promise<void>;
    logout: () => void;
    clearError: () => void;
    isAuthenticated: () => boolean;
  };
};

export const createAuthStore = (): AuthStore => {
  const [state, setState] = createStore<AuthState>(initialState);

  const actions = {
    async login(username: string, password: string) {
      setState("isLoading", true);
      setState("error", null);

      try {
        const { token } = await authApi.login({ username, password });
        localStorage.setItem(AUTH_TOKEN_KEY, token);
        setState("token", token);
      } catch (e) {
        const msg = e instanceof Error ? e.message : "Login failed";
        setState("error", msg);
        throw e;
      } finally {
        setState("isLoading", false);
      }
    },
    logout() {
      localStorage.removeItem(AUTH_TOKEN_KEY);
      setState("token", null);
      setState("error", null);
    },
    clearError() {
      setState("error", null);
    },
    isAuthenticated() {
      return !!state.token;
    },
  };

  return { state, actions };
};

export const authStore = createAuthStore();

