#!/usr/bin/env python3
# This script is executed in Anaconda kickstart %pre section and then contents
# of /tmp/pre-generated.ks file is included and interpreted as the kickstart.
# Keep this Python 2 and Python 3 compatible.

from glob import glob
from pprint import pp
from subprocess import Popen, PIPE
import json
import os
import requests
from syslog import syslog

address = "{{ .Address }}"
log = []

# dev environment overrides
if address.startswith("{{"): address = "localhost:8000"

def run(*cmd):
    try:
        proc = Popen(cmd, stdout=PIPE, stderr=PIPE)
        proc.wait()
        stdout = proc.stdout.read().decode().strip()
        stderr = proc.stderr.read().decode().strip()
        if len(stderr) > 0: log_write(" ".join(cmd), stderr)
        return stdout.strip()
    except FileNotFoundError:
        return ""

def log_write(prefix, message):
    global log
    log.append([prefix, str(message)])
    syslog(': '.join(["hardcap", prefix, str(message)]))

def ks_write(line):
    with open('/tmp/pre-generated.ks', 'a') as ks:
        ks.write(line)


def gather_mac():
    macs = []
    for name in glob("/sys/class/net/*/address"):
        mac = open(name).readline().strip()
        if len(mac) > 0: macs.append(mac)
    return macs


def gather_facts():
    global log
    facts = {
        "mac": gather_mac(),
        "cpu": {
            # TODO try with psutil package contains a lot of useful stuff
            "count": open('/proc/cpuinfo').read().count('processor\t:'),
        },
        "memory": {
            "bytes": os.sysconf('SC_PAGE_SIZE') * os.sysconf('SC_PHYS_PAGES'),
        },
        "dmi": {},
    }
    for keyword in ['bios-vendor', 'bios-version', 'bios-release-date', 'bios-revision', 'firmware-revision',
                    'system-manufacturer', 'system-product-name', 'systemK-version', 'system-serial-number',
                    'system-uuid', 'system-sku-number', 'system-family', 'baseboard-manufacturer',
                    'baseboard-product-name', 'baseboard-version', 'baseboard-serial-number', 'baseboard-asset-tag',
                    'chassis-manufacturer', 'chassis-type', 'chassis-version', 'chassis-serial-number',
                    'chassis-asset-tag', 'processor-family', 'processor-manufacturer', 'processor-version',
                    'processor-frequency']:
        facts["dmi"][keyword] = run("dmidecode", "-s", keyword)
    facts["log"] = log
    return facts

facts = gather_facts()
print(json.dumps(facts, indent=2))

r = requests.post('http://%s/r/host_register' % address, json=facts)
log_write("register upload", r.content)

# ks_write('shutdown')
