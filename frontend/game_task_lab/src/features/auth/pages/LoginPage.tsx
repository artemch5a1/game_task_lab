import "./LoginPage.css";
import { createSignal, Show } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { authStore } from "../store/auth.store";

export const LoginPage = () => {
  const navigate = useNavigate();
  const { state, actions } = authStore;

  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [remember, setRemember] = createSignal(true);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    actions.clearError();

    // "remember" пока влияет только на UX (чекбокс из макета).
    // Токен всегда хранится в localStorage, как ты просил.
    void remember();

    await actions.login(username(), password());
    navigate("/", { replace: true });
  };

  return (
    <div class="login-page">
      <div class="login-card">
        <h1 class="login-title">Login to Account</h1>
        <div class="login-subtitle">Please enter your email and password to continue</div>

        <form class="login-form" onSubmit={handleSubmit}>
          <div class="login-field">
            <label>Email address:</label>
            <input
              class="login-input"
              type="text"
              value={username()}
              onInput={(e) => setUsername(e.currentTarget.value)}
              placeholder="your@email.com"
              disabled={state.isLoading}
              required
            />
          </div>

          <div class="login-field">
            <div class="login-row">
              <label style={{ margin: 0 }}>Password</label>
              <a class="login-link" href="#" onClick={(e) => e.preventDefault()}>
                Forget Password?
              </a>
            </div>
            <input
              class="login-input"
              type="password"
              value={password()}
              onInput={(e) => setPassword(e.currentTarget.value)}
              disabled={state.isLoading}
              required
            />
          </div>

          <div class="login-row">
            <label class="login-remember">
              <input
                type="checkbox"
                checked={remember()}
                onChange={(e) => setRemember(e.currentTarget.checked)}
                disabled={state.isLoading}
              />
              Remember Password
            </label>
          </div>

          <Show when={state.error}>
            <div class="login-error">{state.error}</div>
          </Show>

          <button class="login-button" type="submit" disabled={state.isLoading}>
            {state.isLoading ? "Signing In..." : "Sign In"}
          </button>
        </form>

        <div class="login-footer">
          Don’t have an account?{" "}
          <a href="#" onClick={(e) => e.preventDefault()}>
            Create Account
          </a>
        </div>
      </div>
    </div>
  );
};

