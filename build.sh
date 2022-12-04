VERSION=0.1.0-ci

docker build . -f docker/Dockerfile -t application-template:${VERSION}
helm package helm/application-template --version ${VERSION} --app-version ${VERSION}
