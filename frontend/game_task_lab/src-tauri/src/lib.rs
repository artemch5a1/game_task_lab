// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

use serde::Serialize;
use std::{
    fs::File,
    path::{Path, PathBuf},
    sync::{
        atomic::{AtomicBool, Ordering},
        Arc, Mutex,
    },
    thread,
    time::Duration,
};
use tauri::{path::BaseDirectory, AppHandle, Manager, State};
use tiny_http::{Header, Method, Response, Server, StatusCode};

#[derive(Default)]
struct WebglServerState {
    handle: Mutex<Option<WebglServerHandle>>,
}

struct WebglServerHandle {
    port: u16,
    running: Arc<AtomicBool>,
    thread: Option<thread::JoinHandle<()>>,
}

#[derive(Serialize)]
struct WebglStatus {
    running: bool,
    url: Option<String>,
    port: Option<u16>,
}

fn webgl_status_from_handle(handle: &WebglServerHandle) -> WebglStatus {
    WebglStatus {
        running: handle.running.load(Ordering::Relaxed),
        url: Some(format!("http://127.0.0.1:{}/", handle.port)),
        port: Some(handle.port),
    }
}

fn percent_decode_path(input: &str) -> String {
    // Minimal percent-decoder for paths (e.g. "Breach%20It.loader.js").
    let bytes = input.as_bytes();
    let mut out = Vec::with_capacity(bytes.len());
    let mut i = 0;
    while i < bytes.len() {
        match bytes[i] {
            b'%' if i + 2 < bytes.len() => {
                let h1 = bytes[i + 1];
                let h2 = bytes[i + 2];
                let val = |c: u8| -> Option<u8> {
                    match c {
                        b'0'..=b'9' => Some(c - b'0'),
                        b'a'..=b'f' => Some(c - b'a' + 10),
                        b'A'..=b'F' => Some(c - b'A' + 10),
                        _ => None,
                    }
                };
                if let (Some(a), Some(b)) = (val(h1), val(h2)) {
                    out.push((a << 4) | b);
                    i += 3;
                    continue;
                }
                out.push(bytes[i]);
                i += 1;
            }
            b'+' => {
                // Some clients may encode spaces as '+'. Not typical for paths but harmless.
                out.push(b' ');
                i += 1;
            }
            b => {
                out.push(b);
                i += 1;
            }
        }
    }
    String::from_utf8_lossy(&out).to_string()
}

fn content_type_for(path_without_encoding: &str) -> &'static str {
    let lower = path_without_encoding.to_ascii_lowercase();
    if lower.ends_with(".wasm") {
        return "application/wasm";
    }
    if lower.ends_with(".js") {
        return "application/javascript; charset=utf-8";
    }
    if lower.ends_with(".data") {
        return "application/octet-stream";
    }
    if lower.ends_with(".json") {
        return "application/json; charset=utf-8";
    }
    if lower.ends_with(".css") {
        return "text/css; charset=utf-8";
    }
    if lower.ends_with(".html") || lower.ends_with(".htm") {
        return "text/html; charset=utf-8";
    }
    if lower.ends_with(".png") {
        return "image/png";
    }
    if lower.ends_with(".jpg") || lower.ends_with(".jpeg") {
        return "image/jpeg";
    }
    if lower.ends_with(".svg") {
        return "image/svg+xml";
    }
    if lower.ends_with(".ico") {
        return "image/x-icon";
    }
    "application/octet-stream"
}

fn safe_join(root: &Path, rel: &str) -> Option<PathBuf> {
    let mut out = root.to_path_buf();
    for comp in Path::new(rel).components() {
        match comp {
            std::path::Component::Normal(p) => out.push(p),
            std::path::Component::CurDir => {}
            // Block absolute paths and any ".." traversal.
            _ => return None,
        }
    }
    Some(out)
}

fn resolve_breach_it_dir(app: &AppHandle) -> Option<PathBuf> {
    // 1) Bundled resource: <resource_dir>/breach_it
    if let Ok(dir) = app.path().resolve("breach_it", BaseDirectory::Resource) {
        if dir.join("index.html").exists() {
            return Some(dir);
        }
    }

    // 2) Dev fallback: search up from current executable directory for "breach_it/index.html"
    if let Ok(exe) = std::env::current_exe() {
        if let Some(mut dir) = exe.parent().map(|p| p.to_path_buf()) {
            for _ in 0..10 {
                let candidate = dir.join("breach_it").join("index.html");
                if candidate.exists() {
                    return Some(dir.join("breach_it"));
                }
                if !dir.pop() {
                    break;
                }
            }
        }
    }

    // 3) Dev fallback: current working directory
    if let Ok(cwd) = std::env::current_dir() {
        let candidate = cwd.join("breach_it").join("index.html");
        if candidate.exists() {
            return Some(cwd.join("breach_it"));
        }
    }

    None
}

