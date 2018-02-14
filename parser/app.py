import tornado.gen
import tornado.httpserver
import tornado.ioloop
import tornado.web
from nats.io import Client as NATSClient
from tornado.options import options

from settings import tornado_settings, nats_options
from urls import url_patterns


class DSLHandlerApp(tornado.web.Application):
    def __init__(self):
        tornado.web.Application.__init__(self, url_patterns, **tornado_settings)
        self.nats = NATSClient()
        self.nats.connect(**nats_options)


def main():
    app = DSLHandlerApp()
    http_server = tornado.httpserver.HTTPServer(app)
    http_server.listen(options.port)
    tornado.ioloop.IOLoop.current().start()


if __name__ == "__main__":
    main()
