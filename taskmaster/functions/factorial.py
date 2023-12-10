def factorial(n):
    if n <= 1:
        return n
    return n * factorial(n - 1)

def main(args):
    n = args.get("n", 0)
    return {"greeting": factorial(n)}