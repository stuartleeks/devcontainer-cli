set -e

# Fail if build number not set
if [ -z "$BUILD_NUMBER" ]; then
    echo "Env var 'BUILD_NUMBER' must be set for this script to work correctly"
    exit 1
fi

# If running inside CI login to docker
if [ -z ${IS_CI} ]; then
  echo "Not running in CI, skipping CI setup"
else
  if [ -z $IS_PR ] && [[ $BRANCH == "refs/heads/master" ]]; then
    echo "On master setting PUBLISH=true"
    export PUBLISH=true
  else
    echo "Skipping publish as is from PR: $PR_NUMBER or not 'refs/heads/master' BRANCH: $BRANCH"
  fi
fi

echo "git tags:"
git tag --list
echo

if [ -z ${PUBLISH} ]; then
  echo "Running with --skip-publish as PUBLISH not set"
  goreleaser --skip-publish --rm-dist
else
  echo "Publishing release"
  goreleaser
fi