from xml.dom import minidom

import libvirt
import sys

try:
    conn = libvirt.openReadOnly("qemu:///system")
except libvirt.libvirtError:
    print('Failed to open connection to the hypervisor')
    sys.exit(1)

try:
    domain = conn.lookupByName("bohemia")
except libvirt.libvirtError:
    print('Failed to find the main domain')
    sys.exit(1)

print("Domain 0: id %d running %s" % (domain.ID(), domain.OSType()))
print(domain.UUIDString())
print(domain.name())
# print(domain.XMLDesc())

macs = []
xml = minidom.parseString(domain.XMLDesc())
for iface in xml.getElementsByTagName('domain')[0].getElementsByTagName("interface"):
    if iface.getAttribute('type') != "network": continue
    mac = iface.getElementsByTagName("mac")[0].getAttribute("address")
    macs.append(mac)

print(macs)
