"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PayloadError = exports.ModelError = void 0;
class ModelError extends Error {
    constructor(message) {
        super(`Invalid model: ${message}`);
    }
}
exports.ModelError = ModelError;
class PayloadError extends Error {
    constructor(message) {
        super(`Unexpected payload format: ${message}`);
    }
}
exports.PayloadError = PayloadError;
