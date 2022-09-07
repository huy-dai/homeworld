# Edit this configuration file to define what should be installed on
# your system.  Help is available in the configuration.nix(5) man page
# and in the NixOS manual (accessible by running ‘nixos-help’).

# After successful boot-up, do:
# echo TOKEN | nixos-kubernetes-node-join
#     Use the token we got from the Master node

# Troubleshooting
# 0. Verify that worker node is able to contact master node
# 1. Verify that certmgr.service is running


{ config, lib, pkgs, ... }:

let
  kubeMasterIP = "10.1.1.2";
  kubeMasterHostname = "hyades-kube-master-vm";
  kubeMasterAPIServerPort = 6443;
  kubeNodeHostname = "hyades-kube-node-vm";
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

  networking.hostName = "${kubeNodeHostname}";
  services.sshd.enable = true;

  # Serve built packages, and also cache requests to cache.nixos.org
  services.nginx.enable = true;
  services.nix-serve = {
    enable = true;
  };

  networking.firewall.allowedTCPPorts = [ 22 8888 kubeMasterAPIServerPort ];

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
  #Code taken from: <https://nixos.wiki/wiki/Kubernetes#>

  # resolve master hostname
  networking.extraHosts = "${kubeMasterIP} ${kubeMasterHostname}";
  
  # packages for administration tasks
  environment.systemPackages = with pkgs; [
    kompose
    kubectl
    kubernetes
  ];

  services.kubernetes = let
    api = "https://${kubeMasterHostname}:${toString kubeMasterAPIServerPort}";
  in
  {
    roles = ["node"];
    masterAddress = kubeMasterHostname;
    easyCerts = true; #Use easyCerts for testing for now - (TODO: Later switch to a more robust solution for production)

    # point kubelet and other services to kube-apiserver
    kubelet.kubeconfig.server = api;
    apiserverAddress = api;

    # use coredns
    addons.dns.enable = true;

    # needed if you use swap
    kubelet.extraOpts = "--fail-swap-on=false"; #TODO: Better practice is to disable swap at host level rather than use this flag
  };
}

# vim:set ts=2 sw=2 et:
