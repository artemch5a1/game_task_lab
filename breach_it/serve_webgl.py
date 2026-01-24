#!/usr/bin/env python3
"""
Local static server for Unity WebGL builds with precompressed (.gz/.br) assets.

Why:
- Unity often outputs Build/*.wasm.gz, *.data.gz, *.js.gz
- Many simple servers don't set Content-Encoding for these files,
  causing the browser to fail loading them.

Usage:
  python3 serve_webgl.py
  python3 serve_webgl.py --port 8000 --host 127.0.0.1
Then open:
  http://localhost:8000/
"""

from __future__ import annotations

import argparse
import http.server
import os
import socketserver
from pathlib import Path
from typing import Optional


def _content_type_for(path_without_encoding: str) -> str:
    lower = path_without_encoding.lower()
    if lower.endswith(".wasm"):
        return "application/wasm"
    if lower.endswith(".js"):
        return "application/javascript; charset=utf-8"
    if lower.endswith(".data"):
        return "application/octet-stream"
    if lower.endswith(".json"):
        return "application/json; charset=utf-8"
    if lower.endswith(".css"):
        return "text/css; charset=utf-8"
    if lower.endswith(".html") or lower.endswith(".htm"):
        return "text/html; charset=utf-8"
    if lower.endswith(".png"):
        return "image/png"
    if lower.endswith(".jpg") or lower.endswith(".jpeg"):
        return "image/jpeg"
    if lower.endswith(".svg"):
        return "image/svg+xml"
    if lower.endswith(".ico"):
        return "image/x-icon"
    return "application/octet-stream"


class UnityWebGLRequestHandler(http.server.SimpleHTTPRequestHandler):
    # Reduce noisy logging a bit (still shows requests).
    def log_message(self, fmt: str, *args) -> None:
        super().log_message(fmt, *args)

    def end_headers(self) -> None:
        # Commonly useful for local testing.
        self.send_header("Cache-Control", "no-cache")
        super().end_headers()

    def guess_type(self, path: str) -> str:
        # Called by SimpleHTTPRequestHandler when serving files.
        if path.endswith(".gz"):
            return _content_type_for(path[:-3])
        if path.endswith(".br"):
            return _content_type_for(path[:-3])
        return super().guess_type(path)

    def send_head(self):
        # Largely copied from stdlib, with added Content-Encoding handling.
        path = self.translate_path(self.path)
        f: Optional[object] = None
        if os.path.isdir(path):
            parts = self.path.split("?", 1)
            parts[0] = parts[0].rstrip("/") + "/"
            self.path = parts[0] + ("?" + parts[1] if len(parts) > 1 else "")
            for index in ("index.html", "index.htm"):
                index_path = os.path.join(path, index)
                if os.path.exists(index_path):
                    path = index_path
                    break
            else:
                return self.list_directory(path)

        ctype = self.guess_type(path)
        try:
            f = open(path, "rb")
        except OSError:
            self.send_error(404, "File not found")
            return None

        try:
            fs = os.fstat(f.fileno())
            self.send_response(200)
            self.send_header("Content-type", ctype)
            self.send_header("Content-Length", str(fs.st_size))
            self.send_header("Last-Modified", self.date_time_string(fs.st_mtime))

            if path.endswith(".gz"):
                self.send_header("Content-Encoding", "gzip")
            elif path.endswith(".br"):
                self.send_header("Content-Encoding", "br")

            self.end_headers()
            return f
        except Exception:
            f.close()
            raise


class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="localhost", help="Bind host (default: localhost)")
    parser.add_argument("--port", type=int, default=8000, help="Bind port (default: 8000)")
    args = parser.parse_args()

    root = Path(__file__).resolve().parent
    os.chdir(root)

    with ReusableTCPServer((args.host, args.port), UnityWebGLRequestHandler) as httpd:
        print(f"Serving Unity WebGL from: {root}")
        print(f"Open: http://{args.host}:{args.port}/")
        print("Press Ctrl+C to stop.")
        httpd.serve_forever()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

