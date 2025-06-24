from dataclasses import dataclass
from parser import Parser
from errors import ModelError, DictFormatError, EndOfBytesError, BufferFormatError
from typing import Any

MAJOR_VERSION = 0

INT8_TYPE = "int8"
INT16_TYPE = "int16"
INT32_TYPE = "int32"
INT64_TYPE = "int64"
FLOAT32_TYPE = "float32"
FLOAT64_TYPE = "float64"
BOOL_TYPE = "bool"
STRING_TYPE = "string"

def create_list_type(element_type: str) -> str:
    return f"list:{element_type}"

def create_map_type(key_type: str, value_type: str) -> str:
    return f"map:{key_type}:{value_type}"

def create_model_type(model_name: str) -> str:
    return f"model:{model_name}"

@dataclass
class Field:
    index: int
    label: str
    type: str

@dataclass
class FieldSchema:
    label: str
    type: str

class Model:
    def __init__(self, name: str, *fields: Field) -> None:
        self.name = name

        if list(registered_models.keys()).count(self.name) > 0:
            raise ModelError(f"duplicate model {self.name}")

        self.labels: dict[str, int] = {}
        self.schema: dict[int, FieldSchema] = {}

        for field in fields:
            if list(self.schema.keys()).count(field.index) > 0:
                raise ModelError(f"duplicate index {field.index} in model {self.name}")
            
            if list(self.labels.keys()).count(field.label) > 0:
                raise ModelError(f"duplicate label {field.label} in model {self.name}")
            
            if field.label == "":
                raise ModelError(f"empty label in model {self.name}")
            
            if not 0 <= field.index <= 255:
                raise ModelError(f"index not between 0 and 255 in model {self.name} in field {field.label}, instead {field.index}")

            self.labels[field.label] = field.index
            self.schema[field.index] = FieldSchema(field.label, field.type)

        registered_models[self.name] = self

    def encode(self, values: dict[str, Any]) -> bytes:
        parser = Parser(bytearray())
        self._encode(parser, values)
        return parser.get_bytes()

    def _encode(self, parser: Parser, values: dict[str, Any]) -> None:
        for label, value in values.items():
            try:
                index = self.labels[label]
            except KeyError:
                raise DictFormatError(f"label {label} not found in model {self.name}")
            
            try:
                field = self.schema[index]
            except KeyError:
                raise DictFormatError(f"index {index} not found in model {self.name}")
            
            parser.write_uint8(index)
            self._encode_value(parser, value, field.type)

    def _encode_value(self, parser: Parser, value: Any, field_type: str) -> None:
        if field_type == INT8_TYPE:
            parser.write_int8(value)
        if field_type == INT16_TYPE:
            parser.write_int16(value)
        if field_type == INT32_TYPE:
            parser.write_int32(value)
        if field_type == INT64_TYPE:
            parser.write_int64(value)
        if field_type == FLOAT32_TYPE:
            parser.write_float32(value)
        if field_type == FLOAT64_TYPE:
            parser.write_float64(value)
        if field_type == BOOL_TYPE:
            parser.write_bool(value)
        if field_type == STRING_TYPE:
            parser.write_uint32(len(value))
            parser.write_string(value)

        if field_type.startswith("list:"):
            element_type = field_type.removeprefix("list:")
            parser.write_uint32(len(value))
            for element in value:
                self._encode_value(parser, element, element_type)

        if field_type.startswith("map:"):
            parts = field_type.split(":")
            if len(parts) < 3:
                raise ModelError(f"invalid map type in model {self.name} ({field_type})")
            parser.write_uint32(len(value))

            for key, val in value.items():
                self._encode_value(parser, key, parts[1])
                self._encode_value(parser, val, parts[2])

        if field_type.startswith("model:"):
            model_name = field_type.removeprefix("model:")
            try:
                model = registered_models[model_name]
            except KeyError:
                raise ModelError(f"no model {model_name} registered")
            
            parser.write_uint16(len(value))
            model._encode(parser, value)

    def decode(self, buffer: bytes) -> dict[str, object]:
        parser = Parser(bytearray(buffer))
        return self._decode(parser, len(buffer))

    def _decode(self, parser: Parser, limit: int) -> dict[str, Any]:
        bu: dict[str, Any] = {}
        for _ in range(limit):
            try:
                index = parser.read_uint8()
            except EndOfBytesError:
                break

            try:
                field = self.schema[index]
            except KeyError:
                raise BufferFormatError(f"index {index} not found in model {self.name}")
            
            val = self._decode_value(parser, field.type)
            bu[field.label] = val
        return bu
        
    def _decode_value(self, parser: Parser, field_type: str) -> Any:
        if field_type == INT8_TYPE:
            return parser.read_int8()
        if field_type == INT16_TYPE:
            return parser.read_int16()
        if field_type == INT32_TYPE:
            return parser.read_int32()
        if field_type == INT64_TYPE:
            return parser.read_int64()
        if field_type == FLOAT32_TYPE:
            return parser.read_float32()
        if field_type == FLOAT64_TYPE:
            return parser.read_float64()
        if field_type == BOOL_TYPE:
            return parser.read_bool()
        if field_type == STRING_TYPE:
            size = parser.read_uint32()
            return parser.read_string(size)
        
        if field_type.startswith("list:"):
            element_type = field_type.removeprefix("list:")
            size = parser.read_uint32()

            parse_list: list[Any] = []
            for _ in range(size):
                item = self._decode_value(parser, element_type)
                parse_list.append(item)
            return parse_list

        if field_type.startswith("map:"):
            parts = field_type.split(":")
            if len(parts) < 3:
                raise ModelError(f"invalid map type in model {self.name} ({field_type})")
            size = parser.read_uint32()

            parse_dict: dict[Any, Any] = {}
            for _ in range(size):
                key = self._decode_value(parser, parts[1])
                val = self._decode_value(parser, parts[2])
                parse_dict[key] = val
            return parse_dict

        if field_type.startswith("model:"):
            model_name = field_type.removeprefix("model:")
            model = registered_models[model_name]
            if not model:
                raise ModelError(f"no model {model_name} registered")
            size = parser.read_uint16()
            return model._decode(parser, size)
        

registered_models: dict[str, Model] = {} #! global bad, integrate into model
