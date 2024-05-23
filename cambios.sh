#!/bin/bash

readonly current_dir=$( cd "$(dirname "$0")" ; pwd -P )
readonly command_output_file="/tmp/cmd.out"
readonly cambios_binary="${current_dir}/cambios"
readonly error_cambios_binary_does_not_exist=80

process_command_file() {
    local command=$(awk -F '|' '{print $1}' "${command_output_file}")
    local repo_path=$(awk -F '|' '{print $2}' "${command_output_file}")
    local file_path=$(awk -F '|' '{print $3}' "${command_output_file}")
    case "${command}" in
        "diff")
            git -C "${repo_path}" -p diff "${file_path}"
            ;;
        "v")
            (
                cd "${repo_path}"
                bat "${file_path}"
            )
            ;;
        "clipboard")
            echo "${file_path}" | pbcopy
            ;;
    esac
}

if ((${#} != 1)); then
    echo "error: ${0} <repository path>"

    exit 1
fi

if [[ -f "${command_output_file}" ]]; then
    rm "${command_output_file}" > /dev/null 2>&1
fi

if [[ ! -f "${cambios_binary}" ]]; then
    echo "error: ${cambios_binary} not exist" >&2

    exit ${error_cambios_binary_does_not_exist}
fi

readonly repository_path="${1}"

echo "Running cambios binary"
if ! "${cambios_binary}" "${repository_path}"; then
    echo "error: running ${cambios_binary}"
fi
echo "Done running cambios binary"

if [[ -f "${command_output_file}" ]]; then
    process_command_file
fi

if [[ -f "${command_output_file}" ]]; then
    rm "${command_output_file}" > /dev/null 2>&1
fi

exit 0