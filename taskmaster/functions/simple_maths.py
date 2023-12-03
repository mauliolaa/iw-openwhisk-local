def main(args):
    x = args.get("x", 0)
    y = args.get("y", 0)
    return {'result': str(x + y)}