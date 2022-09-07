

# Edit this configuration file to define what should be installed on
# your system.  Help is available in the configuration.nix(5) man page
# and in the NixOS manual (accessible by running ‘nixos-help’).

# After successful boot-up, do:
# mkdir ~/.kube
# ln -s /etc/kubernetes/cluster-admin.kubeconfig ~/.kube/config
# kubectl cluster-info 
#     Verify that Kubernetes master is running
# kubectl get nodes
#     Verify that master is identified as a node
# cat /var/lib/kubernetes/secrets/apitoken.secret
#     Save this to use for the Worker node to join the cluster


{ config, lib, pkgs, ... }:

let
  kubeMasterIP = "10.1.1.2";
  kubeMasterHostname = "hyades-kube-master-vm";
  kubeMasterAPIServerPort = 6443;
in
{
  fileSystems."/homeworld" = {
    device = "host0";
    fsType = "9p";
    options = [
      "trans=virtio" "version=9p2000.L" "ro" "_netdev"
    ];
  };

  time.timeZone = "America/New_York";

  networking.hostName = "${kubeMasterHostname}";
  services.sshd.enable = true;

  # Serve built packages, and also cache requests to cache.nixos.org
  services.nginx.enable = true;
  services.nix-serve = {
    enable = true;
  };

  networking.firewall.allowedTCPPorts = [ 22 kubeMasterAPIServerPort ];

  users.users.root.password = "root";
  services.openssh.permitRootLogin = lib.mkDefault "yes";
  services.getty.autologinUser = lib.mkDefault "root";
  
  # This value determines the NixOS release from which the default
  # settings for stateful data, like file locations and database versions
  # on your system were taken. It‘s perfectly fine and recommended to leave
  # this value at the release version of the first install of this system.
  # Before changing this value read the documentation for this option
  # (e.g. man configuration.nix or on https://nixos.org/nixos/options.html).
  system.stateVersion = "22.05"; # Did you read the comment?

  system.copySystemConfiguration = true;

  #Set up: Kubernetes master node
  #Code taken from: <https://nixos.wiki/wiki/Kubernetes#Rook_Ceph_storage_cluster>

  # resolve master hostname
  networking.extraHosts = "${kubeMasterIP} ${kubeMasterHostname}";
  
  # packages for administration tasks
  environment.systemPackages = with pkgs; [
    kompose
    kubectl
    kubernetes
  ];

   services.kubernetes = {
    roles = ["master" "node"];
    masterAddress = kubeMasterHostname;
    apiserverAddress = "https://${kubeMasterHostname}:${toString kubeMasterAPIServerPort}";
    easyCerts = true; #Use easyCerts for testing for now - (TODO: Later switch to a more robust solution for production)
    apiserver = {
      securePort = kubeMasterAPIServerPort;
      advertiseAddress = kubeMasterIP;
    };

    # use coredns
    addons.dns.enable = true;

    # needed if you use swap
    kubelet.extraOpts = "--fail-swap-on=false"; #TODO: Better practice is to disable swap at host level rather than use this flag
   };
}

# vim:set ts=2 sw=2 et:
