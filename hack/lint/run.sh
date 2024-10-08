#!/usr/bin/env bash
# Copyright 2019 Iguazio
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

OS_NAME=$(uname -s)
OS_NAME_LOWERCASE=$(echo "${OS_NAME}" | tr "[:upper:]" "[:lower:]")

if [[ -z "${BIN_DIR}" ]]; then
  BIN_DIR=$(pwd)/.bin
fi

echo "Linting @$(pwd)..."
"${BIN_DIR}"/golangci-lint run -v --max-same-issues=100
