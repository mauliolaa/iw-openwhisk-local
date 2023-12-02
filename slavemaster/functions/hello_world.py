# name: str
# place: str
def main(args):
    name = args.get("name", "stranger")
    place = args.get("place", "universe")
    greeting = f"Hello from {name} at {place}!"
    return {"greeting": greeting}