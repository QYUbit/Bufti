class VersionError(Exception):
    def __init__(self, *args):
        super().__init__(*args)

class ModelError(Exception):
    def __init__(self, *args):
        super().__init__(*args)

class BufferFormatError(Exception):
    def __init__(self, *args):
        super().__init__(*args)

class DictFormatError(Exception):
    def __init__(self, *args):
        super().__init__(*args)

class EndOfBytesError(Exception):
    def __init__(self, *args):
        super().__init__(*args)

class UnexpectedEndOfBytesError(Exception):
    def __init__(self, *args):
        super().__init__(*args)
