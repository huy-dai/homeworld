#!/usr/bin/env bash
set -e

# Prereqs
# - x86_64 host
# - KVM Virtualuzation available (scripts/kvm_ok)
# - Docker

cd "$(dirname "$0")"

containsElement () {
  # args: needle haystack
  local e match="$1"
  shift
  for e; do [[ "$e" == "$match" ]] && return 0; done
  return 1
}

# Grab all available guests from the hosts/ dir
guests=(hosts/*/)
guests=("${guests[@]%/}")     # Remove trailing slash
guests=("${guests[@]##*/}")   # Remove path prefixes

# Check if specified VM exists as a guest configuration
if ! containsElement "$1" "${guests[@]}" ; then
  echo "Requested configuration not found: $1" 1>&2
  echo "Usage: run_sdk_vm.sh <configuration> [qemu_args]"
  echo "Available Configurations:"
  printf '* %s\n' "${guests[@]}"
  exit 1
fi

# Save VM locations
HOMEWORLD_PATH=$(realpath .)
HOMEWORLD_VM_ROOT=${HOMEWORLD_VM_ROOT:-"${HOMEWORLD_PATH}/vm"}
VM_GUEST="$1"
VM_FILE="${HOMEWORLD_VM_ROOT}/$1.qcow2"
VM_SOCKET="${HOMEWORLD_VM_ROOT}/$1.socket"

echo "Running VM: ${VM_FILE}"

if [ ! -f "${VM_FILE}" ]; then
    echo "VM Disk ${VM_FILE} does not exist yet, creating..."

    # TODO: Remove dependency from docker/docker.io
    docker run -it --rm \
                --device /dev/kvm \
                -v "$HOMEWORLD_PATH":/homeworld \
                docker.io/nixos/nix \
                /bin/sh /homeworld/scripts/setup_vm.sh "${VM_GUEST}"

    echo "VM created. Moving to HOMEWORLD_VM_ROOT"
    mkdir -p "${HOMEWORLD_VM_ROOT}"
    mv "${VM_FILE}" "${HOMEWORLD_VM_ROOT}"
fi

# Grab VM-specific arguments
mapfile -t <"hosts/${VM_GUEST}/qemu.args"

# Show the qemu command.
# NOTE: To perform a dry run, pass in invalid arguments in the command line, such as `-bios NONE`
set -o xtrace

# Wrapping the last two arguments in quotes removes the whitespace between args that we want to maintain.
# NOTE: We may want to research looking into using the -fw_cfg argument to pass files into /sys/firmware/qemu-fw-cfg
qemu-system-x86_64 -enable-kvm -nographic \
    -cpu host \
    -m 2048 \
    -drive if=virtio,file="${VM_FILE}" \
    -virtfs local,path="${HOMEWORLD_PATH}",mount_tag=host0,security_model=none,id=host0,readonly=on \
    -serial mon:stdio \
    -monitor unix:"${VM_SOCKET}",server,nowait ${MAPFILE[@]} ${@:2}

# TODO: Research solution to this
# After finishing the script, may need to run 'tput smam' to fix text wrapping
# (https://bugs.launchpad.net/qemu/+bug/1857449)
