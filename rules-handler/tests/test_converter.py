from unittest import TestCase

from textx import metamodel_from_file

import converter.json_converter as converter

converted_positive = """
{
    "rules": [
        {
            "conditions": [
                {
                    "operator": ">=", 
                    "parameter": {
                        "property": "temperature", 
                        "deviceId": "a32db207-7236-4e75-abad-7c972f4cfd18"
                    }, 
                    "value": 30
                }
            ], 
            "name": "ruleTemp01", 
            "actions": [
                {
                    "name": "TURN OFF", 
                    "parameters": {
                        "deviceId": "937a7c3e-db39-4d75-b52b-b8442463761a"
                    }
                }, 
                {
                    "name": "SEND EMAIL", 
                    "parameters": {
                        "content": "<b>Alert!</b><br><p>High temperature in living room!</p>", 
                        "recipient": "person01@home.com"
                    }
                }, 
                {
                    "name": "SEND EMAIL", 
                    "parameters": {
                        "content": "<b>Alert!</b><br><p>High temperature in living room!</p>", 
                        "recipient": "person02@home.com"
                    }
                }
            ]
        }
    ]
}""".strip()


class JsonConverter(TestCase):
    class TestExample(object):
        def __init__(self, a, b):
            self.a = a
            self.b = b

    def setUp(self):
        self.object_example = self.TestExample("test", 0)
        self.rules_mm = metamodel_from_file("spec/rules_grammar.tx")
        self.rules = self.rules_mm.model_from_str(
            """
            RULE ruleTemp01:
                a32db207-7236-4e75-abad-7c972f4cfd18["temperature"] >= 30
            TRIGGERS
                TURN OFF 937a7c3e-db39-4d75-b52b-b8442463761a
                SEND EMAIL "<b>Alert!</b><br><p>High temperature in living room!</p>" TO "person01@home.com"
                SEND EMAIL "<b>Alert!</b><br><p>High temperature in living room!</p>" TO "person02@home.com"
            """)

    def test_element_creator(self):
        self.assertEqual({"attribute": "value"}, converter._element_creator(("attribute", "value")))

    def test_object_handler(self):
        self.assertEqual({"a": "test", "b": 0}, converter._object_handler(self.object_example))

    def test_tx_filter_positive(self):
        self.assertTrue(converter._tx_filter(("valid_attribute",)))

    def test_tx_filter_negative(self):
        self.assertFalse(converter._tx_filter(("_private",)))
        self.assertFalse(converter._tx_filter(("parent",)))

    def test_list_handler(self):
        self.assertTrue({"test": [{"first": "first"}, {"second": "second"}]},
                        converter._list_handler(
                            [self.TestExample("first", "first"), self.TestExample("second", "second")]))

    def test_general_handler(self):
        self.assertEqual({"a": "test", "b": "test"}, converter._general_handler(self.TestExample("test", "test")))
        self.assertEqual("test", converter._general_handler("test"))
        self.assertTrue({"test": [{"first": "first"}, {"second": "second"}]},
                        converter._general_handler(
                            [self.TestExample("first", "first"), self.TestExample("second", "second")]))

    def test_converter(self):
        self.assertEqual(converted_positive, converter.convert(self.rules))
