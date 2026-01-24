import { createStore } from "solid-js/store";
import { authApi } from "../api/auth.api";

const AUTH_TOKEN_KEY = "auth_token";

type AuthRole = "admin" | "user" | string;

const decodeRoleFromToken = (token: string | null): AuthRole | null => {
  if (!token) return null;

  try {
    const parts = token.split(".");
    if (parts.length < 2) return null;

    const base64Url = parts[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const padded = base64.padEnd(Math.ceil(base64.length / 4) * 4, "=");
    const json = atob(padded);
    const payload = JSON.parse(json) as { role?: string };
    return payload.role ?? null;
  } catch {
    return null;
  }
};

interface AuthState {
  token: string | null;
  role: AuthRole | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: AuthState = {
  token: localStorage.getItem(AUTH_TOKEN_KEY),
  role: decodeRoleFromToken(localStorage.getItem(AUTH_TOKEN_KEY)),
  isLoading: false,
  error: null,
};

export type AuthStore = {
  state: AuthState;
  actions: {
    login: (username: string, password: string) => Promise<void>;
    register: (username: string, password: string) => Promise<void>;
    logout: () => void;
    clearError: () => void;
    isAuthenticated: () => boolean;
    isAdmin: () => boolean;
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
        setState("role", decodeRoleFromToken(token));
      } catch (e) {
        const msg = e instanceof Error ? e.message : "Login failed";
        setState("error", msg);
        throw e;
      } finally {
        setState("isLoading", false);
      }
    },
    async register(username: string, password: string) {
      setState("isLoading", true);
      setState("error", null);

      try {
        const { token } = await authApi.register({ username, password });
        localStorage.setItem(AUTH_TOKEN_KEY, token);
        setState("token", token);
        setState("role", decodeRoleFromToken(token));
      } catch (e) {
        const msg = e instanceof Error ? e.message : "Register failed";
        setState("error", msg);
        throw e;
      } finally {
        setState("isLoading", false);
      }
    },
    logout() {
      localStorage.removeItem(AUTH_TOKEN_KEY);
      setState("token", null);
      setState("role", null);
      setState("error", null);
    },
    clearError() {
      setState("error", null);
    },
    isAuthenticated() {
      return !!state.token;
    },
    isAdmin() {
      return state.role === "admin";
    },
  };

  return { state, actions };
};

export const authStore = createAuthStore();

