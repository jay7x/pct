#!/usr/bin/env bash
set -euo pipefail

ORG='jay7x'
REPO='pct'
APP_PKG_NAME='pct'

ARCH=''
OS=''
CHECKSUM=''
FILENAME=''
VERSION=''

cleanup() {
	rm -f /tmp/pct_install_*
}
trap cleanup EXIT

logDebug() {
	if [ -n "${PCT_INSTALL_DEBUG-}" ]; then
		echo "$1" >&2
	fi
}

setVars() {
	case "$(arch)" in
	'aarch64' | 'arm64') ARCH='arm64' ;;
	'x86_64' | 'amd64') ARCH='amd64' ;;
	*)
		echo "Unsupported arch '$(arch)'!" >&2
		exit 1
		;;
	esac
	OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
}

fetchUrl() {
	local uri="$1"; shift
	local file="$1"; shift

	resp=$(curl -sL "$uri" -o "$file" --write-out "%{http_code}")
	logDebug "GET ${uri} | Response: ${resp}"

	if [ "$resp" -ne 200 ]; then
		echo "Unable to fetch '${uri}'" >&2
		exit 1
	fi
}

getChecksums() {
	local checksumURL="https://github.com/${ORG}/${REPO}/releases/latest/download/checksums.txt"
	local filenameRegex="${APP_PKG_NAME}_([0-9.]+)_${OS}_${ARCH}.tar.gz"
	local tmpFile="/tmp/pct_install_checksums.txt"

	fetchUrl "$checksumURL" "$tmpFile"

	CHECKSUM=$(awk "/  ${filenameRegex}/{print \$1}" < "$tmpFile")
	FILENAME=$(awk "/  ${filenameRegex}/{print \$2}" < "$tmpFile")
	VERSION=$(awk -F'_' '{print $2}' <<< "$FILENAME")
}

downloadLatestRelease() {
	local destination="${HOME}/.puppetlabs/pct"
	local downloadURL="https://github.com/${ORG}/${REPO}/releases/download/v${VERSION}/${FILENAME}"
	local tmpFile="/tmp/${FILENAME}"

	mkdir -p "$destination"

	echo "Downloading and extracting ${APP_PKG_NAME} v${VERSION} to ${destination}..." >&2

	fetchUrl "$downloadURL" "$tmpFile"

	actual=$(shasum -a 256 "$tmpFile" 2>/dev/null || sha256sum "$tmpFile")
	actualChecksum="${actual%% *}"

	if [ "$actualChecksum" != "$CHECKSUM" ]; then
		echo "Checksum verification failed for ${FILENAME}" >&2
		exit 1
	fi

	logDebug "Extracting ${tmpFile} to ${destination}"
	tar -zxf "$tmpFile" -C "$destination" || {
		echo "Untar unsuccessful" >&2
		exit 1
	}

	echo "Remember to add the pct app to your path:"
	# shellcheck disable=SC2016
	echo 'export PATH=$PATH:'"${destination}"
}

setVars
getChecksums
downloadLatestRelease
