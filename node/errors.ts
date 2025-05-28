export class ModelError extends Error {
    constructor(message: string) {
        super(`Invalid model: ${message}`);
    }
}

export class PayloadError extends Error {
    constructor(message: string) {
        super(`Unexpected payload format: ${message}`);
    }
}
