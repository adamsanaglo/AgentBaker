#!/bin/bash

set -e

VHD_LOGS_FILEPATH=/opt/azure/vhd-install.complete

required_env_vars=(
    "IMG_SKU"
    "SKU_NAME"
)

for v in "${required_env_vars[@]}"
do
    if [ -z "${!v}" ]; then
        echo "$v was not set!"
        exit 1
    else
        echo "$v is set to '${!v}'"
    fi
done

FEATURE_FLAGS="${FEATURE_FLAGS:-}"

echo "Starting COSI ACL build on $(date)" > ${VHD_LOGS_FILEPATH}

source /opt/azure/containers/provision_source.sh

# NOTE: Unlike the azlosguard VHD postinstall, we do NOT flatten the A/B repart
# layout here. The base osguard image's dual-partition (usr-a/usr-b) scheme is
# required for over-the-wire COSI updates managed by Trident.

# Setup a symlink for lg-redirect-sysext
mkdir -p /etc/extensions/lg-redirect-sysext/usr/local/
mkdir -p /opt/bin
ln -s /opt/bin /etc/extensions/lg-redirect-sysext/usr/local/bin
# Bind mount /opt/bin to /usr/local/bin during the build since the redirect sysext is not running
mount --bind /opt/bin /usr/local/bin
trap "umount /usr/local/bin" EXIT

# Link /opt/azure/containers to /home/packer for postinstall
ln -s /opt/azure/containers /home/packer

### pre-install-dependencies ###
echo -e "\nnews.none                          -/var/log/messages" >> /etc/rsyslog.d/60-CIS.conf
# Create dir for update_certs.path
mkdir -p /opt/certs
chmod 755 /opt/certs
# Use AKS Log Collector instead of WALA log collections
echo -e "\n# Disable WALA log collection because AKS Log Collector is installed.\nLogs.Collect=n" >> /etc/waagent.conf

# Use fetch-only mode to avoid unpacking images, saving disk space on COSI's constrained root partition
export IMAGE_FETCH_ONLY=true

source /opt/azure/containers/tool_installs.sh
source /opt/azure/containers/tool_installs_distro.sh

echo "VHD build dev - basicass set up" >> ${VHD_LOGS_FILEPATH}

chmod a-x /opt/azure/containers/provision_start.sh
chmod a-x /opt/azure/containers/provision.sh
chmod a-x /opt/azure/containers/provision_installs.sh
chmod a-x /opt/azure/containers/provision_installs_distro.sh
chmod a-x /opt/azure/containers/provision_configs.sh
chmod a-x /opt/azure/containers/provision_source.sh
chmod a-x /opt/azure/containers/provision_source_distro.sh

/opt/azure/containers/install-dependencies.sh

echo "Build completed successfully on $(date)" >> ${VHD_LOGS_FILEPATH}
