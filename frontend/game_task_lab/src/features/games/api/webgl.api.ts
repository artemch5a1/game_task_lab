import { invoke } from "@tauri-apps/api/core";

export type WebglStatus = {
  running: boolean;
  url: string | null;
  port: number | null;
};

export const webglApi = {
  start: () => invoke<WebglStatus>("webgl_start"),
  stop: () => invoke<WebglStatus>("webgl_stop"),
  status: () => invoke<WebglStatus>("webgl_status"),
};