fn start_webgl_server(root: PathBuf) -> Result<WebglServerHandle, String> {
    let listener = std::net::TcpListener::bind("127.0.0.1:0")
        .map_err(|e| format!("failed to bind localhost port: {e}"))?;
    let port = listener
        .local_addr()
        .map_err(|e| format!("failed to get bound port: {e}"))?
        .port();
    let server =
        Server::from_listener(listener, None).map_err(|e| format!("failed to start server: {e}"))?;

    let running = Arc::new(AtomicBool::new(true));
    let running_thread = running.clone();

    let thread = thread::spawn(move || {
        while running_thread.load(Ordering::Relaxed) {
            let request: tiny_http::Request = match server.recv_timeout(Duration::from_millis(200)) {
                Ok(Some(r)) => r,
                Ok(None) => continue,
                Err(_) => continue,
            };

            // Only GET/HEAD are expected for static assets.
            if request.method() != &Method::Get && request.method() != &Method::Head {
                let _ = request.respond(
                    Response::empty(StatusCode(405))
                        .with_header(Header::from_bytes("Content-Type", "text/plain; charset=utf-8").unwrap()),
                );
                continue;
            }

            let url = request.url();
            let path_part = url.split('?').next().unwrap_or("/");
            let decoded = percent_decode_path(path_part);
            let rel = if decoded == "/" || decoded.is_empty() {
                "index.html".to_string()
            } else {
                decoded.trim_start_matches('/').to_string()
            };

            // Service endpoint: stop the WebGL server (best-effort) when user exits the game.
            if rel == "__exit" {
                let _ = request.respond(Response::empty(StatusCode(204)));
                running_thread.store(false, Ordering::Relaxed);
                continue;
            }

            let Some(fs_path) = safe_join(&root, &rel) else {
                let _ = request.respond(Response::empty(StatusCode(400)));
                continue;
            };

            // If directory, try index.html
            let fs_path = if fs_path.is_dir() {
                fs_path.join("index.html")
            } else {
                fs_path
            };

            let (content_encoding, content_type_path) = if let Some(p) = fs_path.to_str() {
                if p.ends_with(".gz") {
                    ("gzip", p.trim_end_matches(".gz"))
                } else if p.ends_with(".br") {
                    ("br", p.trim_end_matches(".br"))
                } else {
                    ("", p)
                }
            } else {
                ("", "")
            };

            let file = match File::open(&fs_path) {
                Ok(f) => f,
                Err(_) => {
                    let _ = request.respond(Response::empty(StatusCode(404)));
                    continue;
                }
            };

            let mut response = Response::from_file(file).with_status_code(StatusCode(200));
            if !content_type_path.is_empty() {
                let _ = Header::from_bytes("Content-Type", content_type_for(content_type_path))
                    .map(|h| response.add_header(h));
            }
            if !content_encoding.is_empty() {
                let _ = Header::from_bytes("Content-Encoding", content_encoding)
                    .map(|h| response.add_header(h));
            }
            let _ = Header::from_bytes("Cache-Control", "no-cache")
                .map(|h| response.add_header(h));

            let _ = request.respond(response);
        }
    });

    Ok(WebglServerHandle {
        port,
        running,
        thread: Some(thread),
    })
}

#[tauri::command]
fn webgl_status(state: State<WebglServerState>) -> WebglStatus {
    let guard = state.handle.lock().expect("webgl server state poisoned");
    if let Some(handle) = guard.as_ref() {
        webgl_status_from_handle(handle)
    } else {
        WebglStatus {
            running: false,
            url: None,
            port: None,
        }
    }
}

#[tauri::command]
fn webgl_start(app: AppHandle, state: State<WebglServerState>) -> Result<WebglStatus, String> {
    let mut guard = state
        .handle
        .lock()
        .map_err(|_| "state lock failed".to_string())?;

    if let Some(handle) = guard.as_ref() {
        if handle.running.load(Ordering::Relaxed) {
            return Ok(webgl_status_from_handle(handle));
        }
    }

    let root = resolve_breach_it_dir(&app).ok_or_else(|| {
        "WebGL build not found. Expected folder `breach_it/` with `index.html`.".to_string()
    })?;

    let handle = start_webgl_server(root)?;
    let status = webgl_status_from_handle(&handle);
    *guard = Some(handle);
    Ok(status)
}

#[tauri::command]
fn webgl_stop(state: State<WebglServerState>) -> Result<WebglStatus, String> {
    let handle_opt = {
        let mut guard = state
            .handle
            .lock()
            .map_err(|_| "state lock failed".to_string())?;
        guard.take()
    };

    let Some(mut handle) = handle_opt else {
        return Ok(WebglStatus {
            running: false,
            url: None,
            port: None,
        });
    };

    handle.running.store(false, Ordering::Relaxed);
    if let Some(t) = handle.thread.take() {
        let _ = t.join();
    }

    Ok(WebglStatus {
        running: false,
        url: None,
        port: None,
    })
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .manage(WebglServerState::default())
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![
            greet,
            webgl_start,
            webgl_stop,
            webgl_status
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
