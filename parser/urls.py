from handlers import health, rules

url_patterns = [
    (r"/health", health.Handler),
    (r"/users/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})/rules", rules.Handler)
]
