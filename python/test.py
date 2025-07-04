import unittest
from butil import Model, Field, create_list_type, create_map_type, create_model_type, INT32_TYPE, INT64_TYPE, FLOAT64_TYPE, BOOL_TYPE, STRING_TYPE, registered_models

class TestModelFunctionality(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        # Ensure models are created only once
        if "other" not in registered_models:
            cls.other_model = Model("other", Field(0, "aa", FLOAT64_TYPE))
        else:
            cls.other_model = registered_models["other"]
        
        if "myModel" not in registered_models:
            cls.my_model = Model("myModel",
                Field(0, "a", STRING_TYPE),
                Field(1, "b", create_list_type(INT32_TYPE)),
                Field(2, "c", create_map_type(STRING_TYPE, BOOL_TYPE)),
                Field(3, "d", create_model_type("other"))
            )
        else:
            cls.my_model = registered_models["myModel"]
        
        cls.sample_data = {
            "a": "asdasdasd",
            "b": [423, 2324, 444],
            "c": {"iii": True, "uuu": False},
            "d": {"aa": 55.3}
        }
    
    def test_encode(self):
        try:
            buf = self.my_model.encode(self.sample_data)
            self.assertIsInstance(buf, bytearray)
        except Exception as e:
            self.fail(f"Encoding failed with exception: {e}")
    
    def test_decode(self):
        buf = self.my_model.encode(self.sample_data)
        decoded_data = self.my_model.decode(buf)
        self.assertEqual(decoded_data, self.sample_data)
    
    def test_list_encoding_decoding(self):
        list_model = Model("listModel", Field(0, "numbers", create_list_type(INT32_TYPE)))
        sample_list_data = {"numbers": [1, 2, 3, 4, 5]}
        buf = list_model.encode(sample_list_data)
        decoded = list_model.decode(buf)
        self.assertEqual(decoded, sample_list_data)
    
    def test_map_encoding_decoding(self):
        map_model = Model("mapModel", Field(0, "mapping", create_map_type(STRING_TYPE, BOOL_TYPE)))
        sample_map_data = {"mapping": {"key1": True, "key2": False}}
        buf = map_model.encode(sample_map_data)
        decoded = map_model.decode(buf)
        self.assertEqual(decoded, sample_map_data)
    
    def test_model_encoding_decoding(self):
        buf = self.my_model.encode(self.sample_data)
        decoded_data = self.my_model.decode(buf)
        self.assertEqual(decoded_data, self.sample_data)

    def test_empty_list(self):
        list_model = Model("emptyListModel", Field(0, "empty_list", create_list_type(INT32_TYPE)))
        buf = list_model.encode({"empty_list": []})
        decoded = list_model.decode(buf)
        self.assertEqual(decoded, {"empty_list": []})

    def test_empty_map(self):
        map_model = Model("emptyMapModel", Field(0, "empty_map", create_map_type(STRING_TYPE, BOOL_TYPE)))
        buf = map_model.encode({"empty_map": {}})
        decoded = map_model.decode(buf)
        self.assertEqual(decoded, {"empty_map": {}})

    def test_nested_empty_structures(self):
        nested_model = Model("nestedModel", Field(0, "nested_map", create_map_type(STRING_TYPE, create_list_type(INT32_TYPE))))
        buf = nested_model.encode({"nested_map": {}})
        decoded = nested_model.decode(buf)
        self.assertEqual(decoded, {"nested_map": {}})

    def test_large_numbers(self):
        large_number_model = Model("largeNumberModel", Field(0, "big_int", INT64_TYPE))
        buf = large_number_model.encode({"big_int": 9223372036854775807})
        decoded = large_number_model.decode(buf)
        self.assertEqual(decoded, {"big_int": 9223372036854775807})

    #! fails for utf-16 strings
    def test_special_characters(self):
        special_char_model = Model("specialCharModel", Field(0, "text", STRING_TYPE))
        special_text = {"text": "こんにちは世界! 🌍"}
        buf = special_char_model.encode(special_text)
        decoded = special_char_model.decode(buf)
        self.assertEqual(decoded, special_text)

if __name__ == "__main__":
    unittest.main()
