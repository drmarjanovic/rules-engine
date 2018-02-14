import tornado.web

import settings


class Handler(tornado.web.RequestHandler):
    def get(self):
        self.write({"name": "parser", "version": settings.VERSION})
