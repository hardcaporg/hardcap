import os
import re
import sys
from struct import pack, unpack
from typing import Callable, BinaryIO, Dict, Tuple
import cbor2
from xml.dom import minidom
import libvirt
import sys


def log(*a):
    print(*a, file=sys.stderr)


class Server(object):
    _methods = {}

    def __new__(cls, *args, **kwargs):
        obj = super(Server, cls).__new__(cls, *args, **kwargs)
        obj.__dict__ = cls._methods
        return obj

    def register(self, name: str, method: Callable[[Dict], Tuple[Dict, str]]):
        self._methods[name] = method

    def call(self, name: str, args: Dict) -> Tuple[Dict, str]:
        return self._methods[name].__call__(args)


class Request:
    def __init__(self):
        self.header = {}
        self.args = {}

    @staticmethod
    def read_or_raise(io: BinaryIO, n: int):
        buf = io.read(n)
        if len(buf) == 0:
            raise EOFError
        return buf

    @staticmethod
    def load(io: BinaryIO):
        self = Request()
        size = unpack("<I", Request.read_or_raise(io, 4))[0]
        self.header = cbor2.loads(Request.read_or_raise(io, size))
        size = unpack("<I", Request.read_or_raise(io, 4))[0]
        self.args = cbor2.loads(Request.read_or_raise(io, size))
        return self

    def service_method(self):
        return self.header["ServiceMethod"]

    def seq(self):
        return self.header["Seq"]

    def __repr__(self):
        from pprint import pformat
        return pformat(vars(self))


class Response(Request):
    def __init__(self, request: Request, reply: any, error: str):
        super().__init__()
        self.header = request.header
        self.args = request.args
        self.reply = reply
        self.header["Error"] = error

    def error(self):
        return self.header["Error"]

    def dumps(self, io: BinaryIO):
        header = cbor2.dumps(self.header)
        reply = cbor2.dumps(self.reply)
        io.write(pack("<I", len(header)))
        io.write(header)
        io.write(pack("<I", len(reply)))
        io.write(reply)


def multiply(args: Dict) -> Tuple[Dict, str]:
    reply = {"C": args["A"] * args["B"]}
    return reply, ""


def enlist(args: Dict) -> Tuple[Dict, str]:
    url = args["URL"]
    systems = []
    reply = {"Systems": systems}

    try:
        conn = libvirt.openReadOnly(url)
    except libvirt.libvirtError as err:
        return {}, f'failed to open connection to the hypervisor: {err}'

    try:
        for domain in conn.listAllDomains():
            if not re.match(args["NamePattern"], domain.name()):
                continue

            log("Domain 0: id %d running %s" % (domain.ID(), domain.OSType()))
            log(domain.UUIDString())
            log(domain.name())
            # log(domain.XMLDesc())

            macs = []
            xml = minidom.parseString(domain.XMLDesc())
            for iface in xml.getElementsByTagName('domain')[0].getElementsByTagName("interface"):
                if iface.getAttribute('type') != "network": continue
                mac = iface.getElementsByTagName("mac")[0].getAttribute("address")
                macs.append(mac)

            log(macs)

            s = {
                "Name": domain.name(),
                "UID": domain.UUIDString(),
                "MACs": macs,
            }
            systems.append(s)

    except libvirt.libvirtError as err:
        return {}, f'failed to list intances: {err}'

    return reply, ""


server = Server()
server.register("Arith.Multiply", multiply)
server.register("Appliance.Enlist", enlist)

if __name__ == '__main__':
    if os.environ["PLUGIN"]:
        log("Starting libvirt appliance plugin")
        while not sys.stdin.buffer.closed:
            try:
                req = Request.load(sys.stdin.buffer)
                (result, error) = server.call(req.service_method(), req.args)
                resp = Response(req, result, error)
                resp.dumps(sys.stdout.buffer)
                sys.stdout.buffer.flush()
                log(f"Dispatched {req.service_method()} call")
            except EOFError:
                log("Pipeline closed, exiting")
                break
            except Exception:
                raise
    else:
        log("Must be called from hardcap")
