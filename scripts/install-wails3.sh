#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"
wails_module="github.com/wailsapp/wails/v3"

if ! command -v go >/dev/null 2>&1; then
  echo "[Wails3] go is required but was not found in PATH" >&2
  exit 1
fi

version="$(cd "${repo_root}/desktop" && go list -m -f '{{.Version}}' "${wails_module}")"
if [[ -z "${version}" || "${version}" == "<nil>" ]]; then
  echo "[Wails3] failed to resolve ${wails_module} version from desktop/go.mod" >&2
  exit 1
fi

install_target="${wails_module}/cmd/wails3@${version}"
echo "[Wails3] resolved ${wails_module} ${version} from desktop/go.mod"
echo "[Wails3] installing ${install_target}"
go install "${install_target}"
