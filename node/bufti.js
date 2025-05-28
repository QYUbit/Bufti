"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const errors_1 = require("./errors");
const modelRegister = {};
class Model {
    constructor(name, ...fields) {
        this.schema = {};
        this.labels = {};
        if (name === "") {
            throw new errors_1.ModelError(`empty model name`);
        }
        if (modelRegister[name]) {
            throw new errors_1.ModelError(`Model ${name} already exists`);
        }
        this.name = name;
        for (const field of fields) {
            if (field.label === "") {
                throw new errors_1.ModelError(`empty label in model ${this.name}`);
            }
            if (field.index > 255 || field.index < 0) {
                throw new errors_1.ModelError(`index not between 0 and 255 in model ${this.name}, instead ${field.index}`);
            }
            if (this.labels[field.label]) {
                throw new errors_1.ModelError(`duplicate label ${field.label} in model ${this.name}`);
            }
            if (this.schema[field.index]) {
                throw new errors_1.ModelError(`duplicate index ${field.index} in model ${this.name}`);
            }
            this.labels[field.label] = field.index;
            this.schema[field.index] = { label: field.label, type: field.type };
        }
    }
    encode(map) {
        const buf = Buffer.from([]);
        this.encodePayload(buf, map);
        return new Uint8Array(buf.buffer);
    }
    encodePayload(buf, map) {
        for (const key in map) {
            if (!this.labels[key]) {
                throw new errors_1.PayloadError(`label ${key} not found`);
            }
            const index = this.labels[key];
            if (!this.schema[index]) {
                throw new errors_1.PayloadError(`index ${index} not found`);
            }
            const schemaField = this.schema[index];
            buf.writeUint8(index);
            this.encodeValue(buf, map[key], schemaField.type);
        }
    }
    encodeValue(buf, value, type) {
        switch (type) {
            case "int8":
                buf.writeInt8(value);
            case "int16":
                buf.writeInt16BE(value);
            case "int32":
                buf.writeInt32BE(value);
            case "int64":
                buf.writeBigInt64BE(value);
            case "float32":
                buf.writeFloatBE(value);
            case "float64":
                buf.writeDoubleBE(value);
            case "bool":
                buf.writeUint8(value ? 1 : 0);
            case "string":
                buf.write(value, "utf-8");
        }
    }
}
exports.default = Model;
const myModel = new Model("myModel", { index: 0, label: "abc", type: "int32" });
const someMap = { abc: 82765 };
const buf = myModel.encode(someMap);
