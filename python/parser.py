from errors import EndOfBytesError, UnexpectedEndOfBytesError
import struct

class Parser():
    def __init__(self, buffer: bytearray):
        self.off = 0
        self.buf = buffer

    def get_bytes(self) -> bytearray:
        return self.buf
    
    def print(self) -> None:
        print("".join("{:02x} ".format(x).upper() for x in self.buf))

    def read_bytes(self, length: int) -> bytearray:
        if self.off == len(self.buf) and length != 0:
            raise EndOfBytesError()
        if self.off + length > len(self.buf):
            raise UnexpectedEndOfBytesError(len(self.buf), self.off, length)
        
        buffer = self.buf[self.off:self.off + length]
        self.off += length
        return buffer
    
    def read_string(self, size: int) -> str:
        return self.read_bytes(size).decode()

    def read_uint8(self) -> int:
        return struct.unpack("B", self.read_bytes(1))[0]
    
    def read_uint16(self) -> int:
        return struct.unpack("!H", self.read_bytes(2))[0]
    
    def read_uint32(self) -> int:
        return struct.unpack("!L", self.read_bytes(4))[0]
    
    def read_uint64(self) -> int:
        return struct.unpack("!Q", self.read_bytes(8))[0]
    
    def read_int8(self) -> int:
        return struct.unpack("b", self.read_bytes(1))[0]
    
    def read_int16(self) -> int:
        return struct.unpack("!h", self.read_bytes(2))[0]
    
    def read_int32(self) -> int:
        return struct.unpack("!l", self.read_bytes(4))[0]
    
    def read_int64(self) -> int:
        return struct.unpack("!q", self.read_bytes(8))[0]
    
    def read_float32(self) -> float:
        return struct.unpack("!f", self.read_bytes(4))[0]
    
    def read_float64(self) -> float:
        return struct.unpack("!d", self.read_bytes(8))[0]
    
    def read_bool(self) -> bool:
        return struct.unpack("?", self.read_bytes(1))[0]
    
    def write_uint8(self, value: int) -> None:
        self.off = 0
        self.buf.extend(struct.pack("B", value))

    def write_uint16(self, value: int) -> None:
        self.off = 0
        self.buf.extend(struct.pack("!H", value))

    def write_uint32(self, value: int) -> None:
        self.off = 0
        self.buf.extend(struct.pack("!L", value))
    
    def write_uint64(self, value: int) -> None:
        self.off = 0
        self.buf.extend(struct.pack("!Q", value))

    def write_int8(self, value: int) -> None:
        self.off = 0
        self.buf.extend(struct.pack("b", value))

    def write_int16(self, value: int) -> None:
        self.off = 0
        self.buf = struct.pack("!h", value)

    def write_int32(self, value: int) -> None:
        self.off = 0
        self.buf.extend(struct.pack("!l", value))
    
    def write_int64(self, value: int) -> None:
        self.off = 0
        self.buf.extend(struct.pack("!q", value))

    def write_float32(self, value: float) -> None:
        self.off = 0
        self.buf.extend(struct.pack("!f", value))

    def write_float64(self, value: float) -> None:
        self.off = 0
        self.buf.extend(struct.pack("!d", value))

    def write_bool(self, value: bool) -> None:
        self.off = 0
        self.buf.extend(struct.pack("?", value))

    def write_bytes(self, value: bytes) -> None:
        self.off = 0
        self.buf.extend(value)

    def write_string(self, value: str) -> None:
        self.off = 0
        self.buf.extend(bytes(value, "utf-8"))
