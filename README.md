# Project Hardcap

On-premise bare-metal or virtual provisioning for the cloud era with strong focus on Red Hat Anaconda, the installer for Red Hat Enterprise Linux, Fedora and clones.

THIS REPOSITORY IS WORK IN PROGRESS!

## The background

As a software engineer who maintained provisioning capabilities of the [Foreman](https://www.theforeman.org) open-source project for about a decade, I've decided to start a fresh project from scratch. It's not a Foreman replacement at all, actually, this prototype aims to simplify and complement some of Foreman's rich provisioning capabilities. Some day, there might be a Foreman plugin for integration with hardcap.

The main theme is: **simplicity**. 

This document describes some of my initial thoughts of what hardcap could be, it's rather a list of things I would like to avoid, or do differently, than a comprehensive design document or full list of feature. I will update the document as I work on a prototype.

## The architecture

* A **hardcap-controller** is a daemon which holds the whole inventory in a postgres database. It has a gRPC API for a command line client and agents. Multiple controllers can be deployed, but they are completely separate (no data sharing).
* A **hardcap-agent** agent provides plugin gRPC API for inventory operations (list all blades in a rack), power operations (turn on, off) and management operations (boot from network, local drive). An endpoint for rendering templates (Anaconda `kickstart`) or fact updates are available too.
  * A `hca-redfish` plugin provides integration with Redfish hardware.
  * A `hca-libvirt` plugin provides the same but for libvirt.
* A **hardcap-cli** is a command line utility for interacting with multiple controllers (`hcli`).

## The workflow

Hardware must be registered first either manually or automatically via chassis detection over OOB (Redfish), then it's powered on booting into Anaconda which discovers them and provides more hardware details. Cloud portal users initiate OS installation which makes them provisioned until they are acquired for production use. Once released, hosts can be optionally wiped and returned back to the discovered state. The workflow:

1. Host is **registered** by manually entering OOB credentials or via chassis auto-registration. There can be OOB drivers for virtualized environments too (e.g. Intel or z/VM).
2. Host is **discovered** by booting into Anaconda which sends hardware facts and waits for provisioning.
3. Host can be optionally **powered off** until it's requested to be provisioned, or it can stay in the discovered state forever. Host can be optionally provisioned into hardware testing images (firmware update, CPU/memory tests) and then powered off.
4. Host is **provisioned** when Anaconda gets a command for installation with credentials (ssh keys, root password). Any Anaconda provisioning type is supported: RPM, ostree, image-based.
5. Host is **acquired** when the system must be allocated for use, a list of facts about the system is returned (e.g. MAC/IP address). Provisioning and acquisition is usually done in one step.
6. Host is **released** when it is no longer needed. Optional wipe flag can be set.
7. Host is **wiped** after it boots into Anaconda and performs storage cleanup. This is optional. Then it returns to the discovered state.

Hardware detection is implemented as an Anaconda `%pre` install script (Bash/Python) which is important for providers with limited capabilities (IPMI). Host can immediately continue provisioning, or it can be powered off until it is requested by an operator. Notable hardware facts reported:

* MAC addresses
* architecture
* CPUs/cores
* memory
* drives/capacity

## The vision

**Hardware-first approach**. The workflow must be built around bare-metal hosts, their hardware details and capabilities. Instead of "give me five XXL instances", the workflow is "we have these available servers/VMs, make a pick".

**OOB chassis discovery**. New hosts can be registered manually or via OOB chassis discovery for protocols which support it (Redfish). A libvirt implementation will be delivered early in the process to test this feature; customers can use it to actually provision statically-allocated libvirt VMs.

**The OS installer is Anaconda**. The whole workflow is designed around Anaconda's features and capabilities. Hardcap project is focused on provisioning of Red Hat Enterprise Linux, Fedora, CentOS Stream and clones. Other distributions or operating systems can be also installed, but the only option being image-based provisioning. In other words, Anaconda can technically install Ubuntu OS image (untested).

**RPM and image-based provisioning**. Anaconda today supports three ways of installing the OS: traditional (RPM), stateless (OSTree) and image-based (RHV-H image). All the three methods will be supported. All workflows will be unattended out of box.

**Redfish and UEFI first**. It's 2022, and it is safe to say that all modern systems these days support both UEFI (HTTP Boot protocol) and Redfish for out of band management. This is the primary focus, BIOS/PXE will be possible tho.

**Inventory of hardware, not software**. Hardcap's work is done, when OS is installed and optionally a configuration management software is bootstrapped (Ansible has a valid ssh connection or Puppet agent was installed). No details about installed operating system, sofware, repositories or patches should be ever stored in hardcap's database. There are other projects to do that.

**No NIC inventory management**. Instead of trying to model host-NIC-network complex configurations (bond, team, vlan, vlan over bond, two bond pairs), hardcap gives users the power of kickstart. Out of box, it works just fine with a single active network connection.

**No built-in authorization**. While authentication is a hard requirement, authorization (users, roles, organizations, permissions) does not map well for auto-discovery hardware workflow (discovered hosts by definition belong to subnets not users) and makes things unnecessarily complex for ad-hoc provisioning via the CLI. There are no plans for a complex authorization model.

**Federated deployment**. One or multiple controllers can be deployed and connected to one to many agents providing a gateway into backend services like Redfish, VMWare, libvirt or IPMI.

**Extensibility in mind**. Extensibility is at the front, all functionality will be delivered as plugins written either in Go or Python. Any language supporting gRPC protocol can be used to write new OOB agent plugins.

**Limited support for hosts without BMC**. Out-of-band interface is needed for the main workflow, however, there will be possible to register and provision hosts without BMC/Redfish interfaces. There will be limitations tho.

**Anaconda remote logging**. Anaconda supports several remote management features like VNC (outband or inband connection), remote syslog and ssh daemon. These features should be fully utilized from the day one and users should be able to use them for debugging purposes.

**Simple templating**. The aim is to have a very simple templating capabilities that is easy to understand. Little bit of copying is sometimes better than very complicated model with lots of options which is hard to maintain and keep documentation for.

**Snippets support**. Any template with "snippet" flag (or tag) will be considered a snippet and any snippet statement would be open to be overridden by users. This will be useful for network and storage parts of kickstarts.

**Dynamic kickstart generation**. Anaconda offers a nice trick of generating kickstart on the fly depending on hardware configuration (e.g. writing network or storage statements after hardware is initialized for all devices found). Some examples will be shipped for more complex storage setups.

**Rendered templates kept in database**. To be able to reprovision a host with the very same configuration, kickstart must be kept in the database.

**No IPAM/DHCP/DNS/VLAN management**. Hardcap is not a DHCP/DNS management software, an IPAM or SDN solution. It assumes that the network infrastructure is managed with different tools. The API should provide good integration options (e.g. create DHCP leases or DNS records).

**Hook system**. Various events will be exposed via a hook gRPC API so that integration with 3rd party software is possible (e.g. create a DNS record and set a custom key-value with the hostname).

**All-in-one installation**. Although the project aims for simple deployment, there are couple of components needed in order to provision hosts (a controller, an agent, postgres, DHCP). Installation will be done via Ansible and there will be a single-node installation with dnsmasq DHCP service for small-scale customers out of box.

**No boot and image files management**. Network bootloader configuration is always static (boot into Anaconda), there is no need for TFTP/PXE files management. Also, repositories or/and OS images must not be handled by this service. Instead, it will be possible to integrate this service with Pulp, Foreman or Satellite when needed.

**Cockpit ready**. There is no plan for writing a web UI, however, a simple installation and a good API could lead to a cockpit UI for small-scale networks with all-in-one style deployment.

**Testing clusters**. Although this was not the primary goal, there's a chance that the hardcap's design will be a good fit for beaker-like deployments.

## Plugin interface

The API is specifically designed as a gRPC which allows writing of plugins in any language. Here is a rough vision of how such an API would look like:

## Inventory

Database model stored on each controller:

### Host

* friendly id hc-XXXXXXX (globally-unique autogenerated identifier)
* friendly name (globally-unique autogenerated name)
* name (rack/chassis or VM name)
* hardware facts
  * MAC addresses (list)
  * architecture
  * CPUs/cores
  * memory
  * drives/capacity
* state (registered, discovered, provisioned, acquired, released, wiped)
* power (on, off)
* boot (local, network)
* tags
  * architecture
  * bare-metal or hypervisor name
  * controller name
  * agent name

## Key technologies

* Anaconda
* Go
* Postgres / pgx + scany
* gRPC
* zerolog
* Redfish, IPMI
* libvirt, VMWare
* Python (for plugins)

### Anaconda

Red Hat OS installer, Anaconda, is the primary installer for hardcap, although the project should work with any installer that provides a capability of sending a HW/VM serial number (or a MAC address) via HTTP either as a URL path element, parameter or HTTP header.

Anaconda does this plus it is a mature linux installer with advanced features like automation via kickstart or image-based and rpm-ostree deployment.

## Further reading

* https://theforeman.org/
* https://maas.io/
* https://github.com/redhat-performance/badfish
* https://github.com/hashicorp/go-plugin (plugins API)

