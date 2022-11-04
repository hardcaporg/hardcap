#!/usr/bin/env python3
# This script is executed in Anaconda kickstart %pre section and then contents
# of /tmp/pre-generated.ks file is included and interpreted as the kickstart.
# Keep this Python 2 and Python 3 compatible.

from glob import glob
from pprint import pp
from subprocess import Popen, PIPE


def run(*cmd):
    try:
        output = Popen(cmd, stdout=PIPE)
        return output.communicate()
    except FileNotFoundError:
        return ""


def ks_write(line):
    with open('/tmp/pre-generated.ks', 'a') as ks:
        ks.write(line)


def gather_mac():
    macs = []
    for name in glob("/sys/class/net/*/address"):
        macs.append(open(name).readline())
    return macs


def gather_facts():
    facts = {
        "mac": gather_mac(),
        "dmi": {},
    }
    for keyword in ['bios-vendor', 'bios-version', 'bios-release-date', 'bios-revision', 'firmware-revision',
                    'system-manufacturer', 'system-product-name', 'system-version', 'system-serial-number',
                    'system-uuid', 'system-sku-number', 'system-family', 'baseboard-manufacturer',
                    'baseboard-product-name', 'baseboard-version', 'baseboard-serial-number', 'baseboard-asset-tag',
                    'chassis-manufacturer', 'chassis-type', 'chassis-version', 'chassis-serial-number',
                    'chassis-asset-tag', 'processor-family', 'processor-manufacturer', 'processor-version',
                    'processor-frequency']:
        facts["dmi"][keyword] = run("dmidecode", "-s", keyword)
    return facts


pp(gather_facts())

# ks_write('shutdown')
