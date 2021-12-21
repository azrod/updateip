#!/bin/bash

version=$1
REGISTRY_USERNAME=$2
REGISTRY_TOKEN=$3
FROM="gcr.io/distroless/static"

platforms="linux/arm64 linux/amd64 linux/arm/v6 linux/arm/v7 darwin/amd64"

binary="updateip"

golang_version=$(curl -s "https://go.dev/VERSION?m=text" | sed -e 's/go//' | cut -f1,2 -d.)

#Setting up some colors for helping read the demo output
bold=$(tput bold)
red=$(tput setaf 1)
green=$(tput setaf 2)
yellow=$(tput setaf 3)
blue=$(tput setaf 4)
cyan=$(tput setaf 6)
reset=$(tput sgr0)

buildah rmi "localhost/${binary}:${version}"
buildah manifest create "localhost/${binary}:${version}"

echo "${green}Build Golang${reset}"
build=$(buildah from docker.io/golang:"${golang_version}")
buildah config --workingdir /go/src/work "$build"
buildah copy "$build" . .
for platform in $platforms; do
	GOOS=$(echo "$platform" | cut -f1 -d/)
	GOARCH=$(echo "$platform" | cut -f2 -d/)
	GOARM=$(echo "$platform" | cut -f3 -d/ | sed "s/v//")
	VARIANT="--variant $(echo "$platform" | cut -f3 -d/)"
	if [[ -z "$GOARM" ]]; then
		VARIANT=""
	fi

	if [ "$GOARCH" == "arm64" ]; then
		VARIANT="--variant 8"
	fi

	if [ "$GOOS" == "darwin" ]; then
		FROM="scratch"
	fi

	binary_suffix="$GOARCH$(echo "$platform" | cut -f3 -d/)"
	echo "${green}go build...${reset}"
	buildah run "$build" /bin/sh -c "CGO_ENABLED=0 GOOS=$GOOS GOARCH=${GOARCH} go build -ldflags=\"-w -s\"  -o $binary.$binary_suffix"

	echo "${cyan}Create container for ${platform}${reset}"
	rtcntr=$(buildah --os "$GOOS" --arch "$GOARCH" $VARIANT from "$FROM")
	rtmnt=$(buildah mount "$rtcntr")
	buildmnt=$(buildah mount "$build")
	cp "${buildmnt}/go/src/work/${binary}.${binary_suffix}" "${rtmnt}/${binary}"

	buildah config --entrypoint "[\"/${binary}\"]" "${rtcntr}"
	buildah config --label org.opencontainers.image.source="https://github.com/azrod/${binary}" "${rtcntr}"
	buildah config --label maintainer="Azrod <contact@mickael-stanislas.com>" "${rtcntr}"
	buildah commit --rm --squash --manifest "${binary}" "${rtcntr}" "${binary}:${GOOS}-${GOARCH}-${version}"
done

# manifests
echo "${blue}push to github${reset}"
buildah manifest push --creds "$REGISTRY_USERNAME":"$REGISTRY_TOKEN" --all "localhost/${binary}:${version}" docker://ghcr.io/"${REGISTRY_USERNAME}/${binary}:${version}"
if [ "$version" != "dev" ]; then
  buildah manifest push --creds "$REGISTRY_USERNAME":"$REGISTRY_TOKEN" --all "localhost/${binary}:${version}" docker://ghcr.io/"${REGISTRY_USERNAME}/${binary}:latest"
fi

echo "${red}Cleanup...${reset}"
echo "cleanup"
buildah umount --all
buildah rm --all