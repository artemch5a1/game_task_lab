// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

use serde::Serialize;
use std::{
    path::PathBuf,
    process::{Child, Command, Stdio},
    sync::Mutex,
    time::Duration,
};
use tauri::State;

const UNITY_EXECUTABLE_PATH: &str =
    "/traffic_light_alarm_system/Система сигнализации светофоров.x86_64";

#[derive(Default)]
struct UnityProcessState {
    child: Mutex<Option<Child>>,
}

#[derive(Serialize)]
struct UnityStatus {
    running: bool,
    pid: Option<u32>,
}

fn status_from_child(child: &mut Child) -> UnityStatus {
    match child.try_wait() {
        Ok(Some(_)) => UnityStatus {
            running: false,
            pid: None,
        },
        Ok(None) => UnityStatus {
            running: true,
            pid: Some(child.id()),
        },
        Err(_) => UnityStatus {
            running: true,
            pid: Some(child.id()),
        },
    }
}

#[tauri::command]
fn unity_status(state: State<UnityProcessState>) -> UnityStatus {
    let mut guard = state.child.lock().expect("unity process state poisoned");
    if let Some(child) = guard.as_mut() {
        let status = status_from_child(child);
        if !status.running {
            *guard = None;
        }
        status
    } else {
        UnityStatus {
            running: false,
            pid: None,
        }
    }
}

#[tauri::command]
fn unity_start(state: State<UnityProcessState>, executable_path: Option<String>) -> Result<UnityStatus, String> {
    let mut guard = state.child.lock().map_err(|_| "state lock failed".to_string())?;

    if let Some(child) = guard.as_mut() {
        let status = status_from_child(child);
        if status.running {
            return Ok(status);
        }
        *guard = None;
    }

    let exec_path: PathBuf = executable_path
        .map(PathBuf::from)
        .unwrap_or_else(|| PathBuf::from(UNITY_EXECUTABLE_PATH));

    if !exec_path.exists() {
        return Err(format!(
            "Unity executable not found at path: {}",
            exec_path.to_string_lossy()
        ));
    }
    let exec_dir = exec_path
        .parent()
        .ok_or_else(|| "failed to determine Unity executable directory".to_string())?;

    let mut cmd = Command::new(&exec_path);
    cmd.current_dir(exec_dir)
        .stdin(Stdio::null())
        .stdout(Stdio::null())
        .stderr(Stdio::null());

    #[cfg(unix)]
    {
        use std::os::unix::process::CommandExt;
        unsafe {
            cmd.pre_exec(|| {
                // create new process group (pgid == pid)
                if libc::setpgid(0, 0) != 0 {
                    return Err(std::io::Error::last_os_error());
                }
                Ok(())
            });
        }
    }

    let child = cmd.spawn().map_err(|e| format!("failed to start Unity: {e}"))?;
    let pid = child.id();
    *guard = Some(child);
    Ok(UnityStatus {
        running: true,
        pid: Some(pid),
    })
}

#[tauri::command]
fn unity_stop(state: State<UnityProcessState>) -> Result<UnityStatus, String> {
    let child_opt = {
        let mut guard = state.child.lock().map_err(|_| "state lock failed".to_string())?;
        guard.take()
    };

    let Some(mut child) = child_opt else {
        return Ok(UnityStatus {
            running: false,
            pid: None,
        });
    };

    let pid = child.id();

    // Try a graceful stop (TERM), then force (KILL).
    #[cfg(unix)]
    {
        unsafe {
            if pid != 0 {
                // negative pid => process group
                let _ = libc::kill(-(pid as i32), libc::SIGTERM);
            }
        }
        std::thread::sleep(Duration::from_millis(600));
        if let Ok(None) = child.try_wait() {
            unsafe {
                if pid != 0 {
                    let _ = libc::kill(-(pid as i32), libc::SIGKILL);
                }
            }
        }
    }

    // Fallback for non-unix (or if group killing didn't work)
    let _ = child.kill();
    let _ = child.wait();

    Ok(UnityStatus {
        running: false,
        pid: None,
    })
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .manage(UnityProcessState::default())
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![
            greet,
            unity_start,
            unity_stop,
            unity_status
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
