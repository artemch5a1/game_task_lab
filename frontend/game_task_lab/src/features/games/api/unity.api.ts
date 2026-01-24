import { invoke } from "@tauri-apps/api/core";

export type UnityStatus = {
  running: boolean;
  pid: number | null;
};

export const unityApi = {
  start: (executablePath?: string) =>
    invoke<UnityStatus>("unity_start", { executablePath }),
  stop: () => invoke<UnityStatus>("unity_stop"),
  status: () => invoke<UnityStatus>("unity_status"),
};

