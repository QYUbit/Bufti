import { ModelError, PayloadError } from "./errors";

interface SchemaField {
    label: string;
    type: string;
}

interface Field {
    index: number;
    label: string;
    type: string;
}

type ModelRegister = Record<string, Model>;
const modelRegister: ModelRegister = {};

export default class Model {
    readonly name: string;
    private schema: {[index: number]: SchemaField} = {};
    private labels: {[label: string]: number} = {};

    constructor(name: string, ...fields: Field[]) {
        if (name === "") {
            throw new ModelError(`empty model name`);
        }
        if (modelRegister[name]) {
            throw new ModelError(`Model ${name} already exists`);
        }

        this.name = name;

        for (const field of fields) {
            if (field.label === "") {
                throw new ModelError(`empty label in model ${this.name}`);
            }
            if (field.index > 255 || field.index < 0) {
                throw new ModelError(`index not between 0 and 255 in model ${this.name}, instead ${field.index}`);
            }
            if (this.labels[field.label]) {
                throw new ModelError(`duplicate label ${field.label} in model ${this.name}`);
            }
            if (this.schema[field.index]) {
                throw new ModelError(`duplicate index ${field.index} in model ${this.name}`);
            }

            this.labels[field.label] = field.index;
            this.schema[field.index] = {label: field.label, type: field.type};
        }
    }

    public encode(map: {[key: string]: any}): Uint8Array {
        const buf = Buffer.from([]);
        this.encodePayload(buf, map);
        return new Uint8Array(buf.buffer);
    }

    private encodePayload(buf: Buffer, map: {[key: string]: any}) {
        for (const key in map) {
            if (!this.labels[key]) {
                throw new PayloadError(`label ${key} not found`);
            }
            const index = this.labels[key];
            if (!this.schema[index]) {
                throw new PayloadError(`index ${index} not found`);
            }
            const schemaField = this.schema[index];

            buf.writeUint8(index);
            this.encodeValue(buf, map[key], schemaField.type);
        }
    }

    private encodeValue(buf: Buffer, value: any, type: string) {
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
                buf.writeUint8(value?1:0);
            case "string":
                buf.write(value, "utf-8");
        }
    }
}

const myModel = new Model("myModel",
    {index: 0, label: "abc", type: "int32"}
);

const someMap = {abc: 82765};

const buf = myModel.encode(someMap);
