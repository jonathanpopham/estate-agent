from __future__ import annotations

import argparse


def main() -> None:
    parser = argparse.ArgumentParser(description="Estate Agent")
    subcommands = parser.add_subparsers(dest="command", required=True)

    serve = subcommands.add_parser("serve", help="Run the HTTP server")
    serve.add_argument("--host", default="127.0.0.1")
    serve.add_argument("--port", default=8080, type=int)
    serve.add_argument("--reload", action="store_true")

    args = parser.parse_args()

    if args.command == "serve":
        try:
            import uvicorn
        except ImportError as exc:
            raise SystemExit(
                "Server dependencies are missing. Install with: pip install -e '.[server]'"
            ) from exc

        uvicorn.run(
            "estate_agent.server:app",
            host=args.host,
            port=args.port,
            reload=args.reload,
        )


if __name__ == "__main__":
    main()

