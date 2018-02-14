import json

_primitives = [int, float, str, unicode, bool]


def convert(model):
    return json.dumps(_general_handler(model), indent=4)


def _tx_filter(to_filter):
    return not (to_filter[0].startswith("_") or "parent" == to_filter[0])


def _list_handler(to_handle):
    return map(_general_handler, to_handle)


def _object_handler(to_handle):
    handled = {}
    attributes_filtered = filter(_tx_filter, to_handle.__dict__.iteritems())
    elements_created = map(_element_creator, attributes_filtered)
    map(lambda x: handled.update(x), elements_created)
    return handled


def _element_creator(source):
    return {source[0]: _general_handler(source[1])}


def _general_handler(to_handle):
    if type(to_handle) is list:
        return _list_handler(to_handle)
    elif type(to_handle) in _primitives:
        return to_handle
    else:
        return _object_handler(to_handle)
