import "./LoginPage.css";
import { createSignal, Show } from "solid-js";
import { A, useNavigate } from "@solidjs/router";
import { authStore } from "../store/auth.store";

export const RegisterPage = () => {
  const navigate = useNavigate();
  const { state, actions } = authStore;

  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [password2, setPassword2] = createSignal("");

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    actions.clearError();

    if (password() !== password2()) {
      // локальная валидация
      throw new Error("Passwords do not match");
    }

    await actions.register(username(), password());
    navigate("/", { replace: true });
  };

  return (
    <div class="login-page">
      <div class="login-card">
        <h1 class="login-title">Create Account</h1>
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
            <label>Password</label>
            <input
              class="login-input"
              type="password"
              value={password()}
              onInput={(e) => setPassword(e.currentTarget.value)}
              disabled={state.isLoading}
              required
            />
          </div>

          <div class="login-field">
            <label>Repeat password</label>
            <input
              class="login-input"
              type="password"
              value={password2()}
              onInput={(e) => setPassword2(e.currentTarget.value)}
              disabled={state.isLoading}
              required
            />
          </div>

          <Show when={state.error}>
            <div class="login-error">{state.error}</div>
          </Show>

          <button class="login-button" type="submit" disabled={state.isLoading}>
            {state.isLoading ? "Creating..." : "Create Account"}
          </button>
        </form>

        <div class="login-footer">
          Already have an account? <A href="/login">Sign In</A>
        </div>
      </div>
    </div>
  );
};

