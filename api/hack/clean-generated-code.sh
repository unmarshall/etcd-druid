#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
API_GO_MODULE_ROOT="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$API_GO_MODULE_ROOT")"
GENERATED_CLIENT_DIR="${PROJECT_ROOT}/client"

find "$(dirname "$SCRIPT_DIR")" -type f -name "zz_*.go" -exec rm '{}' \;
grep -lr '// Code generated by client-gen. DO NOT EDIT' "$(dirname "$GENERATED_CLIENT_DIR")" | xargs rm -f
grep -lr '// Code generated by informer-gen. DO NOT EDIT' "$(dirname "$GENERATED_CLIENT_DIR")" | xargs rm -f
grep -lr '// Code generated by lister-gen. DO NOT EDIT' "$(dirname "$GENERATED_CLIENT_DIR")" | xargs rm -f
