import tornado.web
from textx import metamodel_from_file, TextXSyntaxError

import settings
from converter.json_converter import convert

rules_mm = metamodel_from_file(settings.GRAMMAR_FILE)


class Handler(tornado.web.RequestHandler):
    @tornado.gen.coroutine
    def post(self, user_id):
        try:
            rules_model = rules_mm.model_from_str(self.request.body)
            for rule in rules_model.rules:
                rule.userId = user_id

            if self.application.nats.is_connected:
                yield self.application.nats.publish(settings.NATS["topic"], convert(rules_model))
                self.set_status(202)
            else:
                self.clear()
                self.set_status(503)

        except TextXSyntaxError:
            self.clear()
            self.set_status(400)

        self.finish()
