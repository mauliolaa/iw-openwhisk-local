def factorial(n)
    if n == 1 || n == 0
        return 1
    else
        return n * factorial(n-1)
    end
end

def main(args)
    n = args["n"] || 1
    result = factorial(n)
    { "result" => result }
  end